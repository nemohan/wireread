// Package wireread provides high-performance readers for binary wire protocol data.
//
// This package offers two implementations of the Reader interface:
//   - SafeReader: Complete boundary checking with detailed error handling
//   - FastReader: Zero-overhead parsing for pre-validated data frames
//
// Both implementations support reading various data types in different endianness
// (big-endian and little-endian), including integers, strings, variable-length
// integers, and protocol-specific formats like MySQL length-encoded integers.
//
// Example usage with SafeReader:
//
//	data := []byte{0x00, 0x01, 0x02, 0x03}
//	reader := wireread.NewSafeReader(data)
//	value, err := reader.ReadUint16BE()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Example usage with FastReader:
//
//	data := []byte{0x00, 0x01, 0x02, 0x03}
//	reader := wireread.NewFastReader(data)
//	value, _ := reader.ReadUint16BE() // No error checking for performance
package wireread

// Reader defines the interface for reading wire protocol data.
// It provides methods for reading various data types from a byte buffer
// with support for different byte orders and protocol-specific formats.
type Reader interface {
	// Bytes returns the remaining unparsed bytes from the current read position
	Bytes() []byte

	// ReadBytes reads n bytes from the buffer
	ReadBytes(n int) ([]byte, error)

	// ReadByte reads a single byte
	ReadByte() (byte, error)

	// Skip skips n bytes in the buffer
	Skip(n int) error

	// ReadUvarint reads a variable-length unsigned integer
	ReadUvarint() (uint64, error)

	// ReadString reads n bytes and returns them as a string
	ReadString(n int) (string, error)
	// ReadStringInto reads n bytes into the provided string pointer
	ReadStringInto(out *string, n int) error

	// ReadNullTerminatedString reads a null-terminated string (C-style string)
	ReadNullTerminatedString() (string, error)

	// ReadLengthEncodedInteger reads a MySQL length-encoded integer
	ReadLengthEncodedInteger() (uint64, error)

	// ReadLine reads a line terminated by \n (handles \r\n)
	ReadLine() (string, error)

	// Big Endian read methods (BE = Big Endian)
	// ReadUint16BE reads a 16-bit unsigned integer in big-endian byte order
	ReadUint16BE() (uint16, error)
	// ReadUint16BEInto reads a 16-bit unsigned integer in big-endian byte order into the provided pointer
	ReadUint16BEInto(out *uint16) error
	// ReadInt16BEInto reads a 16-bit signed integer in big-endian byte order into the provided pointer
	ReadInt16BEInto(out *int16) error

	// ReadUint32BE reads a 32-bit unsigned integer in big-endian byte order
	ReadUint32BE() (uint32, error)
	// ReadUint32BEInto reads a 32-bit unsigned integer in big-endian byte order into the provided pointer
	ReadUint32BEInto(out *uint32) error
	// ReadInt32BEInto reads a 32-bit signed integer in big-endian byte order into the provided pointer
	ReadInt32BEInto(out *int32) error

	// ReadUint64BE reads a 64-bit unsigned integer in big-endian byte order
	ReadUint64BE() (uint64, error)
	// ReadUint64BEInto reads a 64-bit unsigned integer in big-endian byte order into the provided pointer
	ReadUint64BEInto(out *uint64) error

	// Little Endian read methods (LE = Little Endian)
	// ReadUint16LE reads a 16-bit unsigned integer in little-endian byte order
	ReadUint16LE() (uint16, error)
	// ReadUint16LEInto reads a 16-bit unsigned integer in little-endian byte order into the provided pointer
	ReadUint16LEInto(out *uint16) error

	// ReadUint32LE reads a 32-bit unsigned integer in little-endian byte order
	ReadUint32LE() (uint32, error)
	// ReadUint32LEInto reads a 32-bit unsigned integer in little-endian byte order into the provided pointer
	ReadUint32LEInto(out *uint32) error

	// ReadUint64LE reads a 64-bit unsigned integer in little-endian byte order
	ReadUint64LE() (uint64, error)
	// ReadUint64LEInto reads a 64-bit unsigned integer in little-endian byte order into the provided pointer
	ReadUint64LEInto(out *uint64) error
}
