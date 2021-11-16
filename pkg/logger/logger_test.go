package logger

import (
	"errors"
	"testing"
)

func TestJSONLogger(t *testing.T) {

	Config := &ZapConfig{
		Level: 0,
		Fields: map[string]string{
			"servicename": "testservice",
		},
		DisableConsole: true,
		FileName:       "./logs/test.log",
	}

	logger := NewJSONLogger(Config)
	defer logger.Sync()

	err := errors.New("pkg error")
	logger.Error("err occurs", WrapMeta(nil, NewMeta("para1", "value1"), NewMeta("para2", "value2"))...)
	logger.Error("err occurs", WrapMeta(err, NewMeta("para1", "value1"), NewMeta("para2", "value2"))...)

}

func BenchmarkJsonLogger(b *testing.B) {
	b.ResetTimer()
	Config := &ZapConfig{
		Level: 0,
		Fields: map[string]string{
			"defined_key": "defined_value",
		},
		DisableConsole: false,
	}
	logger := NewJSONLogger(Config)

	defer logger.Sync()

}

// logger.Info(path,
// 	zap.Int("status", c.Writer.Status()),
// 	zap.String("method", c.Request.Method),
// 	zap.String("path", path),
// 	zap.String("query", query),
// 	zap.String("ip", c.ClientIP()),
// 	zap.String("user-agent", c.Request.UserAgent()),
// 	zap.String("time", end.Format(timeFormat)),
// 	zap.Duration("latency", latency),
// )
