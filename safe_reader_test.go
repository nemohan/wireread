package wireread

import (
	"io"
	"testing"
)

func TestSafeReader_ReadByte(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    byte
		wantErr bool
	}{
		{"read single byte", []byte{0x42}, 0x42, false},
		{"empty data", []byte{}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSafeReader(tt.data)
			got, err := r.ReadByte()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeReader_ReadBytes(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		n       int
		want    []byte
		wantErr bool
	}{
		{"read 4 bytes", []byte{1, 2, 3, 4, 5}, 4, []byte{1, 2, 3, 4}, false},
		{"insufficient data", []byte{1, 2}, 5, nil, true},
		{"read 0 bytes", []byte{1, 2}, 0, []byte{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSafeReader(tt.data)
			got, err := r.ReadBytes(tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !bytesEqual(got, tt.want) {
				t.Errorf("ReadBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeReader_ReadUint16BE(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    uint16
		wantErr bool
	}{
		{"valid uint16", []byte{0x01, 0x02}, 0x0102, false},
		{"max uint16", []byte{0xFF, 0xFF}, 0xFFFF, false},
		{"insufficient data", []byte{0x01}, 0, true},
		{"empty data", []byte{}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSafeReader(tt.data)
			got, err := r.ReadUint16BE()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUint16BE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUint16BE() = 0x%04x, want 0x%04x", got, tt.want)
			}
		})
	}
}

func TestSafeReader_ReadUint32BE(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    uint32
		wantErr bool
	}{
		{"valid uint32", []byte{0x01, 0x02, 0x03, 0x04}, 0x01020304, false},
		{"max uint32", []byte{0xFF, 0xFF, 0xFF, 0xFF}, 0xFFFFFFFF, false},
		{"insufficient data", []byte{0x01, 0x02}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSafeReader(tt.data)
			got, err := r.ReadUint32BE()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUint32BE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUint32BE() = 0x%08x, want 0x%08x", got, tt.want)
			}
		})
	}
}

func TestSafeReader_ReadUint64BE(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    uint64
		wantErr bool
	}{
		{"valid uint64", []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, 0x0102030405060708, false},
		{"insufficient data", []byte{0x01, 0x02, 0x03, 0x04}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSafeReader(tt.data)
			got, err := r.ReadUint64BE()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUint64BE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUint64BE() = 0x%016x, want 0x%016x", got, tt.want)
			}
		})
	}
}

func TestSafeReader_ReadUint16LE(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    uint16
		wantErr bool
	}{
		{"valid uint16", []byte{0x02, 0x01}, 0x0102, false},
		{"max uint16", []byte{0xFF, 0xFF}, 0xFFFF, false},
		{"insufficient data", []byte{0x01}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSafeReader(tt.data)
			got, err := r.ReadUint16LE()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUint16LE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUint16LE() = 0x%04x, want 0x%04x", got, tt.want)
			}
		})
	}
}

func TestSafeReader_ReadUint32LE(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    uint32
		wantErr bool
	}{
		{"valid uint32", []byte{0x04, 0x03, 0x02, 0x01}, 0x01020304, false},
		{"insufficient data", []byte{0x01, 0x02}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSafeReader(tt.data)
			got, err := r.ReadUint32LE()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUint32LE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUint32LE() = 0x%08x, want 0x%08x", got, tt.want)
			}
		})
	}
}

func TestSafeReader_ReadString(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		n       int
		want    string
		wantErr bool
	}{
		{"read hello", []byte("Hello World"), 5, "Hello", false},
		{"read empty", []byte("test"), 0, "", false},
		{"insufficient data", []byte("Hi"), 5, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSafeReader(tt.data)
			got, err := r.ReadString(tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeReader_ReadNullTerminatedString(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    string
		wantErr bool
	}{
		{"simple string", []byte{'H', 'e', 'l', 'l', 'o', 0}, "Hello", false},
		{"empty string", []byte{0}, "", false},
		{"no null terminator", []byte{'H', 'i'}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSafeReader(tt.data)
			got, err := r.ReadNullTerminatedString()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadNullTerminatedString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadNullTerminatedString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeReader_ReadLine(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    string
		wantErr bool
	}{
		{"unix line", []byte("Hello\n"), "Hello", false},
		{"windows line", []byte("Hello\r\n"), "Hello", false},
		{"no newline", []byte("Hello"), "", true},
		{"empty line", []byte("\n"), "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSafeReader(tt.data)
			got, err := r.ReadLine()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadLine() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSafeReader_ReadLengthEncodedInteger(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    uint64
		wantErr bool
	}{
		{"1-byte value", []byte{0x05}, 5, false},
		{"null value", []byte{0xFB}, 0, false},
		{"2-byte value", []byte{0xFC, 0x01, 0x02}, 0x0201, false},
		{"3-byte value", []byte{0xFD, 0x01, 0x02, 0x03}, 0x030201, false},
		{"8-byte value", []byte{0xFE, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, 0x0807060504030201, false},
		{"insufficient 2-byte", []byte{0xFC, 0x01}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSafeReader(tt.data)
			got, err := r.ReadLengthEncodedInteger()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLengthEncodedInteger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadLengthEncodedInteger() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSafeReader_Skip(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	r := NewSafeReader(data)

	if err := r.Skip(2); err != nil {
		t.Errorf("Skip(2) error = %v", err)
	}

	got, _ := r.ReadByte()
	if got != 3 {
		t.Errorf("After Skip(2), ReadByte() = %d, want 3", got)
	}

	if err := r.Skip(10); err != io.ErrUnexpectedEOF {
		t.Errorf("Skip(10) error = %v, want io.ErrUnexpectedEOF", err)
	}
}

func TestSafeReader_Bytes(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	r := NewSafeReader(data)

	r.ReadBytes(2)
	remaining := r.Bytes()

	want := []byte{3, 4, 5}
	if !bytesEqual(remaining, want) {
		t.Errorf("Bytes() = %v, want %v", remaining, want)
	}
}

func TestSafeReader_MultipleReads(t *testing.T) {
	data := []byte{
		0x00, 0x01, // uint16 BE = 1
		0x02, 0x00, 0x00, 0x00, // uint32 LE = 2
		'H', 'i', 0, // null-terminated string
	}

	r := NewSafeReader(data)

	v16, err := r.ReadUint16BE()
	if err != nil || v16 != 1 {
		t.Errorf("ReadUint16BE() = %d, %v; want 1, nil", v16, err)
	}

	v32, err := r.ReadUint32LE()
	if err != nil || v32 != 2 {
		t.Errorf("ReadUint32LE() = %d, %v; want 2, nil", v32, err)
	}

	str, err := r.ReadNullTerminatedString()
	if err != nil || str != "Hi" {
		t.Errorf("ReadNullTerminatedString() = %q, %v; want \"Hi\", nil", str, err)
	}
}

func TestSafeReader_Into_Methods(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	r := NewSafeReader(data)

	var u16 uint16
	if err := r.ReadUint16BEInto(&u16); err != nil || u16 != 0x0001 {
		t.Errorf("ReadUint16BEInto() = %d, %v; want 1, nil", u16, err)
	}

	var i16 int16
	if err := r.ReadInt16BEInto(&i16); err != nil || i16 != 0x0203 {
		t.Errorf("ReadInt16BEInto() = %d, %v; want 515, nil", i16, err)
	}

	var u32 uint32
	if err := r.ReadUint32BEInto(&u32); err != nil || u32 != 0x04050607 {
		t.Errorf("ReadUint32BEInto() = %d, %v; want 67438087, nil", u32, err)
	}
}

// Helper function
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
