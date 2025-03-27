package log

import (
	"fmt"
	"runtime/debug"
)

func Info(p ...any) {
	if zaploger == nil {
		return
	}
	zaploger.Info(fmt.Sprint(p...))
}

func Warn(p ...interface{}) {
	if zaploger == nil {
		return
	}
	zaploger.Warn(fmt.Sprint(p...))
}

func Error(p ...interface{}) {
	if zaploger == nil {
		return
	}
	zaploger.Error(fmt.Sprint(p...))
}

func Debug(p ...interface{}) {
	if zaploger == nil {
		return
	}
	zaploger.Debug(fmt.Sprint(p...))
}

func DPanic(p ...interface{}) {
	if zaploger == nil {
		return
	}
	zaploger.DPanic(fmt.Sprint(p...))
}

func Panic(p ...interface{}) {
	if zaploger == nil {
		return
	}
	zaploger.Panic(fmt.Sprint(p...))
}

func Fatal(p ...interface{}) {
	if zaploger == nil {
		return
	}
	zaploger.Fatal(fmt.Sprint(p...))
}

// /会抛出异常
func Recover(p ...interface{}) {
	if zaploger == nil {
		return
	}

	err := recover()
	if err != nil {
		zaploger.Panic((fmt.Sprint(p...)) + fmt.Sprint(" error: ", err) + " debug.Stack: " + string(debug.Stack()))
	}
}

// /不会抛出异常
func DRecover(p ...interface{}) {
	if zaploger == nil {
		return
	}
	err := recover()
	if err != nil {
		zaploger.DPanic((fmt.Sprint(p...)) + fmt.Sprint(" error: ", err) + " debug.Stack: " + string(debug.Stack()))
	}
}

func Infof(format string, args ...interface{}) {
	if zaploger == nil {
		return
	}
	zaploger.Info(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...interface{}) {
	if zaploger == nil {
		return
	}
	zaploger.Warn(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	if zaploger == nil {
		return
	}
	zaploger.Error(fmt.Sprintf(format, args...))
}

func Debugf(format string, args ...interface{}) {
	if zaploger == nil {
		return
	}
	zaploger.Debug(fmt.Sprintf(format, args...))
}

func DPanicf(format string, args ...interface{}) {
	if zaploger == nil {
		return
	}
	zaploger.DPanic(fmt.Sprintf(format, args...))
}

func Panicf(format string, args ...interface{}) {
	if zaploger == nil {
		return
	}
	zaploger.Panic(fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...interface{}) {
	if zaploger == nil {
		return
	}
	zaploger.Fatal(fmt.Sprintf(format, args...))
}

// /会抛出异常
func Recoverf(format string, args ...interface{}) {
	if zaploger == nil {
		return
	}
	err := recover()
	if err != nil {
		zaploger.Panic((fmt.Sprintf(format, args...)) + fmt.Sprint(" error: ", err) + " debug.Stack: " + string(debug.Stack()))
	}
}

// /不会抛出异常
func DRecoverf(format string, args ...interface{}) {
	if zaploger == nil {
		return
	}
	err := recover()
	if err != nil {
		zaploger.DPanic((fmt.Sprintf(format, args...)) + fmt.Sprint(" error: ", err) + " debug.Stack: " + string(debug.Stack()))
	}
}
