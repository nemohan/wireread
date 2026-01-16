package wireread

import (
	"bytes"
	"encoding/binary"
	"io"
)

// SafeReader is a safe implementation of Reader with complete boundary checking.
// It validates all read operations and returns errors when data is insufficient.
type SafeReader struct {
	data []byte
	size int
	rpos int
}

// NewSafeReader creates a new SafeReader for the given data.
// All read operations will be validated for boundary conditions.
func NewSafeReader(data []byte) *SafeReader {
	return &SafeReader{
		data: data,
		size: len(data),
		rpos: 0,
	}
}

// Bytes returns the remaining unparsed bytes from the current read position
func (sr *SafeReader) Bytes() []byte {
	return sr.data[sr.rpos:]
}

func (sr *SafeReader) ReadBytes(n int) ([]byte, error) {
	if len(sr.data[sr.rpos:]) < n {
		return nil, io.ErrUnexpectedEOF
	}
	dest := make([]byte, n)
	copy(dest, sr.data[sr.rpos:])
	sr.rpos += n
	return dest, nil
}

func (sr *SafeReader) ReadByte() (byte, error) {
	if sr.rpos+1 > sr.size {
		return 0, io.ErrUnexpectedEOF
	}
	tmp := sr.data[sr.rpos]
	sr.rpos++
	return tmp, nil
}

func (sr *SafeReader) Skip(n int) error {
	if sr.rpos+n > sr.size {
		return io.ErrUnexpectedEOF
	}
	sr.rpos += n
	return nil
}

func (sr *SafeReader) ReadUvarint() (uint64, error) {
	return binary.ReadUvarint(sr)
}

// ReadString reads n bytes and returns them as a string
func (sr *SafeReader) ReadString(n int) (string, error) {
	if n == 0 {
		return "", nil
	}
	if sr.rpos+n > sr.size {
		return "", io.ErrUnexpectedEOF
	}

	result := string(sr.data[sr.rpos : sr.rpos+n])
	sr.rpos += n
	return result, nil
}

// ReadStringInto reads n bytes into the provided string pointer
func (sr *SafeReader) ReadStringInto(out *string, n int) error {
	result, err := sr.ReadString(n)
	if err != nil {
		return err
	}
	*out = result
	return nil
}

// ReadNullTerminatedString reads a null-terminated string (C-style string)
func (sr *SafeReader) ReadNullTerminatedString() (string, error) {
	for i, b := range sr.data[sr.rpos:] {
		if b == 0 {
			result := string(sr.data[sr.rpos : sr.rpos+i])
			sr.rpos += i + 1
			return result, nil
		}
	}
	return "", io.ErrUnexpectedEOF
}

// ReadLengthEncodedInteger reads a MySQL length-encoded integer
func (sr *SafeReader) ReadLengthEncodedInteger() (uint64, error) {
	if len(sr.data[sr.rpos:]) == 0 {
		return 0, io.ErrUnexpectedEOF
	}

	data := sr.data[sr.rpos:]
	switch data[0] {
	case 0xFB: // NULL
		sr.rpos++
		return 0, nil
	case 0xFC: // 2-byte integer
		if len(data) < 3 {
			return 0, io.ErrUnexpectedEOF
		}
		sr.rpos += 2
		return uint64(binary.LittleEndian.Uint16(data[1:3])), nil
	case 0xFD: // 3-byte integer
		if len(data) < 4 {
			return 0, io.ErrUnexpectedEOF
		}
		sr.rpos += 3
		return uint64(data[1]) | uint64(data[2])<<8 | uint64(data[3])<<16, nil
	case 0xFE: // 8-byte integer
		if len(data) < 9 {
			return 0, io.ErrUnexpectedEOF
		}
		sr.rpos += 8
		return binary.LittleEndian.Uint64(data[1:9]), nil
	default: // 1-byte integer
		sr.rpos++
		return uint64(data[0]), nil
	}
}

func (sr *SafeReader) ReadLine() (string, error) {
	begin := sr.rpos
	idx := bytes.Index(sr.data[sr.rpos:], []byte{'\n'})
	if idx < 0 {
		return "", io.ErrUnexpectedEOF
	}
	end := sr.rpos + idx
	if idx > 0 && sr.data[end-1] == '\r' {
		end--
	}
	sr.rpos += idx + 1
	return string(sr.data[begin:end]), nil
}

// ReadUint16BE reads a 16-bit unsigned integer in big-endian byte order
func (sr *SafeReader) ReadUint16BE() (uint16, error) {
	if sr.rpos+2 > sr.size {
		return 0, io.ErrUnexpectedEOF
	}
	tmp := binary.BigEndian.Uint16(sr.data[sr.rpos:])
	sr.rpos += 2
	return tmp, nil
}

