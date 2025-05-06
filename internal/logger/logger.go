package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log is the global logger instance
	Log  *zap.SugaredLogger
	once sync.Once
)

// Init initializes the logger with the specified level
func Init(level string) error {
	var initErr error
	once.Do(func() {
		// Parse log level
		var zapLevel zapcore.Level
		if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
			zapLevel = zapcore.InfoLevel
		}

		// Create encoder config for console output
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// Create a custom encoder that adds a prefix to our logs
		encoder := zapcore.NewConsoleEncoder(encoderConfig)
		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stderr),
			zapLevel,
		)

		// Create logger with additional options
		logger := zap.New(core,
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
			zap.Development(),
			zap.Fields(zap.String("plugin", "ti-scaffold")), // Add plugin name to all logs
		)
		Log = logger.Sugar()
	})

	return initErr
}

// Sync flushes any buffered log entries
func Sync() error {
	if Log != nil {
		return Log.Sync()
	}
	return nil
}
