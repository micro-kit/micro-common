package cache

import (
	"errors"
	"time"

	"github.com/micro-kit/micro-common/config"
	"github.com/micro-kit/micro-common/etcdcli"
	"github.com/micro-kit/micro-common/logger"
	goredis "gopkg.in/redis.v5"
)

/* 缓存连接 redis */

var (
	// Client redis 连接对象
	Client goredis.Cmdable
)

var (
	// NilErr 没有错误
	NilErr = goredis.Nil
)

// GetClient 获取redis连接对象
func GetClient() goredis.Cmdable {
	if Client == nil {
		initRedis()
	}
	return Client
}

func init() {
	initRedis()
}

// 初始化redis
func initRedis() {
	var err error
	err = config.GetRedisConfg(etcdcli.EtcdCli, func(cfg *config.RedisConfg) {
		Client, err = NewClient(cfg)
		if err != nil {
			logger.Logger.Panicw("Creating redis connection errors", "err", err)
		}
	})
	if err != nil {
		logger.Logger.Panicw("Get redis configuration error", "err", err)
	}
}

// NewClient 创建客户端连接
func NewClient(cfg *config.RedisConfg) (client goredis.Cmdable, err error) {
	if cfg == nil {
		err = errors.New("The redis configuration file can not be empty.")
		return
	}
	logger.Logger.Infow("Start connecting to redis database")
	if cfg.IsCluster == true {
		// redis集群
		client = goredis.NewClusterClient(&goredis.ClusterOptions{
			Addrs:    cfg.Address,
			Password: cfg.Password,
			PoolSize: cfg.PoolSize,
		})
	} else {
		// redis单机
		client = goredis.NewClient(&goredis.Options{
			Addr:     cfg.Address[0],
			Password: cfg.Password,
			DB:       cfg.Db,
			PoolSize: cfg.PoolSize,
		})
	}
	// ping 防止断开
	go func() {
		for {
			err := client.Ping().Err()
			if err != nil {
				logger.Logger.Errorw("redis ping error", "err", err)
			}

			time.Sleep(time.Second * 30)
		}
	}()
	logger.Logger.Infow("Connect to redis database successfully")
	return
}
