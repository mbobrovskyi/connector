package logger

type Logger interface {
	Debug(args ...any)
	Error(args ...any)
}
