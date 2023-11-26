package murmur

import (
	"strconv"
	"testing"
)

func BenchmarkMM32(b *testing.B) {
	buf := make([]byte, 8192)
	for length := 1; length <= cap(buf); length *= 2 {
		b.Run(
			strconv.Itoa(length),
			func(b *testing.B) {
				buf = buf[:length]
				b.SetBytes(int64(length))
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					Murmur32(buf, 0)
				}
			},
		)
	}
}
