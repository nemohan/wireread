package wireread

import (
	"testing"
)

// Benchmark SafeReader operations
func BenchmarkSafeReader_ReadUint16BE(b *testing.B) {
	data := make([]byte, b.N*2)
	for i := 0; i < len(data); i += 2 {
		data[i] = 0x01
		data[i+1] = 0x02
	}
	r := NewSafeReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReadUint16BE()
	}
}

func BenchmarkSafeReader_ReadUint32BE(b *testing.B) {
	data := make([]byte, b.N*4)
	r := NewSafeReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReadUint32BE()
	}
}

func BenchmarkSafeReader_ReadUint64BE(b *testing.B) {
	data := make([]byte, b.N*8)
	r := NewSafeReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReadUint64BE()
	}
}

func BenchmarkSafeReader_ReadUint32LE(b *testing.B) {
	data := make([]byte, b.N*4)
	r := NewSafeReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReadUint32LE()
	}
}

func BenchmarkSafeReader_ReadBytes(b *testing.B) {
	data := make([]byte, b.N*100)
	r := NewSafeReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReadBytes(100)
	}
}

func BenchmarkSafeReader_ReadString(b *testing.B) {
	data := make([]byte, b.N*50)
	r := NewSafeReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReadString(50)
	}
}

// Benchmark FastReader operations
func BenchmarkFastReader_ReadUint16BE(b *testing.B) {
	data := make([]byte, b.N*2)
	for i := 0; i < len(data); i += 2 {
		data[i] = 0x01
		data[i+1] = 0x02
	}
	r := NewFastReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReadUint16BE()
	}
}

func BenchmarkFastReader_ReadUint32BE(b *testing.B) {
	data := make([]byte, b.N*4)
	r := NewFastReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReadUint32BE()
	}
}

func BenchmarkFastReader_ReadUint64BE(b *testing.B) {
	data := make([]byte, b.N*8)
	r := NewFastReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReadUint64BE()
	}
}

func BenchmarkFastReader_ReadUint32LE(b *testing.B) {
	data := make([]byte, b.N*4)
	r := NewFastReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReadUint32LE()
	}
}

func BenchmarkFastReader_ReadBytes(b *testing.B) {
	data := make([]byte, b.N*100)
	r := NewFastReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReadBytes(100)
	}
}

func BenchmarkFastReader_ReadString(b *testing.B) {
	data := make([]byte, b.N*50)
	r := NewFastReader(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReadString(50)
	}
}

// Comparison benchmarks
func BenchmarkComparison_Uint32BE(b *testing.B) {
	data := make([]byte, 1000000)

	b.Run("SafeReader", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := NewSafeReader(data)
			for j := 0; j < 1000; j++ {
				r.ReadUint32BE()
			}
		}
	})

	b.Run("FastReader", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := NewFastReader(data)
			for j := 0; j < 1000; j++ {
				r.ReadUint32BE()
			}
		}
	})
}

func BenchmarkComparison_MixedOperations(b *testing.B) {
	data := make([]byte, 1000000)

	b.Run("SafeReader", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := NewSafeReader(data)
			for j := 0; j < 100; j++ {
				r.ReadUint16BE()
				r.ReadUint32LE()
				r.ReadBytes(10)
				r.ReadUint64BE()
			}
		}
	})

	b.Run("FastReader", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := NewFastReader(data)
			for j := 0; j < 100; j++ {
				r.ReadUint16BE()
				r.ReadUint32LE()
				r.ReadBytes(10)
				r.ReadUint64BE()
			}
		}
	})
}
