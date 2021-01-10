package aliyun

import (
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/micro-kit/micro-common/config"
	"github.com/micro-kit/micro-common/etcdcli"
	"github.com/micro-kit/micro-common/logger"
)

/* 阿里云相关操作客户端，短信、oss、推送、mq等 */

var (
	AliyunConfig *config.AliyunConfig
	AliyunSdkClient *sdk.Client
)

// 初始化redis
func init() {
	var err error
	err = config.GetAliyunConfig(etcdcli.EtcdCli, func(cfg *config.AliyunConfig) {
		AliyunConfig = cfg
		// 初始化阿里云sdk客户端
		AliyunSdkClient,err = NewAliyunClient();
	})
	if err != nil {
		logger.Logger.Panicw("Get aliyun configuration error", "err", err)
	}
}

// NewAliyunClient 创建阿里云客户端
func NewAliyunClient() (*sdk.Client, error) {
	if AliyunConfig == nil {
		return nil, errors.New("阿里云配置为nil")
	}
	return sdk.NewClientWithAccessKey(AliyunConfig.RegionId, AliyunConfig.SecretId, AliyunConfig.SecretKey)
}

