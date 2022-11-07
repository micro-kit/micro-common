package config

import (
	"os"
	"strconv"
)

// 全局配置
const (
	// DEFAULT_ETCD_ADDR 默认etcd地址
	DEFAULT_ETCD_ADDR string = "127.0.0.1:2379"
	// DEFAULT_ETCD_USER 默认etcd用户名
	DEFAULT_ETCD_USER string = "root"
	// DEFAULT_ETCD_PASSWORD 默认etcd密码
	DEFAULT_ETCD_PASSWORD string = ""
	// DEFAULT_HTTP_ADDR 默认http监听地址
	DEFAULT_HTTP_ADDR string = "127.0.0.1:18080"
	// DEFAULT_GRPC_ADDR 默认grpc监听地址
	DEFAULT_GRPC_ADDR string = "127.0.0.1:28080"
	// DEFAULT_TCP_ADDR 默认tcp监听地址
	DEFAULT_TCP_ADDR string = "127.0.0.1:38080"
	// DEFAULT_SVC_NAME 默认服务名
	DEFAULT_SVC_NAME string = "default"
	// 当前运行环境，dev or pro or test
	DEFAULT_MODE string = "dev"
	// 默认注册ttl 秒
	DEFAULT_REGISTER_TTL int64 = 5
	// 服务调用地址
	DEFAULT_GRPC_ADVERTISE_ADDR string = "127.0.0.1:28080"
	// 服务注册根路径
	DEFAULT_SCHEMA string = "microkit"
	// DEFAULT_SVC_ID 默认服务ID
	DEFAULT_SVC_ID = "1"
	// 链路追踪服务地址
	DEFAULT_JAEGER_AGENTHOSTPORT = "127.0.0.1:5775"
)

// GetETCDAddr 读取etcd服务地址
func GetETCDAddr() string {
	etcdAddr := os.Getenv("ETCD_ADDR")
	if etcdAddr == "" {
		return DEFAULT_ETCD_ADDR
	}
	return etcdAddr
}

// GetETCDPassword etcd密码
func GetETCDUser() string {
	etcdAddr := os.Getenv("ETCD_USER")
	if etcdAddr == "" {
		return DEFAULT_ETCD_USER
	}
	return etcdAddr
}

// GetETCDPassword etcd密码
func GetETCDPassword() string {
	etcdAddr := os.Getenv("ETCD_PASSWORD")
	if etcdAddr == "" {
		return DEFAULT_ETCD_PASSWORD
	}
	return etcdAddr
}

// GetHTTPAddr 获取配置值
func GetHTTPAddr() string {
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		return DEFAULT_HTTP_ADDR
	}
	return httpAddr
}

// GetGRPCAddr 读取grpc地址
func GetGRPCAddr() string {
	grpcAddr := os.Getenv("GRPC_ADDR")
	if grpcAddr == "" {
		return DEFAULT_GRPC_ADDR
	}
	return grpcAddr
}

// GetTCPAddr 读取tcp地址
func GetTCPAddr() string {
	grpcAddr := os.Getenv("TCP_ADDR")
	if grpcAddr == "" {
		return DEFAULT_TCP_ADDR
	}
	return grpcAddr
}

// GetSvcName 获取服务名 - [redis 使用不传type，一个服务部分类型使用key]
func GetSvcName() string {
	svcName := os.Getenv("SVC_NAME")
	if svcName == "" {
		svcName = DEFAULT_SVC_NAME
	}
	return svcName
}

// GetSvcID 获取服务id
func GetSvcID() string {
	id := os.Getenv("SVC_ID")
	if id == "" {
		id = DEFAULT_SVC_ID
	}
	return id
}

// GetMode 当前运行环境
func GetMode() string {
	mode := os.Getenv("MODE")
	if mode == "" {
		return DEFAULT_MODE
	}
	return mode
}

// GetRegisterTTL 获取注册ttl
func GetRegisterTTL() int64 {
	registerTTL := os.Getenv("DEFAULT_REGISTER_TTL")
	if registerTTL == "" {
		return DEFAULT_REGISTER_TTL
	}
	// 转数字
	registerTTLNum, _ := strconv.Atoi(registerTTL)
	if registerTTLNum == 0 {
		return DEFAULT_REGISTER_TTL
	}
	return int64(registerTTLNum)
}

// GetGRPCAdvertiseAddr 读取grpc注册地址
func GetGRPCAdvertiseAddr() string {
	grpcAdvertiseAddr := os.Getenv("DEFAULT_GRPC_ADVERTISE_ADDR")
	if grpcAdvertiseAddr == "" {
		if DEFAULT_GRPC_ADVERTISE_ADDR == "" {
			return DEFAULT_GRPC_ADDR
		}
		return DEFAULT_GRPC_ADVERTISE_ADDR
	}
	return grpcAdvertiseAddr
}

// GetSchema 获取注册根地址
func GetSchema() string {
	schema := os.Getenv("DEFAULT_SCHEMA")
	if schema == "" {
		return DEFAULT_SCHEMA
	}
	return schema
}

// GetJaegerAgentHostPort 链路追踪服务地址
func GetJaegerAgentHostPort() string {
	schema := os.Getenv("JAEGER_AGENTHOSTPORT")
	if schema == "" {
		return DEFAULT_JAEGER_AGENTHOSTPORT
	}
	return schema
}
