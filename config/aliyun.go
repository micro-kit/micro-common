package config

import (
	"context"
	"errors"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/naoina/toml"
	"go.etcd.io/etcd/clientv3"
)

/* 阿里云相关操作配置 */

// AliyunConfig 阿里云sdk配置
type AliyunConfig struct {
	RegionId  string `toml:"region_id"`  // 阿里云接入区域id default | cn-hangzhou
	SecretId  string `toml:"secret_id"`  // 阿里云 SecretId
	SecretKey string `toml:"secret_key"` // 阿里云 SecretKey
}

// GetAliyunConfig 获取aliyun配置
func GetAliyunConfig(cli *clientv3.Client, updateConfig func(*AliyunConfig)) error {
	key := "root/config/" + GetSvcName() + "/aliyun/cfg.toml"
	// fmt.Println(key)
	// 读etcd配置
	etcdResp, err := cli.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	if len(etcdResp.Kvs) == 0 {
		return errors.New("The aliyun configuration configured to be empty. key: " + key)
	}
	aliyunConfig := new(AliyunConfig)
	err = toml.Unmarshal(etcdResp.Kvs[0].Value, aliyunConfig)
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
				err = toml.Unmarshal(ev.Kv.Value, aliyunConfig)
				if err != nil {
					continue
				}
				updateConfig(aliyunConfig)
			}
		}
	}()

	// 调用更新配置回调函数
	updateConfig(aliyunConfig)
	return nil
}
