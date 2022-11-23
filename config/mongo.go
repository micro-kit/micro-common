package config

import (
	"context"
	"errors"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/naoina/toml"
	"go.etcd.io/etcd/clientv3"
)

// MongoConfig mongo配置
type MongoConfig struct {
	Address     []string `toml:"address"`       // mongo 服务器地址,包括地址和端口 127.0.0.1:6379
	Db          string   `toml:"db"`            // mongo db
	MaxPoolSize int      `toml:"max_pool_size"` // 连接池最大值
	MinPoolSize int      `toml:"min_pool_size"` // 连接池最小值
	ReplicaSet  string   `toml:"replica_set"`   // ReplicaSet
	Username    string   `toml:"username"`      // 用户名
	Password    string   `toml:"password"`      // 密码
	AuthSource  string   `toml:"auth_source"`   // 验证数据库
}

// GetMongoConfig 获取mongo配置
func GetMongoConfig(cli *clientv3.Client, updateConfig func(*MongoConfig)) error {
	key := "root/config/" + GetSvcName() + "/mongo/cfg.toml"
	// fmt.Println(key)
	// 读etcd配置
	etcdResp, err := cli.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	if len(etcdResp.Kvs) == 0 {
		return errors.New("The mongo configuration configured to be empty. key: " + key)
	}
	mongoConfig := new(MongoConfig)
	err = toml.Unmarshal(etcdResp.Kvs[0].Value, mongoConfig)
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
				err = toml.Unmarshal(ev.Kv.Value, mongoConfig)
				if err != nil {
					continue
				}
				updateConfig(mongoConfig)
			}
		}
	}()

	// 调用更新配置回调函数
	updateConfig(mongoConfig)
	return nil
}
