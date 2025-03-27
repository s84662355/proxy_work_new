package log

import (
	"fmt"
	"runtime"

	"go.uber.org/zap/zapcore"
)

var callerEncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	pc, file, line, ok := runtime.Caller(6)
	if !ok {
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	enc.AppendString(fmt.Sprintf("file: %s  line: %d    func: %s  ", file, line, funcName))
}
