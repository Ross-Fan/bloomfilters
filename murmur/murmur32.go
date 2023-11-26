package murmur

import "unsafe"

const (
	c1 uint32 = 0xcc9e2d51
	c2 uint32 = 0x1b873593
	c3 uint32 = 0x85ebca6b
	c4 uint32 = 0xc2b2ae35
	r1 uint32 = 15
	r2 uint32 = 13
	m  uint32 = 5
	n  uint32 = 0xe6546b64
)

func Murmur32(data []byte, seed uint32) (h uint32) {
	h = seed
	length := len(data)
	pos := 0

	for length >= 4 {
		k := *(*uint32)(unsafe.Pointer(&data[pos]))
		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2

		h ^= k
		h = (h << r2) | (h >> (32 - r2))
		h = h*m + n
		length -= 4
		pos += 4
	}
	if length > 0 {
		var k uint32 = 0
		for i := 0; i < length; i++ {
			k |= uint32(data[pos+i]) << (8 * i)
		}
		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2
		h ^= k
	}
	h ^= uint32(len(data))
	h ^= h >> 16
	h *= c3
	h ^= h >> 13
	h *= c4
	h ^= h >> 16

	return h
}
