package log

import (
	"github.com/golang-module/carbon"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"sync"
)

/*
日志初始化步骤：
1. 设置 log 写入的文件 (writer)
2. 设置为日志编码的方法（encoder）
3. 创建日志核心
*/

type FileImpl struct {
	options    *Options
	l          *zap.SugaredLogger
	lock       sync.RWMutex
	preGetTime string
}

type Channel string

func (i Channel) String() string {
	return string(i)
}

const (
	RUNTIME Channel = "runtime" // 运行时日志
	SQL     Channel = "sql"     // sql 日志
	ACCESS  Channel = "access"  // 访问日志
)

type Options struct {
	StorePath  string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Level      zapcore.Level

	dateFileName string // 计算出来，每天一个文件
	Channel      Channel
}

var (
	fileLoggerMap  = make(map[Channel]*FileImpl)
	defaultOptions = &Options{
		StorePath:    "./storage",
		dateFileName: "0000-00-00",
		Channel:      RUNTIME,
		MaxSize:      50,
		MaxBackups:   100,
		MaxAge:       30,
		Compress:     true,
		Level:        zap.DebugLevel,
	}
)

func NewFileImpl(options *Options) *FileImpl {
	impl, ok := fileLoggerMap[options.Channel]
	if ok {
		return impl
	}

	if options == nil {
		options = defaultOptions
	}

	l := &FileImpl{
		options: options,
	}
	fileLoggerMap[options.Channel] = l
	return l
}

func (f *FileImpl) Emergency(content ...any) {
	f.setLogger()
	f.l.DPanicln(content)
}

func (f *FileImpl) Alert(content ...any) {
	f.setLogger()
	f.l.Errorln(content)
}

func (f *FileImpl) Critical(content ...any) {
	f.setLogger()
}

func (f *FileImpl) Error(content ...any) {
	f.setLogger()
	f.l.Errorln(content)
}

func (f *FileImpl) Warning(content ...any) {
	f.setLogger()
	f.l.Warnln(content)
}

func (f *FileImpl) Notice(content ...any) {
	f.setLogger()
	f.l.Warnln(content)
}

func (f *FileImpl) Info(content ...any) {
	f.setLogger()
	f.l.Infoln(content)
}

func (f *FileImpl) Debug(content ...any) {
	f.setLogger()
	f.l.Debugln(content)
}

// 根据时间轮换，因为随时都有可能发生时间变化，调用日志方法之前需要先调用这个方法
func (f *FileImpl) setLogger() {

	f.lock.Lock()
	defer f.lock.Unlock()

	now := carbon.Now()
	nowDate := now.ToDateString()

	// 重新换一个文件来记录日志
	if f.preGetTime != "" {
		if nowDate != f.preGetTime {
			f.l = nil
			f.preGetTime = nowDate
		}
	} else {
		f.preGetTime = nowDate
	}

	if f.l != nil {
		return
	}

	f.options.dateFileName = f.options.Channel.String() + "-" + nowDate

	encoder := getEncoder()
	writeSyncer := getLogWriter(f.options)
	// 打印
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(writeSyncer)),
		f.options.Level,
	)

	// zap.AddCaller()  添加将调用函数信息记录到日志中的功能。
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	f.l = logger.Sugar()
}

func getLogWriter(options *Options) *lumberjack.Logger {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   options.StorePath + "/" + options.dateFileName + ".log", // 日志文件路径
		MaxSize:    options.MaxSize,                                         // 最大尺寸, M
		MaxBackups: options.MaxBackups,                                      // 备份数 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxAge:     options.MaxAge,                                          // 保留旧文件的最大天数
		Compress:   options.Compress,                                        // 是否压缩/归档旧文件
	}

	return lumberJackLogger
}

func getEncoder() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000") // 修改时间编码器

	// 在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return encoderConfig
}
