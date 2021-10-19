package middlewares

import (
	logger_interface "base/src/applications/interfaces/logger"
	"net/http"
)

type logInfo struct {
	Method string
	Host   string
	From   string
}

type LoggingMiddleware struct {
	logger logger_interface.ILogger
}

func (lm LoggingMiddleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logInfo := &logInfo{
			Method: r.Method,
			Host:   r.Host,
			From:   r.RemoteAddr,
		}

		lm.logger.Info(r.RequestURI, logInfo)

		next.ServeHTTP(w, r)
	})
}

func NewLoggingMiddleware(logger logger_interface.ILogger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}
