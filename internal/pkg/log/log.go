package log

import (
	"QA-System/internal/global/config"
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var zapStacktraceMutex sync.Mutex
var LogDir string
var LogName string

type Config struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
	Name              string
	Writers           string
	LoggerDir         string
	LogRollingPolicy  string
	LogBackupCount    uint
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

const (
	// WriterConsole console输出
	WriterConsole = "console"
	// WriterFile 文件输出
	WriterFile = "file"

	LogSuffix      = ".log"
	WarnLogSuffix  = "_warn.log"
	ErrorLogSuffix = "_error.log"
)

const (
	RotateTimeDaily  = "daily"
	RotateTimeHourly = "hourly"
)

func loadConfig() *Config {
	return &Config{
		Development:       global.Config.GetBool("log.development"),
		DisableCaller:     global.Config.GetBool("log.disableCaller"),
		DisableStacktrace: global.Config.GetBool("log.disableStacktrace"),
		Encoding:          global.Config.GetString("log.encoding"),
		Level:             global.Config.GetString("log.level"),
		Name:              global.Config.GetString("log.name"),
		Writers:           global.Config.GetString("log.writers"),
		LoggerDir:         global.Config.GetString("log.loggerDir"),
		LogRollingPolicy:  global.Config.GetString("log.logRollingPolicy"),
		LogBackupCount:    global.Config.GetUint("log.logBackupCount"),
	}
}

func ZapInit() {
	cfg := loadConfig()

	LogDir = cfg.LoggerDir
	LogName = cfg.Name
	if strings.HasSuffix(LogDir, "/") {
		LogDir = strings.TrimRight(LogDir, "/")
	}

	if err := os.MkdirAll(cfg.LoggerDir, 0755); err != nil {
		zap.S().Error("创建日志目录失败:", err)
		return
	}

	var encoderCfg zapcore.EncoderConfig
	if cfg.Development {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if cfg.Encoding == WriterConsole {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	var cores []zapcore.Core
	var options []zap.Option

	options = append(options, zap.Fields(zap.String("serviceName", cfg.Name)))

	writers := strings.Split(cfg.Writers, ",")
	for _, w := range writers {
		switch w {
		case WriterConsole:
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLoggerLevel(cfg)))
		case WriterFile:
			cores = append(cores, getInfoCore(encoder, cfg))
			core, option := getWarnCore(encoder, cfg)
			cores = append(cores, core)
			if option != nil {
				options = append(options, option)
			}

			core, option = getErrorCore(encoder, cfg)
			cores = append(cores, core)
			if option != nil {
				options = append(options, option)
			}
		default:
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLoggerLevel(cfg)))
			cores = append(cores, getAllCore(encoder, cfg))
		}
	}

	combinedCore := zapcore.NewTee(cores...)

	if !cfg.DisableCaller {
		options = append(options, zap.AddCaller())
	}

	if !cfg.DisableStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

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

func getAllCore(encoder zapcore.Encoder, cfg *Config) zapcore.Core {
	allWriter := getLogWriterWithTime(cfg, GetLogFile(cfg.Name, LogSuffix))
	allLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.FatalLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(allWriter), allLevel)
}

func getInfoCore(encoder zapcore.Encoder, cfg *Config) zapcore.Core {
	infoWrite := getLogWriterWithTime(cfg, GetLogFile(cfg.Name, LogSuffix))
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.InfoLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(infoWrite), infoLevel)
}

func getWarnCore(encoder zapcore.Encoder, cfg *Config) (zapcore.Core, zap.Option) {
	warnWrite := getLogWriterWithTime(cfg, GetLogFile(cfg.Name, WarnLogSuffix))
	var stacktrace zap.Option
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if !cfg.DisableCaller {
			zapStacktraceMutex.Lock()
			stacktrace = zap.AddStacktrace(zapcore.WarnLevel)
			zapStacktraceMutex.Unlock()
		}
		return lvl == zapcore.WarnLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(warnWrite), warnLevel), stacktrace
}

func getErrorCore(encoder zapcore.Encoder, cfg *Config) (zapcore.Core, zap.Option) {
	errorWrite := getLogWriterWithTime(cfg, GetLogFile(cfg.Name, ErrorLogSuffix))
	var stacktrace zap.Option
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if !cfg.DisableCaller {
			zapStacktraceMutex.Lock()
			stacktrace = zap.AddStacktrace(zapcore.ErrorLevel)
			zapStacktraceMutex.Unlock()
		}
		return lvl >= zapcore.ErrorLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(errorWrite), errorLevel), stacktrace
}

func getLogWriterWithTime(cfg *Config, filename string) io.Writer {
	logFullPath := filename
	rotationPolicy := cfg.LogRollingPolicy
	backupCount := cfg.LogBackupCount

	var (
		rotateDuration time.Duration
		timeFormat     string
	)
	if rotationPolicy == RotateTimeHourly {
		rotateDuration = time.Hour
		timeFormat = ".%Y%m%d%H"
	} else if rotationPolicy == RotateTimeDaily {
		rotateDuration = time.Hour * 24
		timeFormat = ".%Y%m%d"
	}

	// 检查日志文件是否存在
	if _, err := os.Stat(logFullPath); os.IsNotExist(err) {
		// 如果日志文件不存在，创建它
		file, err := os.Create(logFullPath)
		if err != nil {
			zap.S().Error("Failed to create log file:", err)
			panic(err)
		}
		file.Close()

		// 设置日志文件权限为 0644
		err = os.Chmod(logFullPath, 0644)
		if err != nil {
			zap.S().Error("Failed to set log file permissions:", err)
			panic(err)
		}
	}

	hook, err := rotatelogs.New(
		logFullPath+time.Now().Format(timeFormat),
		rotatelogs.WithLinkName(logFullPath),
		rotatelogs.WithRotationCount(backupCount),
		rotatelogs.WithRotationTime(rotateDuration),
	)

	if err != nil {
		zap.S().Error("Failed to initialize log rotation:", err)
		panic(err)
	}
	return hook
}

func GetLogFile(filename string, suffix string) string {
	return ConcatString(global.Config.GetString("log.loggerDir"), "/", filename, suffix)
}

func ConcatString(s ...string) string {
	if len(s) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for _, i := range s {
		buffer.WriteString(i)
	}
	return buffer.String()
}