// ReadUint16BEInto reads a 16-bit unsigned integer in big-endian byte order into the provided pointer
func (sr *SafeReader) ReadUint16BEInto(out *uint16) error {
	tmp, err := sr.ReadUint16BE()
	if err != nil {
		return err
	}
	*out = tmp
	return nil
}

// ReadInt16BEInto reads a 16-bit signed integer in big-endian byte order into the provided pointer
func (sr *SafeReader) ReadInt16BEInto(out *int16) error {
	tmp, err := sr.ReadUint16BE()
	if err != nil {
		return err
	}
	*out = int16(tmp)
	return nil
}

// ReadUint32BE reads a 32-bit unsigned integer in big-endian byte order
func (sr *SafeReader) ReadUint32BE() (uint32, error) {
	if sr.rpos+4 > sr.size {
		return 0, io.ErrUnexpectedEOF
	}

	tmp := binary.BigEndian.Uint32(sr.data[sr.rpos:])
	sr.rpos += 4
	return tmp, nil
}

// ReadUint32BEInto reads a 32-bit unsigned integer in big-endian byte order into the provided pointer
func (sr *SafeReader) ReadUint32BEInto(out *uint32) error {
	tmp, err := sr.ReadUint32BE()
	if err != nil {
		return err
	}
	*out = tmp
	return nil
}

// ReadInt32BEInto reads a 32-bit signed integer in big-endian byte order into the provided pointer
func (sr *SafeReader) ReadInt32BEInto(out *int32) error {
	tmp, err := sr.ReadUint32BE()
	if err != nil {
		return err
	}
	*out = int32(tmp)
	return nil
}

// ReadUint64BE reads a 64-bit unsigned integer in big-endian byte order
func (sr *SafeReader) ReadUint64BE() (uint64, error) {
	if sr.rpos+8 > sr.size {
		return 0, io.ErrUnexpectedEOF
	}
	sr.rpos += 8
	return binary.BigEndian.Uint64(sr.data[sr.rpos-8:]), nil
}

// ReadUint64BEInto reads a 64-bit unsigned integer in big-endian byte order into the provided pointer
func (sr *SafeReader) ReadUint64BEInto(out *uint64) error {
	tmp, err := sr.ReadUint64BE()
	if err != nil {
		return err
	}
	*out = tmp
	return nil
}

// ReadUint32LE reads a 32-bit unsigned integer in little-endian byte order
func (sr *SafeReader) ReadUint32LE() (uint32, error) {
	if sr.rpos+4 > sr.size {
		return 0, io.ErrUnexpectedEOF
	}

	tmp := binary.LittleEndian.Uint32(sr.data[sr.rpos:])
	sr.rpos += 4
	return tmp, nil
}

// ReadUint32LEInto reads a 32-bit unsigned integer in little-endian byte order into the provided pointer
func (sr *SafeReader) ReadUint32LEInto(out *uint32) error {
	tmp, err := sr.ReadUint32LE()
	if err != nil {
		return err
	}
	*out = tmp
	return nil
}

// ReadUint16LE reads a 16-bit unsigned integer in little-endian byte order
func (sr *SafeReader) ReadUint16LE() (uint16, error) {
	if sr.rpos+2 > sr.size {
		return 0, io.ErrUnexpectedEOF
	}
	tmp := binary.LittleEndian.Uint16(sr.data[sr.rpos:])
	sr.rpos += 2
	return tmp, nil
}

// ReadUint16LEInto reads a 16-bit unsigned integer in little-endian byte order into the provided pointer
func (sr *SafeReader) ReadUint16LEInto(out *uint16) error {
	tmp, err := sr.ReadUint16LE()
	if err != nil {
		return err
	}
	*out = tmp
	return nil
}

// ReadUint64LE reads a 64-bit unsigned integer in little-endian byte order
func (sr *SafeReader) ReadUint64LE() (uint64, error) {
	if sr.rpos+8 > sr.size {
		return 0, io.ErrUnexpectedEOF
	}
	sr.rpos += 8
	return binary.LittleEndian.Uint64(sr.data[sr.rpos-8:]), nil
}

// ReadUint64LEInto reads a 64-bit unsigned integer in little-endian byte order into the provided pointer
func (sr *SafeReader) ReadUint64LEInto(out *uint64) error {
	tmp, err := sr.ReadUint64LE()
	if err != nil {
		return err
	}
	*out = tmp
	return nil
}
