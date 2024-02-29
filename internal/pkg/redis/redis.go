package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"notification-service/configs"

	"github.com/go-redis/redis/v8"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
)

// var redisClient *redis.Client

type RedisClient struct {
	Client interface{}
}

func InitConnection(redisDB, redisHost, redisPort, redisPassword string, appConfig string) Collections {
	var client interface{}

	if appConfig != "cluster" {
		// Create Redis Client
		db := 0
		parseRedisDb, err := strconv.ParseInt(redisDB, 10, 32)

		if err == nil {
			db = int(parseRedisDb)
		}

		c := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%v:%v", redisHost, redisPort),
			Password: redisPassword,
			DB:       db,
		})

		if configs.GetConfig().Datadog.DatadogEnabled == "true" {
			redistrace.WrapClient(c)
		}

		if c.Ping(context.Background()).Err() != nil {
			panic("cannot connect redis")
		}
		client = c
	} else {
		// Create Redis Cluster Client
		hostArray := strings.Split(redisHost, ",")
		c := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    hostArray,
			Password: redisPassword,
		})

		// Test Connection
		for _, addr := range hostArray {
			nodeClient := redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: redisPassword,
			})

			_, err := nodeClient.Ping(context.Background()).Result()
			if err != nil {
				fmt.Println("REDIS CLUSTER ERROR:", err.Error())
				panic("Cannot connect to Redis Cluster")
			}
			nodeClient.Close()
		}
		client = c
	}
	return &RedisClient{Client: client}
}

type Collections interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Conn(ctx context.Context) *redis.Conn
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd

	Close() error
}

func (r *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return r.Client.(*redis.Client).SetNX(ctx, key, value, expiration)
}

func (r *RedisClient) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return r.Client.(*redis.Client).EvalSha(ctx, sha1, keys, args...)
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.Client.(*redis.Client).Del(ctx, keys...)
}

func (r *RedisClient) Conn(ctx context.Context) *redis.Conn {
	return r.Client.(*redis.Client).Conn(ctx)
}

func (r *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.Client.(*redis.Client).Get(ctx, key)
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.Client.(*redis.Client).Set(ctx, key, value, expiration)
}

func (r *RedisClient) Close() error {
	switch c := r.Client.(type) {
	case *redis.Client:
		return c.Close()
	case *redis.ClusterClient:
		return c.Close()
	default:
		return fmt.Errorf("unsupported Redis client type")
	}
}
