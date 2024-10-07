package logger

type Logger interface {
	Debug(message string, keyValues ...any)
	Info(message string, keyValues ...any)
	Warn(message string, keyValues ...any)
	Error(message string, keyValues ...any)
	Fatal(message string, keyValues ...any)
}
