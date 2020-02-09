package config

import (
	"context"
	"errors"

	"github.com/naoina/toml"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

// RedisConfg redis配置
type RedisConfg struct {
	Address   []string `toml:"address"`    // redis 服务器地址,包括地址和端口 127.0.0.1:6379
	Password  string   `toml:"password"`   // redis 密码
	Db        int      `toml:"db"`         // 连接的数据库
	PoolSize  int      `toml:"pool_size"`  // 连接池大小
	IsCluster bool     `toml:"is_cluster"` // 是否集群模式
}

// GetRedisConfg 获取redis配置
func GetRedisConfg(cli *clientv3.Client, updateConfig func(*RedisConfg)) error {
	key := "root/config/" + GetSvcName() + "/redis/cfg.toml"
	// fmt.Println(key)
	// 初始化master连接
	etcdResp, err := cli.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	if len(etcdResp.Kvs) == 0 {
		return errors.New("The redis configuration configured to be empty. key: " + key)
	}
	redisConfg := new(RedisConfg)
	err = toml.Unmarshal(etcdResp.Kvs[0].Value, redisConfg)
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
				err = toml.Unmarshal(ev.Kv.Value, redisConfg)
				if err != nil {
					continue
				}
				updateConfig(redisConfg)
			}
		}
	}()

	// 调用更新配置回调函数
	updateConfig(redisConfg)
	return nil
}
