package logovi

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func LogInit(path string) (*log.Logger, any) {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		log.Error("Ne mogu postaviti fajl za čuvanje logova. Koristiću osnovni izlaz.")
	}

	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Do stuff here
			logger.WithFields(log.Fields{
				"method": r.Method,
				"url":    r.URL,
				"ip":     r.Host,
			}).Info("Failed to send event")
			//log.Info(r.Method, r.URL, r.Host)

			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		})
	}
	return logger, loggingMiddleware
}
