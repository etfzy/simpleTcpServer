package simpletcp

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func createDefaultLog() *zap.Logger {
	//获取编码器
	encoderConfig := zap.NewProductionEncoderConfig() //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	//日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.InfoLevel
	})

	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./simpletcp.json", //日志文件存放目录
		MaxSize:    100,                //文件大小限制,单位MB
		MaxBackups: 20,                 //最大保留日志文件数量
		MaxAge:     30,                 //日志文件保留天数
		Compress:   false,              //是否压缩处理
	})

	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer), highPriority)
	//errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer,zapcore.AddSync(os.Stdout)), highPriority)
	logzap := zap.New(errorFileCore)
	return logzap
}
