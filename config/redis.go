package config

import (
	"context"
	"errors"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/naoina/toml"
	"go.etcd.io/etcd/clientv3"
)

// RedisConfig redis配置
type RedisConfig struct {
	Address   []string `toml:"address"`    // redis 服务器地址,包括地址和端口 127.0.0.1:6379
	Password  string   `toml:"password"`   // redis 密码
	Db        int      `toml:"db"`         // 连接的数据库
	PoolSize  int      `toml:"pool_size"`  // 连接池大小
	IsCluster bool     `toml:"is_cluster"` // 是否集群模式
}

// GetRedisConfig 获取redis配置
func GetRedisConfig(cli *clientv3.Client, updateConfig func(*RedisConfig)) error {
	key := "root/config/" + GetSvcName() + "/redis/cfg.toml"
	// fmt.Println(key)
	// 读etcd配置
	etcdResp, err := cli.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	if len(etcdResp.Kvs) == 0 {
		return errors.New("The redis configuration configured to be empty. key: " + key)
	}
	redisConfig := new(RedisConfig)
	err = toml.Unmarshal(etcdResp.Kvs[0].Value, redisConfig)
	if err != nil {
		return err
	}

	// 监视key变化
	go func() {
		rch := cli.Watch(context.Background(), key)
		for wresp := range rch {
			for _, ev := range wresp.Events {
				if string(ev.Kv.Key) != key || ev.Type != mvccpb.PUT {
					continue
				}
				err = toml.Unmarshal(ev.Kv.Value, redisConfig)
				if err != nil {
					continue
				}
				updateConfig(redisConfig)
			}
		}
	}()

	// 调用更新配置回调函数
	updateConfig(redisConfig)
	return nil
}
