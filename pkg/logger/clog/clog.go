package clog

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Level int

type Leveler interface {
	Level() Level
	String() string
}

const (
	levelDebug Level = iota
	levelInfo
	levelWarn
	levelError
	levelFatal
)

var (
	LevelDebug Leveler = newLeveler(levelDebug, "DEBUG")
	LevelInfo          = newLeveler(levelInfo, "INFO")
	LevelWarn          = newLeveler(levelWarn, "WARN")
	LevelError         = newLeveler(levelError, "ERROR")
	LevelFatal         = newLeveler(levelFatal, "FATAL")
)

type stdLeveler struct {
	l Level
	s string
}

func newLeveler(l Level, s string) Leveler {
	return &stdLeveler{l: l, s: s}
}

func (l *stdLeveler) Level() Level   { return l.l }
func (l *stdLeveler) String() string { return l.s }

const (
	FlagTime = 1 << iota
	FlagDate
	FlagSource
	FlagLevel
	FlagStd = FlagTime | FlagDate | FlagSource | FlagLevel
)

type Logger struct {
	w      io.Writer
	mu     *sync.Mutex
	ctx    context.Context
	level  Leveler
	flags  int
	prefix string
}

func New(w io.Writer, level Leveler, flags int) *Logger {
	return NewWithContext(context.Background(), w, level, flags)
}

func NewWithContext(ctx context.Context, w io.Writer, level Leveler, flags int) *Logger {
	return &Logger{
		w:     w,
		mu:    &sync.Mutex{},
		ctx:   ctx,
		level: level,
		flags: flags,
	}
}

func (l *Logger) Debug(message string, keyValues ...any) {
	l.print(l.ctx, LevelDebug, message, keyValues...)
}

func (l *Logger) Info(message string, keyValues ...any) {
	l.print(l.ctx, LevelInfo, message, keyValues...)
}

func (l *Logger) Warn(message string, keyValues ...any) {
	l.print(l.ctx, LevelWarn, message, keyValues...)
}

func (l *Logger) Error(message string, keyValues ...any) {
	l.print(l.ctx, LevelError, message, keyValues...)
}

func (l *Logger) Fatal(message string, keyValues ...any) {
	l.print(l.ctx, LevelFatal, message, keyValues...)
	os.Exit(1)
}

func (l *Logger) enabled(_ context.Context, level Leveler) bool {
	ok := l.level.Level() <= level.Level()
	return ok
}

func (l *Logger) hasFlag(flag int) bool {
	return l.flags&flag != 0
}

// TODO отформатировать строку источника, дабы она была кликабельной
// TODO добавить поддержку как обычного текста, так и формата JSON
// TODO улучшить структуру пакета
func (l *Logger) print(ctx context.Context, level Leveler, message string, args ...any) {
	if !l.enabled(ctx, level) {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("02.01.2006 15:04:05")
	timeAndDate := strings.Split(timestamp, " ")
	if l.hasFlag(FlagDate) {
		if _, err := l.w.Write([]byte(timeAndDate[0] + " ")); err != nil {
			return
		}
	}
	if l.hasFlag(FlagTime) {
		if _, err := l.w.Write([]byte(timeAndDate[1] + " ")); err != nil {
			return
		}
	}
	if l.hasFlag(FlagSource) {
		pc, _, _, ok := runtime.Caller(2)
		details := runtime.FuncForPC(pc)
		if !ok || details == nil {
			return
		}
		if _, err := l.w.Write([]byte(details.Name() + " ")); err != nil {
			return
		}
	}
	if l.hasFlag(FlagLevel) {
		if _, err := l.w.Write([]byte(level.String() + " ")); err != nil {
			return
		}
	}

	var attrsString string
	if len(args) > 1 {
		attrs := make([]string, 0)
		for i, v := range args {
			if (i+1)%2 != 0 {
				continue
			}
			attrs = append(attrs, fmt.Sprintf("%s=\"%+v\"", args[i-1], v))
		}
		attrsString = " " + strings.Join(attrs, " ")
	}

	if _, err := l.w.Write([]byte(message + attrsString + "\n")); err != nil {
		return
	}
}
