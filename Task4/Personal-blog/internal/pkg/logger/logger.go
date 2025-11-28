package logger

import (
	"Personal-blog/configs"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// Init 初始化日志
func Init() error {
	// 创建日志目录
	logPath := configs.Config.Logger.Path
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return err
	}

	// 日志文件路径
	logFile := filepath.Join(logPath, configs.Config.Logger.Filename)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// 初始化 logrus
	Log = logrus.New()
	Log.Out = file
	Log.Formatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}

	// 设置日志级别
	level, err := logrus.ParseLevel(configs.Config.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.Level = level

	return nil
}

// 封装常用日志方法（简化调用）
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}
