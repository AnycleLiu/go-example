package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: []string{"192.168.70.131:26379"},
	})

	defer rdb.Close()

	do := func() {
		_, err := rdb.HGetAll(ctx, "hopex:usersystem:user.UDZU0ODZER7TK9TYYLTP.ios").Result()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	N := 10000
	n := 0

	do()

	t := time.Now()

	for {
		n = 0
		t = time.Now()

		for time.Since(t).Milliseconds() < 1000 {
			do()

			n++
		}
		fmt.Printf("%d, 完成%d\n", N, n)
		if n < N {
			break
		}
		N = N * 2
	}

}
