package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type Level int8

type Fields map[string]interface{}

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "panic"
	}
	return ""
}

type Logger struct {
	newLogger *log.Logger
	ctx       context.Context
	fields    Fields
	callers   []string
}

func NewLogger(w io.Writer, perfix string, flag int) *Logger {
	l := log.New(w, perfix, flag)
	return &Logger{newLogger: l}
}

func (l *Logger) clone() *Logger {
	nl := *l
	return &nl
}

func (l *Logger) WithFields(f Fields) *Logger {
	l1 := l.clone()
	if l1.fields == nil {
		l1.fields = make(Fields)
	}

	for k, v := range f {
		l1.fields[k] = v
	}
	return l1
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	l1 := l.clone()
	l1.ctx = ctx
	return l1
}

func (l *Logger) WitchCaller(skip int) *Logger {
	l1 := l.clone()
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		f := runtime.FuncForPC(pc)
		l1.callers = []string{fmt.Sprintf("%s: %d %s", file, line, f.Name())}
	}
	return l1
}

func (l *Logger) WithCallersFrames() *Logger {
	maxCallerDepth := 25
	minCallerDepth := 1
	callers := []string{}
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		s := fmt.Sprintf("%s: %d %s", frame.File, frame.Line, frame.Function)
		callers = append(callers, s)
		if !more {
			break
		}
	}
	l1 := l.clone()
	l1.callers = callers
	return l1
}

func (l *Logger) WithTrace() *Logger {
	ginCtx, ok := l.ctx.(*gin.Context)
	if ok {
		return l.WithFields(Fields{
			"trace_id": ginCtx.MustGet("X-Trace-ID"),
			"span_id":  ginCtx.MustGet("X-Span-ID"),
		})
	}
	return l
}

func (l *Logger) JSONFormat(level Level, message string) map[string]interface{} {
	data := make(Fields, len(l.fields)+4)
	data["level"] = level.String()
	data["time"] = time.Now().Local().UnixNano()
	data["message"] = message
	data["callers"] = l.callers
	if len(l.fields) > 0 {
		for k, v := range l.fields {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}
	return data
}

func (l *Logger) Output(level Level, message string) {
	body, _ := json.Marshal(l.JSONFormat(level, message))
	content := string(body)
	switch level {
	case LevelDebug:
		l.newLogger.Print(content)
	case LevelInfo:
		l.newLogger.Print(content)
	case LevelWarn:
		l.newLogger.Print(content)
	case LevelError:
		l.newLogger.Print(content)
	case LevelFatal:
		l.newLogger.Print(content)
	case LevelPanic:
		l.newLogger.Print(content)
	}
}

func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelInfo, fmt.Sprint(v...))
}

func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelFatal, fmt.Sprint(v...))
}

func (l *Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelFatal, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelError, fmt.Sprint(v...))
}

func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelError, fmt.Sprintf(format, v...))
}

func (l *Logger) Panic(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelPanic, fmt.Sprint(v...))
}

func (l *Logger) Panicf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelPanic, fmt.Sprintf(format, v...))
}
