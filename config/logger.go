package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

func InitLogger() *zap.SugaredLogger {
	//设置日志等级(读取配置文件的log.level来设置对应的日志等级)
	logModel := zapcore.InfoLevel
	switch viper.GetString("logs.level") {
	case "debug":
		logModel = zapcore.DebugLevel
	case "info":
		logModel = zapcore.InfoLevel
	case "warn":
		logModel = zapcore.WarnLevel
	case "error":
		logModel = zapcore.ErrorLevel
	}
	encoder := encoder()
	loginfo := ws()
	core := zapcore.NewCore(encoder, loginfo, logModel)
	return zap.New(core).Sugar()
}

func encoder() zapcore.Encoder {
	//获取一个适用于生产环境的默认编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	//设置时间字段在日志消息中的键名为 “time
	encoderConfig.TimeKey = "time"

	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	return zapcore.NewJSONEncoder(encoderConfig) //使用配置创建一个新的JSON编码器
}

func ws() zapcore.WriteSyncer {
	//logPath := viper.GetString("logs.path") + time.Now().Format(time.DateOnly) + ".log"
	stSeparator := string(filepath.Separator) //分隔符
	stRootDir, _ := os.Getwd()
	stLogFilePath := stRootDir + stSeparator + "log" + stSeparator + time.Now().Format(time.DateOnly)
	//日志分割
	l := &lumberjack.Logger{
		Filename:   stLogFilePath,
		MaxSize:    viper.GetInt("logs.max_size"), // megabytes
		MaxBackups: viper.GetInt("logs.max_backups"),
		MaxAge:     viper.GetInt("logs.max_age"),   // 保留旧文件最大天数
		Compress:   viper.GetBool("logs.compress"), // disabled by default
	}
	return zapcore.AddSync(l)
}
