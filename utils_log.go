package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ufwfqpdgv/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Log struct {
	*zap.Logger
	Config log.Log_info
}

var l *Log

func LogInit(c log.Log_info) {
	hook := lumberjack.Logger{
		Filename:   c.Path_filename, // 日志文件路径
		MaxSize:    c.Max_size,      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: c.Max_backups,   // 日志文件最多保存多少个备份
		MaxAge:     c.Max_age,       // 文件最多保存多少天
		Compress:   c.Compress,      // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     customkTimeEncoder,             // 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 编码器配置
	var encoder zapcore.Encoder
	switch c.Encoding {
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		panic("input on of 'console、json'")
	}
	// 打印到控制台和文件
	var writeSyncer zapcore.WriteSyncer
	if c.Stdout {
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook))
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch c.Level {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zap.ErrorLevel)
	case "panic":
		atomicLevel.SetLevel(zap.PanicLevel)
	case "fatal":
		atomicLevel.SetLevel(zap.FatalLevel)
	default:
		panic("input on of 'debug、info、warn、error、panic、fatal'")
	}
	core := zapcore.NewCore(
		encoder,     // 编码器配置
		writeSyncer, // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	optionArr := make([]zap.Option, 0)
	if c.Development_mode {
		optionArr = append(optionArr, zap.AddCaller())
	}
	// 开启文件及行号
	optionArr = append(optionArr, zap.Development())
	// 设置初始化字段，设置后每行都会带上
	// filed := zap.Fields(zap.String("serviceName", "serviceName"))
	// optionArr = append(optionArr, filed)
	// notice：如用原生的zap下面这里是不用设置的，因自己再加封装了一层，故得caller加1，不然调用的line全是如下面的116
	optionArr = append(optionArr, zap.AddCallerSkip(1))
	// 构造日志
	l = &Log{}
	l.Logger = zap.New(core, optionArr...)
	l.Config = c

	return
}

func customkTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

//下面有点别扭，但方便使用吧
func With(fields ...zap.Field) *Log {
	l.Logger = l.Logger.With(fields...)
	return l
}

func Debug(msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Debug(spew.Sdump(msg...))
		return
	}
	l.Logger.Debug(fmt.Sprint(msg...))
}

func Debugf(format string, msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Debug(spew.Sdump(msg...))
		return
	}
	l.Logger.Debug(fmt.Sprintf(format, msg...))
}

func Info(msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Info(spew.Sdump(msg...))
		return
	}
	l.Logger.Info(fmt.Sprint(msg...))
}

func Infof(format string, msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Info(spew.Sdump(msg...))
		return
	}
	l.Logger.Info(fmt.Sprintf(format, msg...))
}

func Warn(msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Warn(spew.Sdump(msg...))
		return
	}
	l.Logger.Warn(fmt.Sprint(msg...))
}

func Warnf(format string, msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Warn(spew.Sdump(msg...))
		return
	}
	l.Logger.Warn(fmt.Sprintf(format, msg...))
}

func Error(msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Error(spew.Sdump(msg...))
		return
	}
	l.Logger.Error(fmt.Sprint(msg...))
}

func Errorf(format string, msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Error(spew.Sdump(msg...))
		return
	}
	l.Logger.Error(fmt.Sprintf(format, msg...))
}

func Panic(msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Panic(spew.Sdump(msg...))
		return
	}
	l.Logger.Panic(fmt.Sprint(msg...))
}

func Panicf(format string, msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Panic(spew.Sdump(msg...))
		return
	}
	l.Logger.Panic(fmt.Sprintf(format, msg...))
}

func Fatal(msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Fatal(spew.Sdump(msg...))
		return
	}
	l.Logger.Fatal(fmt.Sprint(msg...))
}

func Fatalf(format string, msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Fatal(spew.Sdump(msg...))
		return
	}
	l.Logger.Fatal(fmt.Sprintf(format, msg...))
}

func (*Log) With(fields ...zap.Field) *Log {
	l.Logger = l.Logger.With(fields...)
	return l
}

func (*Log) Debug(msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Debug(spew.Sdump(msg...))
		return
	}
	l.Logger.Debug(fmt.Sprint(msg...))
}

func (*Log) Debugf(format string, msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Debug(spew.Sdump(msg...))
		return
	}
	l.Logger.Debug(fmt.Sprintf(format, msg...))
}

func (*Log) Info(msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Info(spew.Sdump(msg...))
		return
	}
	l.Logger.Info(fmt.Sprint(msg...))
}

func (*Log) Infof(format string, msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Info(spew.Sdump(msg...))
		return
	}
	l.Logger.Info(fmt.Sprintf(format, msg...))
}

func (*Log) Warn(msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Warn(spew.Sdump(msg...))
		return
	}
	l.Logger.Warn(fmt.Sprint(msg...))
}

func (*Log) Warnf(format string, msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Warn(spew.Sdump(msg...))
		return
	}
	l.Logger.Warn(fmt.Sprintf(format, msg...))
}

func (*Log) Error(msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Error(spew.Sdump(msg...))
		return
	}
	l.Logger.Error(fmt.Sprint(msg...))
}

func (*Log) Errorf(format string, msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Error(spew.Sdump(msg...))
		return
	}
	l.Logger.Error(fmt.Sprintf(format, msg...))
}

func (*Log) Panic(msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Panic(spew.Sdump(msg...))
		return
	}
	l.Logger.Panic(fmt.Sprint(msg...))
}

func (*Log) Panicf(format string, msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Panic(spew.Sdump(msg...))
		return
	}
	l.Logger.Panic(fmt.Sprintf(format, msg...))
}

func (*Log) Fatal(msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Fatal(spew.Sdump(msg...))
		return
	}
	l.Logger.Fatal(fmt.Sprint(msg...))
}

func (*Log) Fatalf(format string, msg ...interface{}) {
	if l.Config.Level == "debug" {
		l.Logger.Fatal(spew.Sdump(msg...))
		return
	}
	l.Logger.Fatal(fmt.Sprintf(format, msg...))
}
