package logger

type NopLogger struct{}

func (n *NopLogger) Debug(args ...any) {}
func (n *NopLogger) Error(args ...any) {}

func NewNopLogger() *NopLogger {
	return &NopLogger{}
}
