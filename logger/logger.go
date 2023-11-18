package logger

type Logger interface {
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Fatalln(msg string, args ...interface{})
	Named(name string) Logger
}
