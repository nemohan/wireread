package wireread

import (
	"bytes"
	"encoding/binary"
)

// FastReader is a high-performance reader for complete, trusted data frames.
// It skips all boundary checks and assumes the data is always sufficient.
// WARNING: This reader will panic if data is incomplete. Only use when you are
// certain the data frame is complete and valid.
//
// Use cases:
//   - Parsing pre-validated message frames
//   - Processing data that has already been length-checked
//   - High-performance scenarios where data integrity is guaranteed
type FastReader struct {
	data []byte
	rpos int
}

// NewFastReader creates a new FastReader for the given complete data frame.
// The caller must ensure the data is complete and valid.
func NewFastReader(data []byte) *FastReader {
	return &FastReader{
		data: data,
		rpos: 0,
	}
}

// Bytes returns the remaining unparsed bytes from the current read position
func (fr *FastReader) Bytes() []byte {
	return fr.data[fr.rpos:]
}

// ReadBytes reads n bytes from the buffer without boundary checks
func (fr *FastReader) ReadBytes(n int) ([]byte, error) {
	dest := make([]byte, n)
	copy(dest, fr.data[fr.rpos:])
	fr.rpos += n
	return dest, nil
}

// ReadByte reads a single byte without boundary checks
func (fr *FastReader) ReadByte() (byte, error) {
	b := fr.data[fr.rpos]
	fr.rpos++
	return b, nil
}

// Skip skips n bytes in the buffer without boundary checks
func (fr *FastReader) Skip(n int) error {
	fr.rpos += n
	return nil
}

// ReadUvarint reads a variable-length unsigned integer
func (fr *FastReader) ReadUvarint() (uint64, error) {
	return binary.ReadUvarint(fr)
}

// ReadString reads n bytes and returns them as a string without boundary checks
func (fr *FastReader) ReadString(n int) (string, error) {
	if n == 0 {
		return "", nil
	}
	result := string(fr.data[fr.rpos : fr.rpos+n])
	fr.rpos += n
	return result, nil
}

// ReadStringInto reads n bytes into the provided string pointer without boundary checks
func (fr *FastReader) ReadStringInto(out *string, n int) error {
	if n == 0 {
		*out = ""
		return nil
	}
	*out = string(fr.data[fr.rpos : fr.rpos+n])
	fr.rpos += n
	return nil
}

// ReadNullTerminatedString reads a null-terminated string (C-style string)
func (fr *FastReader) ReadNullTerminatedString() (string, error) {
	for i, b := range fr.data[fr.rpos:] {
		if b == 0 {
			result := string(fr.data[fr.rpos : fr.rpos+i])
			fr.rpos += i + 1
			return result, nil
		}
	}
	// If no null terminator found, return rest of data
	result := string(fr.data[fr.rpos:])
	fr.rpos = len(fr.data)
	return result, nil
}

// ReadLengthEncodedInteger reads a MySQL length-encoded integer without boundary checks
func (fr *FastReader) ReadLengthEncodedInteger() (uint64, error) {
	b := fr.data[fr.rpos]
	switch b {
	case 0xFB: // NULL
		fr.rpos++
		return 0, nil
	case 0xFC: // 2-byte integer
		fr.rpos++
		val := uint64(binary.LittleEndian.Uint16(fr.data[fr.rpos:]))
		fr.rpos += 2
		return val, nil
	case 0xFD: // 3-byte integer
		fr.rpos++
		val := uint64(fr.data[fr.rpos]) | uint64(fr.data[fr.rpos+1])<<8 | uint64(fr.data[fr.rpos+2])<<16
		fr.rpos += 3
		return val, nil
	case 0xFE: // 8-byte integer
		fr.rpos++
		val := binary.LittleEndian.Uint64(fr.data[fr.rpos:])
		fr.rpos += 8
		return val, nil
	default: // 1-byte integer
		fr.rpos++
		return uint64(b), nil
	}
}

// ReadLine reads a line terminated by \n (handles \r\n) without boundary checks
func (fr *FastReader) ReadLine() (string, error) {
	begin := fr.rpos
	idx := bytes.Index(fr.data[fr.rpos:], []byte{'\n'})
	if idx < 0 {
		// No newline found, return rest of data
		result := string(fr.data[fr.rpos:])
		fr.rpos = len(fr.data)
		return result, nil
	}
	end := fr.rpos + idx
	if idx > 0 && fr.data[end-1] == '\r' {
		end--
	}
	fr.rpos += idx + 1
	return string(fr.data[begin:end]), nil
}

