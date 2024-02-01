package owl

import (
	"go.uber.org/zap/zapcore"
	"owl/contract"
	"owl/log"
)

type LoggerFactory struct {
	stage  *Stage
	logCfg GetConfigFunc
}

func NewLoggerFactory(stage *Stage) *LoggerFactory {
	return &LoggerFactory{stage: stage}
}

// RuntimeLogger 返回运行日志
func (i *LoggerFactory) RuntimeLogger() contract.Logger {
	i.logCfg = i.stage.ConfManager.GetConfig("logs").Get
	level := zapcore.Level(i.logCfg("level").ToInt())
	return i.getLogger(log.RUNTIME, level)
}

// SqlLogger 返回 SQL 日志
func (i *LoggerFactory) SqlLogger() contract.Logger {
	return i.getLogger(log.SQL, zapcore.InfoLevel)
}

// AccessLogger 返回访问日志
func (i *LoggerFactory) AccessLogger() contract.Logger {
	return i.getLogger(log.ACCESS, zapcore.InfoLevel)
}

func (i *LoggerFactory) getLogger(channel log.Channel, level zapcore.Level) contract.Logger {
	logCfg := i.logCfg
	options := log.Options{
		StorePath:  LogsPath,
		Channel:    channel,
		MaxSize:    logCfg("max-size").ToInt(),
		MaxBackups: logCfg("max-backups").ToInt(),
		MaxAge:     logCfg("max-age").ToInt(),
		Compress:   logCfg("compress").ToBool(),
		Level:      level,
	}

	return log.NewFileImpl(&options)
}
