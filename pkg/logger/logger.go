package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

type Options struct {
	Mode  string
	Level string
	Path  string
	Name  string
}

var logging = new(zap.SugaredLogger)

func NewLogger(o *Options) error {
	var (
		logger = new(zap.Logger)
		err    error
	)

	switch o.Mode {
	case "prod", "PROD":
		level, err := zapcore.ParseLevel(o.Level)
		if err != nil {
			return err
		}

		if o.Path == "" {
			o.Path = "./logs"
		}
		if o.Name == "" {
			path, _ := os.Executable()
			_, exec := filepath.Split(path)
			o.Name = exec
		}
		fileName := filepath.Join(o.Path, o.Name+".log")
		syncWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    100,
			MaxBackups: 10,
			LocalTime:  true,
			Compress:   true,
		})
		encoder := zap.NewProductionEncoderConfig()
		encoder.EncodeTime = zapcore.ISO8601TimeEncoder
		core := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), syncWriter, zap.NewAtomicLevelAt(level))
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	default:
		logger, err = zap.NewDevelopment(zap.AddCallerSkip(1))
		if err != nil {
			return err
		}
	}
	logging = logger.Sugar()
	return nil
}

func Debug(args ...interface{}) {
	logging.Debug(args...)
}
func Debugw(msg string, keysAndValues ...interface{}) {
	logging.Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	logging.Info(args...)
}
func Infow(msg string, keysAndValues ...interface{}) {
	logging.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	logging.Warn(args...)
}
func Warnw(msg string, keysAndValues ...interface{}) {
	logging.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	logging.Error(args...)
}
func Errorw(msg string, keysAndValues ...interface{}) {
	logging.Errorw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	logging.Fatal(args...)
}
