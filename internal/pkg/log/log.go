package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func ZapInit() {
	logFilePath := "logs/app.log"
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		// 如果日志文件不存在，则创建新的日志文件
		file, err := os.Create(logFilePath)
		if err != nil {
			// 创建日志文件失败，记录错误并返回
			zap.S().Error("创建日志文件失败:", err)
			return
		}
		file.Close()
	}

	// 以追加模式打开日志文件
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// 打开日志文件失败，记录错误并返回
		zap.S().Error("打开日志文件失败:", err)
		return
	}

	writeSyncer := zapcore.AddSync(logFile)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
