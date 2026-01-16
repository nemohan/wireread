# wireread

A high-performance Go library for reading binary wire protocol data with both safe and fast implementations.

[![Go Reference](https://pkg.go.dev/badge/github.com/nemohan/wireread.svg)](https://pkg.go.dev/github.com/nemohan/wireread)
[![Go Report Card](https://goreportcard.com/badge/github.com/nemohan/wireread)](https://goreportcard.com/report/github.com/nemohan/wireread)

## Features

- 🛡️ **SafeReader**: Complete boundary checking with detailed error handling
- ⚡ **FastReader**: Zero-overhead parsing for pre-validated data
- 🔄 **Flexible Endianness**: Support for both big-endian and little-endian byte orders
- 📦 **Rich Data Types**: Integers, strings, variable-length integers, and protocol-specific formats
- 🎯 **Interface-based**: Both implementations satisfy the same `Reader` interface
- 🚀 **Zero Dependencies**: Only uses Go standard library

## Installation

```bash
go get github.com/nemohan/wireread
```

## Quick Start

### SafeReader - With Error Checking

Use `SafeReader` when you need robust error handling:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/nemohan/wireread"
)

func main() {
    data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
    reader := wireread.NewSafeReader(data)
    
    // Read a 16-bit unsigned integer (big-endian)
    value, err := reader.ReadUint16BE()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Value: 0x%04x\n", value) // Output: Value: 0x0001
    
    // Read a 32-bit unsigned integer (little-endian)
    value32, err := reader.ReadUint32LE()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Value32: 0x%08x\n", value32) // Output: Value32: 0x05040302
}
```

### FastReader - High Performance

Use `FastReader` when data integrity is guaranteed and performance is critical:

```go
package main

import (
    "fmt"
    
    "github.com/nemohan/wireread"
)

func main() {
    // Assume we have a complete, validated message frame
    completeFrame := []byte{0x00, 0x0A, 0x48, 0x65, 0x6C, 0x6C, 0x6F}
    reader := wireread.NewFastReader(completeFrame)
    
    // No error checking needed - data is guaranteed to be complete
    length, _ := reader.ReadUint16BE()
    message, _ := reader.ReadString(int(length))
    
    fmt.Printf("Message: %s\n", message) // Output: Message: Hello
}
```

## API Overview

### Reader Interface

Both `SafeReader` and `FastReader` implement the `Reader` interface:

```go
type Reader interface {
    // Basic operations
    Bytes() []byte
    ReadBytes(n int) ([]byte, error)
    ReadByte() (byte, error)
    Skip(n int) error
    
    // String operations
    ReadString(n int) (string, error)
    ReadStringInto(out *string, n int) error
    ReadNullTerminatedString() (string, error)
    ReadLine() (string, error)
    
    // Variable-length integer
    ReadUvarint() (uint64, error)
    
    // Big-endian integers
    ReadUint16BE() (uint16, error)
    ReadUint32BE() (uint32, error)
    ReadUint64BE() (uint64, error)
    ReadUint16BEInto(out *uint16) error
    ReadUint32BEInto(out *uint32) error
    ReadUint64BEInto(out *uint64) error
    ReadInt16BEInto(out *int16) error
    ReadInt32BEInto(out *int32) error
    
    // Little-endian integers
    ReadUint16LE() (uint16, error)
    ReadUint32LE() (uint32, error)
    ReadUint64LE() (uint64, error)
    ReadUint16LEInto(out *uint16) error
    ReadUint32LEInto(out *uint32) error
    ReadUint64LEInto(out *uint64) error
    
    // Protocol-specific
    ReadLengthEncodedInteger() (uint64, error) // MySQL format
}
```

## Usage Examples

### Parsing a Binary Protocol

```go
func parseMessage(data []byte) error {
    reader := wireread.NewSafeReader(data)
    
    // Read message header
    var msgType uint16
    if err := reader.ReadUint16BEInto(&msgType); err != nil {
        return err
    }
    
    // Read payload length
    payloadLen, err := reader.ReadUint32BE()
    if err != nil {
        return err
    }
    
    // Read payload
    payload, err := reader.ReadBytes(int(payloadLen))
    if err != nil {
        return err
    }
    
    fmt.Printf("Type: %d, Length: %d, Payload: %x\n", 
        msgType, payloadLen, payload)
    return nil
}
```

### Reading Different Data Types

```go
reader := wireread.NewSafeReader(data)

// Read individual bytes
b, _ := reader.ReadByte()

// Read byte slices
bytes, _ := reader.ReadBytes(10)

// Read fixed-length strings
str, _ := reader.ReadString(5)

// Read null-terminated strings (C-style)
cstr, _ := reader.ReadNullTerminatedString()

// Read line-delimited data
line, _ := reader.ReadLine()

// Read variable-length integers
varint, _ := reader.ReadUvarint()

// Read MySQL length-encoded integers
mysqlInt, _ := reader.ReadLengthEncodedInteger()

// Skip unwanted data
_ = reader.Skip(4)

// Get remaining data
remaining := reader.Bytes()
```

### Using Pointer-Based Methods

```go
reader := wireread.NewSafeReader(data)

var (
    id      uint32
    version uint16
    flags   uint8
)

// Read into existing variables
reader.ReadUint32BEInto(&id)
reader.ReadUint16BEInto(&version)
reader.ReadByte() // for flags

fmt.Printf("ID: %d, Version: %d\n", id, version)
```

## Performance Comparison

| Operation      | SafeReader | FastReader | Speedup |
| -------------- | ---------- | ---------- | ------- |
| ReadUint32BE   | ~8 ns/op   | ~3 ns/op   | ~2.7x   |
| ReadUint64LE   | ~9 ns/op   | ~3 ns/op   | ~3x     |
| ReadBytes(100) | ~45 ns/op  | ~25 ns/op  | ~1.8x   |

*Benchmarks run on Go 1.22, AMD64 architecture*

## When to Use Which Reader

### Use SafeReader when:
- ✅ Parsing untrusted or external data
- ✅ Data length is unknown or variable
- ✅ You need detailed error information
- ✅ Debugging protocol implementations
- ✅ Safety is more important than performance

### Use FastReader when:
- ✅ Data has been pre-validated
- ✅ Processing complete message frames
- ✅ Performance is critical (hot paths)
- ✅ Data integrity is guaranteed by upper layers
- ✅ Willing to accept panics on invalid data

## Protocol-Specific Features

### MySQL Length-Encoded Integers

```go
reader := wireread.NewSafeReader(data)
value, err := reader.ReadLengthEncodedInteger()
```

Supports MySQL's variable-length integer encoding:
- `< 0xFB`: 1-byte integer
- `0xFC`: 2-byte integer follows
- `0xFD`: 3-byte integer follows
- `0xFE`: 8-byte integer follows
- `0xFB`: NULL value (returns 0)

## Error Handling

`SafeReader` returns `io.ErrUnexpectedEOF` when there's insufficient data:

```go
reader := wireread.NewSafeReader([]byte{0x01})
value, err := reader.ReadUint32BE()
if err == io.ErrUnexpectedEOF {
    fmt.Println("Not enough data")
}
```

`FastReader` will panic with index out of bounds if data is insufficient.

## Thread Safety

Neither `SafeReader` nor `FastReader` is thread-safe. Each goroutine should have its own reader instance.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Acknowledgments

Inspired by the need for efficient binary protocol parsing in network applications.
