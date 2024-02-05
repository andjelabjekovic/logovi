package logovi

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
)

func LogInit(path string, appName string) (*log.Logger, mux.MiddlewareFunc, func(msg string), func(msg string)) {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	maxSizeMB := 0.002
	logger.Out = &lumberjack.Logger{
		Filename:   path,
		MaxSize:    int(maxSizeMB), // megabytesloglib
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	/*file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		log.Error("Ne mogu postaviti fajl za čuvanje logova. Koristiću osnovni izlaz.")
	}*/

	writeError := func(msg string) {
		logger.WithFields(log.Fields{
			"id":  uuid.New().String(),
			"app": appName,
		}).Error(msg)
	}
	writeInfo := func(msg string) {
		logger.WithFields(log.Fields{
			"id":  uuid.New().String(),
			"app": appName,
		}).Info(msg)
	}

	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Do stuff here
			logger.WithFields(log.Fields{
				"id":     uuid.New().String(),
				"app":    appName,
				"method": r.Method,
				"url":    r.RequestURI,
				"ip":     r.RemoteAddr,
			}).Info("Request expected")
			//log.Info(r.Method, r.URL, r.Host)

			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		})
	}
	return logger, loggingMiddleware, writeInfo, writeError
}
