package bf

import (
	"bloomfilters/murmur"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

const (
	n = "n"
	b = "b"
	k = 2
	f = 4
	l = 3000
	s = 40000
)

// var Rdb = redis.NewClient(
// 	&redis.Options{
// 		Addr:     rddr,
// 		DB:       DB,
// 		PoolSize: 2,
// 	},
// )

type (
	bloomFilter struct {
		n      int // input number
		bitSlc [s]byte
	}

	BloomTuple struct {
		Tup [k]*bloomFilter
	}
)

type (
	UserStr struct {
		Uid string
	}
)

func NewBloomFilter() *bloomFilter {
	b := new(bloomFilter)
	b.n = 0
	b.bitSlc = [s]byte{0}
	return b
}

func GetBf(rds *redis.Client, id string) (bt *BloomTuple, err error) {
	res, err := rds.HGetAll(id).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	bt = new(BloomTuple)
	for i := 0; i < k; i++ {
		bt.Tup[i] = NewBloomFilter()
	}
	if res != nil {
		for i := 0; i < k; i++ {
			keyn := n + strconv.Itoa(i)
			keyb := b + strconv.Itoa(i)
			n_, ok1 := res[keyn]
			b_, ok2 := res[keyb]
			if ok1 && ok2 {
				// var err error
				bt.Tup[i].n, err = strconv.Atoi(n_)
				bs_ := []byte(b_)
				if len(bs_) <= s {
					copy(bt.Tup[i].bitSlc[:], []byte(b_)[:])
				} else {
					copy(bt.Tup[i].bitSlc[:], []byte(b_)[:s])
				}
				if err != nil {
					bt.Tup[i].n = 0
					bt.Tup[i].bitSlc = [s]byte{0}
				}
			}
		}
	}
	return bt, nil
}

func SetBf(rds *redis.Client, bt *BloomTuple, id string, item string) {
	h := getHashed(item)

	j := 0
	for i := 0; i < k; i++ {
		if bt.Tup[i].n >= l {
			j += 1
			continue
		}
	}
	if j >= k { // obsolete the first tuple
		for i := 1; i < k; i++ {
			bt.Tup[i-1] = bt.Tup[i]
		}
		bt.Tup[k-1] = NewBloomFilter()
	}

	for i := 0; i < k; i++ {
		if bt.Tup[i].n < l {
			for p := 0; p < f; p++ {
				idx := (h[p] % s) / 8
				offset := h[p] & 0x07
				bt.Tup[i].bitSlc[idx] |= (1 << offset)
			}
			bt.Tup[i].n += 1
		}

	}

	// set to redis
	for i := 0; i < k; i++ {
		keyn := n + strconv.Itoa(i)
		keyb := b + strconv.Itoa(i)
		_, err1 := rds.HSet(id, keyn, bt.Tup[i].n).Result()
		if err1 != nil {
			fmt.Printf("HSet Failure ... %s, %s \n", id, keyn)
		}

		_, err2 := rds.HSet(id, keyb, string(bt.Tup[i].bitSlc[:])).Result()
		if err2 != nil {
			fmt.Printf("HSet Failure ... %s, %s, %s \n", id, keyb, err2)
		}
	}

}

// true -- item in the Bloom,
// false -- item not in the Bloom
func CheckBf(bt *BloomTuple, item string) bool {

	h := getHashed(item)

	for i := 0; i < k; i++ {
		bFlag := 0
		for p := 0; p < f; p++ {
			idx := (h[p] % s) / 8
			offset := h[p] & 0x07

			if (bt.Tup[i].bitSlc[idx] & (1 << offset)) == (1 << offset) {

				bFlag += 1
			}
		}
		// fmt.Println("bFlag:", bFlag)
		if bFlag >= f {
			return true
		}
	}

	return false
}

func getHashed(item string) []uint32 {
	rst := make([]uint32, 0, f)
	item_byte := []byte(item)
	for i := 0; i < f; i++ {
		h := murmur.Murmur32(item_byte, 0)
		rst = append(rst, h)
		if i == f-1 {
			break
		}
		item_byte = Uint32toByte(h)
	}
	return rst
}

func Uint32toByte(x uint32) []byte {
	b := [4]byte{0}
	b[0] = uint8(x & 0xff)
	b[1] = uint8((x >> 8) & 0xff)
	b[2] = uint8((x >> 16) & 0xff)
	b[3] = uint8((x >> 24) & 0xff)
	return b[:]
}
