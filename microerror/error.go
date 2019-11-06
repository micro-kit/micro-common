package microerror

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	goredis "gopkg.in/redis.v5"
)

/* 微服务返回错误代码列表 */

// TODO 此处错误代码应在某区间内为公共错误码、各自程序可以在init中将错误码补充到错误列表

// 基础错误 - 错误码 +-32768
var (
	errors = map[int32]*MicroError{
		10000: NewMicroError(10000, "Unknown server error", nil), // 服务端错误

		10001: NewMicroError(10001, "record not found", gorm.ErrRecordNotFound), // db数据未查询到
		10101: NewMicroError(10101, "key not found", goredis.Nil),               // redis key为空
		10102: NewMicroError(10102, "Redis operation error", nil),               // redis 操作错误

	}
)

// MicroError 错误类型
type MicroError struct {
	Msg  string `json:"msg"`  // 错误信息
	Code int32  `json:"code"` // 错误代码
	Err  error  `json:"err"`  // 原始err信息
}

// NewMicroError 创建MicroError
func NewMicroError(code int32, msg string, err error) *MicroError {
	return &MicroError{
		Code: code,
		Msg:  msg,
		Err:  err,
	}
}

// Error 实现error接口
func (err *MicroError) Error() string {
	js, _ := json.Marshal(err)
	return string(js)
}

// GetMicroError 通过MicroError创建错误，补充err部分
func GetMicroError(code int32, errs ...error) *MicroError {
	var err error
	if len(errs) > 0 {
		err = errs[0]
	}
	microErr := &MicroError{
		Code: code,
		Err:  err,
	}
	if mErr, ok := errors[code]; ok == true {
		microErr.Msg = mErr.Msg
		if err == nil {
			microErr.Err = mErr.Err
		}
	} else {
		microErr.Msg = errors[10000].Msg
	}
	return microErr
}

// InitError 用于注册某个服务的自定义错误
func InitError(errs ...*MicroError) {
	for _, err := range errs {
		errors[err.Code] = err
	}
}
