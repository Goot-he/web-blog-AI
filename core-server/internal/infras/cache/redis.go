package cache

import (
	"context"
	"core-server/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	Cache *redis.Client
}

func NewCacheClient(cfg config.Config) *CacheClient {
	options := redis.Options{
		Addr:     cfg.RedisConf.Host + ":" + cfg.RedisConf.Port,
		Username: cfg.RedisConf.Username,
		Password: cfg.RedisConf.Password,
		DB:       cfg.RedisConf.DB,

		// 设置连接池以及闲置连接数
		// 设置最小剩余空闲连接数 保证当前剩余不小于它
		MinIdleConns: 10,

		// 设置超时时间
		DialTimeout:  5 * time.Second, // 连接超时 默认五秒
		ReadTimeout:  3 * time.Second, // 读超时 默认三秒 -1表示取消超时
		WriteTimeout: 3 * time.Second, // 写超时

		// 命令执行失败的重试策略
		MaxRetries:      3,                      // 最多重试次数 默认为0
		MinRetryBackoff: 8 * time.Millisecond,   // 重试间隔最小值 设为-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, // 重试间隔最大值
	}

	rdb := redis.NewClient(&options)

	// 测试链接
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		//log.Errorf("redis connect err: %v", err)
	}

	// 连接成功

	return &CacheClient{
		Cache: rdb,
	}
}

func (c *CacheClient) Close() error {
	return c.Cache.Close()
}