// ReadUint16BE reads a 16-bit unsigned integer in big-endian byte order
func (fr *FastReader) ReadUint16BE() (uint16, error) {
	val := binary.BigEndian.Uint16(fr.data[fr.rpos:])
	fr.rpos += 2
	return val, nil
}

// ReadUint16BEInto reads a 16-bit unsigned integer in big-endian byte order into the provided pointer
func (fr *FastReader) ReadUint16BEInto(out *uint16) error {
	*out = binary.BigEndian.Uint16(fr.data[fr.rpos:])
	fr.rpos += 2
	return nil
}

// ReadInt16BEInto reads a 16-bit signed integer in big-endian byte order into the provided pointer
func (fr *FastReader) ReadInt16BEInto(out *int16) error {
	*out = int16(binary.BigEndian.Uint16(fr.data[fr.rpos:]))
	fr.rpos += 2
	return nil
}

// ReadUint32BE reads a 32-bit unsigned integer in big-endian byte order
func (fr *FastReader) ReadUint32BE() (uint32, error) {
	val := binary.BigEndian.Uint32(fr.data[fr.rpos:])
	fr.rpos += 4
	return val, nil
}

// ReadUint32BEInto reads a 32-bit unsigned integer in big-endian byte order into the provided pointer
func (fr *FastReader) ReadUint32BEInto(out *uint32) error {
	*out = binary.BigEndian.Uint32(fr.data[fr.rpos:])
	fr.rpos += 4
	return nil
}

// ReadInt32BEInto reads a 32-bit signed integer in big-endian byte order into the provided pointer
func (fr *FastReader) ReadInt32BEInto(out *int32) error {
	*out = int32(binary.BigEndian.Uint32(fr.data[fr.rpos:]))
	fr.rpos += 4
	return nil
}

// ReadUint64BE reads a 64-bit unsigned integer in big-endian byte order
func (fr *FastReader) ReadUint64BE() (uint64, error) {
	val := binary.BigEndian.Uint64(fr.data[fr.rpos:])
	fr.rpos += 8
	return val, nil
}

// ReadUint64BEInto reads a 64-bit unsigned integer in big-endian byte order into the provided pointer
func (fr *FastReader) ReadUint64BEInto(out *uint64) error {
	*out = binary.BigEndian.Uint64(fr.data[fr.rpos:])
	fr.rpos += 8
	return nil
}

// ReadUint16LE reads a 16-bit unsigned integer in little-endian byte order
func (fr *FastReader) ReadUint16LE() (uint16, error) {
	val := binary.LittleEndian.Uint16(fr.data[fr.rpos:])
	fr.rpos += 2
	return val, nil
}

// ReadUint16LEInto reads a 16-bit unsigned integer in little-endian byte order into the provided pointer
func (fr *FastReader) ReadUint16LEInto(out *uint16) error {
	*out = binary.LittleEndian.Uint16(fr.data[fr.rpos:])
	fr.rpos += 2
	return nil
}

// ReadUint32LE reads a 32-bit unsigned integer in little-endian byte order
func (fr *FastReader) ReadUint32LE() (uint32, error) {
	val := binary.LittleEndian.Uint32(fr.data[fr.rpos:])
	fr.rpos += 4
	return val, nil
}

// ReadUint32LEInto reads a 32-bit unsigned integer in little-endian byte order into the provided pointer
func (fr *FastReader) ReadUint32LEInto(out *uint32) error {
	*out = binary.LittleEndian.Uint32(fr.data[fr.rpos:])
	fr.rpos += 4
	return nil
}

// ReadUint64LE reads a 64-bit unsigned integer in little-endian byte order
func (fr *FastReader) ReadUint64LE() (uint64, error) {
	val := binary.LittleEndian.Uint64(fr.data[fr.rpos:])
	fr.rpos += 8
	return val, nil
}

// ReadUint64LEInto reads a 64-bit unsigned integer in little-endian byte order into the provided pointer
func (fr *FastReader) ReadUint64LEInto(out *uint64) error {
	*out = binary.LittleEndian.Uint64(fr.data[fr.rpos:])
	fr.rpos += 8
	return nil
}
