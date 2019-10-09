package utils

import (
	// "fmt"
	// "time"

	"github.com/go-redis/redis"
)

func InitRedisClient(rc Redis) (client *redis.Client) {
	Debug(NowFunc())

	client = redis.NewClient(&redis.Options{
		Network:    rc.Network,
		Addr:       rc.Addr,
		Password:   rc.Password,
		MaxRetries: rc.Max_retries,
		PoolSize:   rc.Pool_size,
	})
	err := client.Ping().Err()
	if err != nil {
		Panic(err)
	}

	return
}

func InitRedisCluster(rc RedisCluster) (clusterClient *redis.ClusterClient) {
	Debug(NowFunc())

	clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
		Password:   rc.Password,
		MaxRetries: rc.Max_retries,
		PoolSize:   rc.Pool_size,
		ClusterSlots: func() (slots []redis.ClusterSlot, err error) {
			len := len(rc.Master_addr_arr)
			gap := 16384 / len
			for i := 0; i < len; i++ {
				slots = append(slots, redis.ClusterSlot{
					Start: i * gap,
					End:   (i+1)*gap - 1,
					Nodes: []redis.ClusterNode{{
						Addr: rc.Master_addr_arr[i],
					}, {
						Addr: rc.Slave_addr_arr[i],
					}},
				})
			}
			return
		},
	})
	var err error
	err = clusterClient.Ping().Err()
	if err != nil {
		Panic(err)
	}
	err = clusterClient.ReloadState()
	if err != nil {
		Panic(err)
	}

	return
}
