package redistool

import (
	"context"
	"net"
	"strconv"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/extconfig"
	"github.com/danenmao/pterergate-dtf/internal/config"
)

// go-redis连接池
var gs_RedisClient *goredis.Client

// 获取go-redis实例
func DefaultRedis() *goredis.Client {
	return gs_RedisClient
}

// 连接默认的MySQL数据库
func ConnectToDefaultRedis() {
	gs_RedisClient = InitRedisClient(&config.DefaultRedisServer)
}

// 初始化go-redis连接池
func InitRedisClient(cfg *extconfig.RedisAddress) *goredis.Client {
	dbNo, err := strconv.Atoi(cfg.DB)
	if err != nil {
		panic(err.Error())
	}

	client := goredis.NewClient(&goredis.Options{
		// 连接信息
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       dbNo,

		// 连接池的容量
		PoolSize:     20,
		MinIdleConns: 10,

		// 超时值
		DialTimeout:  5 * time.Second, // 连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, // 读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, // 写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, // 当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		// 闲置连接检查
		IdleCheckFrequency: 60 * time.Second, // 闲置连接检查的周期
		IdleTimeout:        5 * time.Minute,  // 闲置超时
		MaxConnAge:         0 * time.Second,  // 连接存活时长

		// 命令执行失败时的重试策略
		MaxRetries:      0,
		MinRetryBackoff: 8 * time.Millisecond,   // 每次计算重试间隔时间的下限
		MaxRetryBackoff: 512 * time.Millisecond, // 每次计算重试间隔时间的上限

		// 连接函数
		Dialer: func(ctx context.Context, network string, addr string) (net.Conn, error) {
			netDialer := &net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Minute,
			}

			return netDialer.DialContext(ctx, cfg.Type, addr)
		},

		// 钩子函数
		// 当客户端执行命令时需要从连接池获取连接，且连接池需要新建连接时, 会调用此钩子函数
		OnConnect: func(ctx context.Context, conn *goredis.Conn) error {
			return nil
		},
	})

	// 激活连接
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		glog.Fatal("failed to connect to redis", err)
		return nil
	}

	return client
}
