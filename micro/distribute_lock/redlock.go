package main

import (
	"flag"
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"strings"
	"time"
)

var addrs = flag.String("addrs", "", "comma separated list of addrs")

func init() {
	flag.Parse()
}

// redis 实现分布式锁
// https://pandaychen.github.io/2020/06/01/REDIS-DISTRIBUTED-LOCK/
func main() {
	var poolList []redis.Pool
	var addrList []string
	if len(*addrs) == 0 {
		addrList = []string{"localhost:6379"}
	} else {
		addrList = strings.Split(*addrs, ",")
	}
	fmt.Println("addrList", addrList)

	for _, addr := range addrList {
		client := goredislib.NewClient(
			&goredislib.Options{Addr: addr},
		)
		pool := goredis.NewPool(client)
		poolList = append(poolList, pool)
	}

	// 每个pool就是一个redis实例
	rs := redsync.New(poolList...)

	mutexName := "my-global-mutex"
	// 设置超时时间
	mutex := rs.NewMutex(mutexName)

	err := mutex.Lock()
	if err != nil {
		panic(err)
	}

	// Do your work that requires the lock.
	doLogin()

	// 如果只有一台实例，quorum等于1
	// 释放时会判断 n < m.quorum，所以ok=false
	if ok, err := mutex.Unlock(); !ok || err != nil {
		fmt.Println("unlock", ok, err)
		panic("unlock failed:")
	}
}

func doLogin() {
	fmt.Println("do logic==========>")
	time.Sleep(2 * time.Second)
}
