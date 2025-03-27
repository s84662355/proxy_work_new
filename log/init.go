package log

import (
	"io"
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
	encoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		CallerKey:   "caller",
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		LineEnding: "",
		//	EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeCaller: callerEncodeCaller,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
	infoLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})
	warnLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})
	errorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel
	})

	debugLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel
	})

	dPanicLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DPanicLevel
	})

	panicLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.PanicLevel
	})

	fatalLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.FatalLevel
	})

	isDay := true

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
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(debugWriter), debugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(dPanicWriter), dPanicLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(panicWriter), panicLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(fatalWriter), fatalLevel),
	)

	zaploger = zap.New(core, zap.AddCaller())
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
			rotatelogs.WithRotationTime(time.Hour*24),
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
