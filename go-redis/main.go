package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "39.101.244.245:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})

	// String
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	// 列表 List
	rdb.LPush(ctx, "list", 1, 2, 3)

	// 集合
	rdb.SAdd(ctx, "team", "kobe", "jordan")
	rdb.SAdd(ctx, "team", "curry")
	rdb.SAdd(ctx, "team", "kobe")

	// hash
	rdb.HSet(ctx, "user", "key1", "value1", "key2", "value2")
	rdb.HSet(ctx, "user", []string{"key3", "value3", "key4", "value4"})
	rdb.HSet(ctx, "user", map[string]interface{}{"key5": "value5", "key6": "value6"})

	// 有序集合
	rdb.ZAdd(ctx, "zSet", &redis.Z{
		Score:  0,
		Member: 1,
	})
	rdb.ZAdd(ctx, "zSet", &redis.Z{
		Score:  0,
		Member: 2,
	})
	rdb.ZAdd(ctx, "zSet", &redis.Z{
		Score:  0,
		Member: 3,
	})
	// val, err := rdb.Get(ctx, "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("key", val)

	// val2, err := rdb.Get(ctx, "key2").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key2 does not exist")
	// } else if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println("key2", val2)
	// }
}
