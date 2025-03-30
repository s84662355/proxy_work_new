package log

import (
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var encoder zapcore.Encoder

var infoLevel, warnLevel,
	errorLevel,
	debugLevel,
	dPanicLevel,
	panicLevel,
	fatalLevel zap.LevelEnablerFunc

var zaploger *zap.Logger

func Init(filepath string) {
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂

	encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		CallerKey:   "file",
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		SkipLineEnding: false,
		LineEnding:     "",
		FunctionKey:    "func",
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	isDay := false

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := getWriter(filepath+"/info", isDay)
	warnWriter := getWriter(filepath+"/warn", isDay)
	errorWriter := getWriter(filepath+"/error", isDay)
	debugWriter := getWriter(filepath+"/debug", isDay)
	dPanicWriter := getWriter(filepath+"/dPanic", isDay)
	panicWriter := getWriter(filepath+"/panic", isDay)
	fatalWriter := getWriter(filepath+"/fatal", isDay)

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), zapcore.InfoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), zapcore.WarnLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), zapcore.ErrorLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(debugWriter), zapcore.DebugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(dPanicWriter), zapcore.DPanicLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(panicWriter), zapcore.PanicLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(fatalWriter), zapcore.FatalLevel),

		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.WarnLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.ErrorLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DPanicLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.PanicLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.FatalLevel),
	)

	zaploger = zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
}

func getWriter(filename string, isDay bool) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志

	if isDay {
		hook, err := rotatelogs.New(
			filename+"_%Y-%m-%d.log", // 没有使用go风格反人类的format格式
			rotatelogs.WithLinkName(filename),
			rotatelogs.WithMaxAge(time.Hour*24*3),
			rotatelogs.WithRotationTime(time.Hour),
		)
		if err != nil {
			panic(err)
		}
		return hook
	}

	hook, err := rotatelogs.New(
		filename+"_%Y-%m-%d-%H.log", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return hook
}
