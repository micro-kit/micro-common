package common

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/micro-kit/micro-common/config"
)

/* 公共函数库 */

// GetRootDir 获取当前可执行程序路径
func GetRootDir() string {
	// 文件不存在获取执行路径
	file, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		file = fmt.Sprintf(".%s", string(os.PathSeparator))
	} else {
		file = fmt.Sprintf("%s%s", file, string(os.PathSeparator))
	}
	return file
}

// PathExists 判断文件或目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetWithTimeout 公共rpc请求超时
func GetWithTimeout(c context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c, time.Duration(config.GetRegisterTTL())*time.Second)
}
