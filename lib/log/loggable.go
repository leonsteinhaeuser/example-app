package log

type Logger interface {
	SetModule(module string) Logger
	Trace() Field
	Debug() Field
	Info() Field
	Warn() Field
	Error(error) Field
	Panic(error) Field
}

type Field interface {
	Field(key string, value interface{}) Field
	Log(message string)
	Logf(format string, args ...interface{})
}
