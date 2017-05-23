package gerry

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/negroni"
)

var Logger *log.Logger
var ErrorLogger *log.Logger

func DefaultBeforeLogging(entry *log.Entry, req *http.Request, remoteAddr string) *log.Entry {
	return entry.WithFields(log.Fields{
		"REMOTE_ADDR": remoteAddr,
		// "YELLING":        true,
		"REQUEST_METHOD": req.Method,
		"REQUEST":        req.RequestURI,
	})
}

func DefaultAfterLogging(entry *log.Entry, res negroni.ResponseWriter, latency time.Duration, name string) *log.Entry {
	fields := log.Fields{
		// "ALL_DONE":        true,
		"RESPONSE_STATUS": res.Status(),

		fmt.Sprintf("%s_LATENCY", strings.ToUpper(name)): latency,
	}

	// one way to replace an existing entry key
	if requestId, ok := entry.Data["request_id"]; ok {
		fields["REQUEST_ID"] = requestId
		delete(entry.Data, "request_id")
	}

	return entry.WithFields(fields)
}

func doSetupLogger(logger *log.Logger, file string, level log.Level) {
	logger.Level = level
	if file == "" {
		return
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = f
	} else {
		logger.Error("Failed to log to error file, using default stderr")
	}
}

func SetupErrorLogger(file string) {
	ErrorLogger = log.New()
	doSetupLogger(ErrorLogger, file, log.InfoLevel)
}

func SetupLogger(file string) {
	Logger = log.New()
	doSetupLogger(Logger, file, log.InfoLevel)
}
