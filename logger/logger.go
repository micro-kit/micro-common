package logger

import (
	"github.com/micro-kit/micro-common/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/* zap 日志对象 */

var (
	Logger *zap.SugaredLogger
)

func init() {
	mode := config.GetMode()
	level := zapcore.DebugLevel
	if mode == "pro" {
		level = zapcore.InfoLevel
	}
	// 定时整理日志
	syncWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:  "./logs/error.log",
		MaxSize:   100,
		LocalTime: true,
		Compress:  true,
	})
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.EpochMillisTimeEncoder // 时间格式
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), syncWriter, zap.NewAtomicLevelAt(zapcore.Level(level)))
	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
}
