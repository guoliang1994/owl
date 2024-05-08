package owl

import (
	"go.uber.org/zap/zapcore"
	"owl/contract"
	"owl/log"
)

type LoggerFactory struct {
	stage *Stage
	opt   *Options
}

type Options struct {
	Level      int  `json:"level"`
	MaxSize    int  `json:"max-size"`
	MaxBackups int  `json:"max-backups"`
	MaxAge     int  `json:"max-age"`
	Compress   bool `json:"compress"`
}

func NewLoggerFactory(stage *Stage, cfgManager *ConfManager) *LoggerFactory {
	var opt *Options
	_ = cfgManager.GetConfig("logs", &opt)
	return &LoggerFactory{
		opt:   opt,
		stage: stage,
	}
}

// RuntimeLogger 返回运行日志
func (i *LoggerFactory) RuntimeLogger() contract.Logger {
	if i.opt == nil {
		return log.ConsoleImpl{}
	}
	level := zapcore.Level(i.opt.Level)
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
	if i.opt == nil {
		return log.ConsoleImpl{}
	}
	options := log.Options{
		StorePath:  LogsPath,
		Channel:    channel,
		MaxSize:    i.opt.MaxSize,
		MaxBackups: i.opt.MaxBackups,
		MaxAge:     i.opt.MaxAge,
		Compress:   i.opt.Compress,
		Level:      level,
	}

	return log.NewFileImpl(&options)
}
