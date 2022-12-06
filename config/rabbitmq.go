package config

import (
	"context"
	"errors"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/naoina/toml"
	"go.etcd.io/etcd/clientv3"
)

// RabbitmqConfig rabbitmq配置
type RabbitmqConfig struct {
	Host     string `toml:"host"`     // rabbitmq 服务器地址
	Port     int    `toml:"port"`     // rabbitmq 端口
	User     string `toml:"user"`     // 用户
	Password string `toml:"password"` // 密码
	Vhost    string `toml:"vhost"`    // mq vhost
}

// GetRabbitmqConfig 获取rabbitmq配置
func GetRabbitmqConfig(cli *clientv3.Client, updateConfig func(*RabbitmqConfig)) error {
	key := "root/config/" + GetSvcName() + "/rabbitmq/cfg.toml"
	// fmt.Println(key)
	// 读etcd配置
	etcdResp, err := cli.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	if len(etcdResp.Kvs) == 0 {
		return errors.New("The rabbitmq configuration configured to be empty. key: " + key)
	}
	rabbitmqConfig := new(RabbitmqConfig)
	err = toml.Unmarshal(etcdResp.Kvs[0].Value, rabbitmqConfig)
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
				err = toml.Unmarshal(ev.Kv.Value, rabbitmqConfig)
				if err != nil {
					continue
				}
				updateConfig(rabbitmqConfig)
			}
		}
	}()

	// 调用更新配置回调函数
	updateConfig(rabbitmqConfig)
	return nil
}
