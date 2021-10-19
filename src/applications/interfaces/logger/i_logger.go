package logger_interface

type ILogger interface {
	Info(message string, optional interface{})
	Debug(message string, optional interface{})
	Warn(message string, optional interface{})
	Error(message string, optional interface{})
}
