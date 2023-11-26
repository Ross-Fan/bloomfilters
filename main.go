package main

import (
	"bloomfilters/bf"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	rddr = "127.0.0.1:6379"
	port = 6379
	DB   = 0
)

func main() {

	rdb := redis.NewClient(
		&redis.Options{
			Addr:     rddr,
			DB:       DB,
			PoolSize: 2,
		},
	)

	for i := 1; i < 3010; i++ {
		bmt, err := bf.GetBf(rdb, "115")
		if err != nil {
			fmt.Println("GetBf Err ... ", err.Error())
		}
		bf.SetBf(rdb, bmt, "115", strconv.Itoa(i))
	}

	cnt := 0
	start := time.Now().UnixMilli()
	bmt, err := bf.GetBf(rdb, "115")
	if err != nil {
		fmt.Println("GetBf Err ... ", err.Error())
	}
	for i := 10000; i < 20000; i++ {
		flag := bf.CheckBf(bmt, strconv.Itoa(i))
		if flag {
			cnt += 1
		}
	}
	fmt.Println("time: ", time.Now().UnixMilli()-start)
	fmt.Println("false positives count:", cnt)

	time.Sleep(1 * time.Second)
	rdb.Close()
}
