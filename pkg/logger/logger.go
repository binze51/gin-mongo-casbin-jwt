package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DefaultLevel      = zapcore.InfoLevel
	DefaultTimeLayout = time.RFC3339
)

type ZapConfig struct {
	Level          uint              `yaml:"level"` //配置日志等级，0默认全部
	Fields         map[string]string // 自定义消息json的字段
	DisableConsole bool              `yaml:"disableconsole"`
	FileName       string            `yaml:"logfile"`
}

// NewJSONLogger json日志器
func NewJSONLogger(conf *ZapConfig) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger", // used by logger.Named(key); optional; useless
		CallerKey:      "caller", //
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace", // use by zap.AddStacktrace; optional; useless
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // time.RFC3339
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 调用者路径
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// lowPriority usd by info\debug\warn
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.Level(conf.Level) && lvl < zapcore.ErrorLevel
	})

	// highPriority usd by error\panic\fatal
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.Level(conf.Level) && lvl >= zapcore.ErrorLevel
	})

	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe

	core := zapcore.NewTee()

	if !conf.DisableConsole {
		core = zapcore.NewTee(
			zapcore.NewCore(jsonEncoder,
				zapcore.NewMultiWriteSyncer(stderr),
				highPriority,
			),
			zapcore.NewCore(jsonEncoder,
				zapcore.NewMultiWriteSyncer(stdout),
				lowPriority,
			),
		)
	}

	if conf.FileName != "" {
		dir := filepath.Dir(conf.FileName)
		if err := os.MkdirAll(dir, 0766); err != nil {
			panic(err)
		}
		logfile := &lumberjack.Logger{ // concurrent-safed
			Filename:   conf.FileName, // 文件路径
			MaxSize:    128,           // 单个文件最大尺寸，默认单位 M
			MaxBackups: 300,           // 最多保留 300 个备份
			MaxAge:     30,            // 最大时间，默认单位 day
			LocalTime:  true,          // 使用本地时间
			Compress:   true,          // 是否压缩 disabled by default
		}

		core = zapcore.NewTee(core,
			zapcore.NewCore(jsonEncoder,
				zapcore.AddSync(logfile),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl >= zapcore.Level(conf.Level)
				}),
			),
		)
	}
	logger := zap.New(core,
		zap.AddCaller(),
		zap.ErrorOutput(stderr),
	)

	for key, value := range conf.Fields {
		logger = logger.WithOptions(zap.Fields(zapcore.Field{Key: key, Type: zapcore.StringType, String: value}))
	}
	return logger
}
