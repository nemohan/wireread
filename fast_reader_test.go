package wireread

import (
	"testing"
)

func TestFastReader_ReadByte(t *testing.T) {
	data := []byte{0x42, 0x43, 0x44}
	r := NewFastReader(data)

	got, _ := r.ReadByte()
	if got != 0x42 {
		t.Errorf("ReadByte() = %v, want 0x42", got)
	}

	got, _ = r.ReadByte()
	if got != 0x43 {
		t.Errorf("ReadByte() = %v, want 0x43", got)
	}
}

func TestFastReader_ReadBytes(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	r := NewFastReader(data)

	got, _ := r.ReadBytes(3)
	want := []byte{1, 2, 3}

	if !bytesEqual(got, want) {
		t.Errorf("ReadBytes() = %v, want %v", got, want)
	}
}

func TestFastReader_ReadUint16BE(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want uint16
	}{
		{"simple", []byte{0x01, 0x02}, 0x0102},
		{"max", []byte{0xFF, 0xFF}, 0xFFFF},
		{"zero", []byte{0x00, 0x00}, 0x0000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFastReader(tt.data)
			got, _ := r.ReadUint16BE()
			if got != tt.want {
				t.Errorf("ReadUint16BE() = 0x%04x, want 0x%04x", got, tt.want)
			}
		})
	}
}

func TestFastReader_ReadUint32BE(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want uint32
	}{
		{"simple", []byte{0x01, 0x02, 0x03, 0x04}, 0x01020304},
		{"max", []byte{0xFF, 0xFF, 0xFF, 0xFF}, 0xFFFFFFFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFastReader(tt.data)
			got, _ := r.ReadUint32BE()
			if got != tt.want {
				t.Errorf("ReadUint32BE() = 0x%08x, want 0x%08x", got, tt.want)
			}
		})
	}
}

func TestFastReader_ReadUint64BE(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	r := NewFastReader(data)

	got, _ := r.ReadUint64BE()
	want := uint64(0x0102030405060708)

	if got != want {
		t.Errorf("ReadUint64BE() = 0x%016x, want 0x%016x", got, want)
	}
}

func TestFastReader_ReadUint16LE(t *testing.T) {
	data := []byte{0x02, 0x01}
	r := NewFastReader(data)

	got, _ := r.ReadUint16LE()
	want := uint16(0x0102)

	if got != want {
		t.Errorf("ReadUint16LE() = 0x%04x, want 0x%04x", got, want)
	}
}

func TestFastReader_ReadUint32LE(t *testing.T) {
	data := []byte{0x04, 0x03, 0x02, 0x01}
	r := NewFastReader(data)

	got, _ := r.ReadUint32LE()
	want := uint32(0x01020304)

	if got != want {
		t.Errorf("ReadUint32LE() = 0x%08x, want 0x%08x", got, want)
	}
}

func TestFastReader_ReadUint64LE(t *testing.T) {
	data := []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}
	r := NewFastReader(data)

	got, _ := r.ReadUint64LE()
	want := uint64(0x0102030405060708)

	if got != want {
		t.Errorf("ReadUint64LE() = 0x%016x, want 0x%016x", got, want)
	}
}

func TestFastReader_ReadString(t *testing.T) {
	data := []byte("Hello World")
	r := NewFastReader(data)

	got, _ := r.ReadString(5)
	want := "Hello"

	if got != want {
		t.Errorf("ReadString() = %v, want %v", got, want)
	}

	got, _ = r.ReadString(0)
	if got != "" {
		t.Errorf("ReadString(0) = %v, want empty string", got)
	}
}

func TestFastReader_ReadNullTerminatedString(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{"simple", []byte{'H', 'e', 'l', 'l', 'o', 0}, "Hello"},
		{"empty", []byte{0}, ""},
		{"with extra", []byte{'H', 'i', 0, 'x', 'y'}, "Hi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFastReader(tt.data)
			got, _ := r.ReadNullTerminatedString()
			if got != tt.want {
				t.Errorf("ReadNullTerminatedString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFastReader_ReadLine(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{"unix", []byte("Hello\n"), "Hello"},
		{"windows", []byte("Hello\r\n"), "Hello"},
		{"empty", []byte("\n"), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFastReader(tt.data)
			got, _ := r.ReadLine()
			if got != tt.want {
				t.Errorf("ReadLine() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFastReader_ReadLengthEncodedInteger(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want uint64
	}{
		{"1-byte", []byte{0x05}, 5},
		{"null", []byte{0xFB}, 0},
		{"2-byte", []byte{0xFC, 0x01, 0x02}, 0x0201},
		{"3-byte", []byte{0xFD, 0x01, 0x02, 0x03}, 0x030201},
		{"8-byte", []byte{0xFE, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, 0x0807060504030201},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFastReader(tt.data)
			got, _ := r.ReadLengthEncodedInteger()
			if got != tt.want {
				t.Errorf("ReadLengthEncodedInteger() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFastReader_Skip(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	r := NewFastReader(data)

	r.Skip(2)
	got, _ := r.ReadByte()

	if got != 3 {
		t.Errorf("After Skip(2), ReadByte() = %d, want 3", got)
	}
}

func TestFastReader_Bytes(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	r := NewFastReader(data)

	r.ReadBytes(2)
	remaining := r.Bytes()

	want := []byte{3, 4, 5}
	if !bytesEqual(remaining, want) {
		t.Errorf("Bytes() = %v, want %v", remaining, want)
	}
}

func TestFastReader_MultipleReads(t *testing.T) {
	data := []byte{
		0x00, 0x01, // uint16 BE = 1
		0x02, 0x00, 0x00, 0x00, // uint32 LE = 2
		'H', 'i', 0, // null-terminated string
	}

	r := NewFastReader(data)

	v16, _ := r.ReadUint16BE()
	if v16 != 1 {
		t.Errorf("ReadUint16BE() = %d, want 1", v16)
	}

	v32, _ := r.ReadUint32LE()
	if v32 != 2 {
		t.Errorf("ReadUint32LE() = %d, want 2", v32)
	}

	str, _ := r.ReadNullTerminatedString()
	if str != "Hi" {
		t.Errorf("ReadNullTerminatedString() = %q, want \"Hi\"", str)
	}
}

func TestFastReader_Into_Methods(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	r := NewFastReader(data)

	var u16 uint16
	r.ReadUint16BEInto(&u16)
	if u16 != 0x0001 {
		t.Errorf("ReadUint16BEInto() = %d, want 1", u16)
	}

	var i16 int16
	r.ReadInt16BEInto(&i16)
	if i16 != 0x0203 {
		t.Errorf("ReadInt16BEInto() = %d, want 515", i16)
	}

	var u32 uint32
	r.ReadUint32BEInto(&u32)
	if u32 != 0x04050607 {
		t.Errorf("ReadUint32BEInto() = %d, want 67438087", u32)
	}
}

func TestFastReader_StringInto(t *testing.T) {
	data := []byte("Hello World")
	r := NewFastReader(data)

	var str string
	r.ReadStringInto(&str, 5)
	if str != "Hello" {
		t.Errorf("ReadStringInto() = %v, want Hello", str)
	}
}

// Test that FastReader satisfies Reader interface
func TestFastReader_ImplementsReader(t *testing.T) {
	var _ Reader = (*FastReader)(nil)
}

// Test that SafeReader satisfies Reader interface
func TestSafeReader_ImplementsReader(t *testing.T) {
	var _ Reader = (*SafeReader)(nil)
}
