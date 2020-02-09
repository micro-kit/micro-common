package microerror

import (
	"encoding/json"

	"google.golang.org/grpc/status"
)

/* 微服务返回错误代码列表 */

// TODO 此处错误代码应在某区间内为公共错误码、各自程序可以在init中将错误码补充到错误列表

const (
	// UnknownServerError 服务端默认错误
	UnknownServerError uint32 = 10000
)

// 基础错误 - 错误码为int32数字
var (
	errors = map[uint32]*MicroError{
		/* 常用基础错误 */
		10000: NewMicroError(10000, "Unknown server error"), // 服务端错误
		10001: NewMicroError(10001, "record not found"),     // db数据未查询到
		10002: NewMicroError(10002, "Parameter error"),      // 参数错误

	}
)

// MicroError 错误类型
type MicroError struct {
	Msg  string `json:"msg"`  // 错误信息
	Code uint32 `json:"code"` // 错误代码
}

// NewMicroError 创建MicroError
func NewMicroError(code uint32, msg string) *MicroError {
	return &MicroError{
		Code: code,
		Msg:  msg,
	}
}

// Error 实现error接口
func (err *MicroError) Error() string {
	js, _ := json.Marshal(err)
	return string(js)
}

// GetMicroError 通过MicroError创建错误，补充err部分
func GetMicroError(code uint32, errs ...error) *MicroError {
	var err error
	msg := ""
	if len(errs) > 0 {
		err = errs[0]
		if err != nil {
			msg = err.Error()
		}
	}
	microErr := &MicroError{
		Code: code,
		Msg:  msg,
	}
	if mErr, ok := errors[code]; ok == true {
		microErr.Msg = mErr.Msg
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

// Convert 解析错误信息 - 统一客户端创建时返回此错误
func Convert(err error) *MicroError {
	if err == nil {
		return nil
	}
	grpcStatus := status.Convert(err)
	return &MicroError{
		Code: uint32(grpcStatus.Code()),
		Msg:  grpcStatus.Message(),
	}
}
