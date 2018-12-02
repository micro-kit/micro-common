package cache

import (
	"strings"

	"github.com/micro-kit/micro-common/config"
	uuid "github.com/satori/go.uuid"
)

/* 格式化redis key */

// GetPrefixKey 给key添加前缀
func GetPrefixKey(key string) string {
	svcName := config.GetSvcName("")
	svcName = strings.Trim(svcName, "/")
	return "vivi:" + svcName + ":" + key
}

// GetUserLoginToken 获取token存储key
func GetUserLoginToken() (key, token string) {
	tokenUUID, _ := uuid.NewV4()
	token = tokenUUID.String()
	key = GetPrefixKey("login:token:" + token)
	return
}
