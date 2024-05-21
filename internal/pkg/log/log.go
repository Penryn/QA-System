package log

import (
	"QA-System/internal/global/config"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// Config 结构体定义了日志配置参数
type Config struct {
	Development       bool   // 是否是开发模式
	DisableCaller     bool   // 是否禁用调用者信息
	DisableStacktrace bool   // 是否禁用堆栈跟踪信息
	Encoding          string // 日志编码格式，可以是 "console" 或 "json"
	Level             string // 日志记录的最低级别
	Name              string // 服务或应用程序的名称
	Writers           string // 日志输出目标，可以是 "console", "file" 或者它们的组合
	LoggerDir         string // 日志文件的存储目录
}

const (
	WriterConsole = "console"
	WriterFile    = "file"
)

// ZapInit 初始化 Zap 日志记录器

func  loadConfig() *Config {
	return &Config{
		Development:       global.Config.GetBool("log.development"),
		DisableCaller:     global.Config.GetBool("log.disableCaller"),
		DisableStacktrace: global.Config.GetBool("log.disableStacktrace"),
		Encoding:          global.Config.GetString("log.encoding"),
		Level:             global.Config.GetString("log.level"),
		Name:              global.Config.GetString("log.name"),
		Writers:           global.Config.GetString("log.writers"),
		LoggerDir:         global.Config.GetString("log.loggerDir"),
	}
}

func ZapInit() {
	// Ensure the log directory exists
	cfg := loadConfig()
	if err := os.MkdirAll(cfg.LoggerDir, 0755); err != nil {
		zap.S().Error("创建日志目录失败:", err)
		return
	}

	logFilePath := cfg.LoggerDir + "/app.log"

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

	// 设置编码器配置
	var encoderCfg zapcore.EncoderConfig
	if cfg.Development {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// 根据配置选择编码器类型
	var encoder zapcore.Encoder
	if cfg.Encoding == WriterConsole {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	var cores []zapcore.Core
	var options []zap.Option

	// 添加默认字段
	options = append(options, zap.Fields(zap.String("serviceName", cfg.Name)))

	writers := strings.Split(cfg.Writers, ",")
	for _, w := range writers {
		switch w {
		case WriterConsole:
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLoggerLevel(cfg)))
		case WriterFile:
			cores = append(cores, zapcore.NewCore(encoder, writeSyncer, getLoggerLevel(cfg)))

			// info
			cores = append(cores, getInfoCore(encoder, cfg))

			// warning
			core, option := getWarnCore(encoder, cfg)
			cores = append(cores, core)
			if option != nil {
				options = append(options, option)
			}

			// error
			core, option = getErrorCore(encoder, cfg)
			cores = append(cores, core)
			if option != nil {
				options = append(options, option)
			}
		default:
			// console
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLoggerLevel(cfg)))
			// file
			cores = append(cores, zapcore.NewCore(encoder, writeSyncer, getLoggerLevel(cfg)))
		}
	}

	combinedCore := zapcore.NewTee(cores...)

	// 开启开发模式，堆栈跟踪
	if !cfg.DisableCaller {
		options = append(options, zap.AddCaller())
	}

	if !cfg.DisableStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// 构造日志
	Logger = zap.New(combinedCore, options...)
	Logger.Info("Logger initialized")
}

func getLoggerLevel(cfg *Config) zapcore.Level {
	level, exist := loggerLevelMap[strings.ToLower(cfg.Level)]
	if !exist {
		return zapcore.DebugLevel
	}
	return level
}

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getInfoCore(encoder zapcore.Encoder, cfg *Config) zapcore.Core {
	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
}

func getWarnCore(encoder zapcore.Encoder, cfg *Config) (zapcore.Core, zap.Option) {
	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.WarnLevel), nil
}

func getErrorCore(encoder zapcore.Encoder, cfg *Config) (zapcore.Core, zap.Option) {
	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.ErrorLevel), nil
}
