package log

import (
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"path"
	"runtime"
	"strings"
)

// Fields allows passing key value pairs to Logrus
type Fields map[string]interface{}

type contextHook struct {
	ReqID string
}

var contextHookObj = contextHook{}

func (hook *contextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *contextHook) Fire(entry *logrus.Entry) error {
	pc := make([]uintptr, 5, 5)
	cnt := runtime.Callers(6, pc)

	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 2)
		name := fu.Name()
		if !strings.Contains(name, "github.com/sirupsen/logrus") &&
			!strings.Contains(name, "framework/log") {
			file, line := fu.FileLine(pc[i] - 2)
			entry.Data["file"] = path.Base(file)
			entry.Data["func"] = path.Base(name)
			entry.Data["line"] = line
			if hook.ReqID != "" {
				entry.Data["reqid"] = hook.ReqID
			}
			break
		}
	}
	return nil
}

func init() {
	logrus.SetFormatter(&prefixed.TextFormatter{TimestampFormat: "Jan 02 03:04:05.000"})
	logrus.AddHook(&contextHookObj)
}

// WithField adds a field to the logrus entry
func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

// WithFields add fields to the logrus entry
func WithFields(fields Fields) *logrus.Entry {
	sendfields := make(logrus.Fields)
	for k, v := range fields {
		sendfields[k] = v
	}
	return logrus.WithFields(sendfields)
}

// WithAddFields add fields to the logrus entry
func WithAddFields(fields Fields, fields2 Fields) *logrus.Entry {
	return withAddFields(fields, fields2)
}

func withAddFields(fields Fields, fields2 Fields) *logrus.Entry {
	sendfields := make(logrus.Fields)
	for k, v := range fields {
		sendfields[k] = v
	}
	for k, v := range fields2 {
		sendfields[k] = v
	}
	return logrus.WithFields(sendfields)
}

// FuncStart logs the start of a function
func FuncStart(logFields Fields) {
	withAddFields(logFields, Fields{"status": "start"}).Info()
}

// FuncEnd logs the "end" of a function, typically at the end of a function, before a return
func FuncEnd(logFields Fields) {
	withAddFields(logFields, Fields{"status": "end"}).Info()
}

// FuncFail logs fail, typically at the end of a function, before a return
func FuncFail(logFields Fields) {
	withAddFields(logFields, Fields{"status": "fail"}).Info()
}

// FuncSucc logs success, typically at the end of a function, before a return
func FuncSucc(logFields Fields) {
	withAddFields(logFields, Fields{"status": "success"}).Info()
}

// WithError adds an error field to the logrus entry
func WithError(err error) *logrus.Entry {
	return logrus.WithError(err)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, v ...interface{}) {
	logrus.Warningf(format, v...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, v ...interface{}) {
	logrus.Errorf(format, v...)
}

// Error logs a message at level Error on the standard logger.
func Error(v ...interface{}) {
	logrus.Error(v...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(v ...interface{}) {
	logrus.Warning(v...)
}

// Info logs a message at level Info on the standard logger.
func Info(v ...interface{}) {
	logrus.Info(v...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(v ...interface{}) {
	logrus.Debug(v...)
}

// there is no fatal on purpose - log and panic instead

// EnableJSONOutput enables JSON output on the logger
func EnableJSONOutput() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.AddHook(&contextHookObj)
}

// EnableTextOutput enables plain text output on the logger
func EnableTextOutput() {
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "Jan 02 03:04:05.000"})
}

// SetOutput sets the standard logger output.
func SetOutput(name string) {
	out, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		logrus.SetOutput(os.Stderr)
	}
	logrus.SetOutput(out)
}

// SetDebug sets the log level to debug
func SetDebug(on bool) {
	if on {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.AddHook(&contextHookObj)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

// SetWarn sets the log level to warn
func SetWarn() {
	logrus.SetLevel(logrus.WarnLevel)
}

// SetError sets the log level to error
func SetError() {
	logrus.SetLevel(logrus.ErrorLevel)
}

// SetReqID sets the request id that comes from bouncer
func SetReqID(reqID string) {
	contextHookObj.ReqID = reqID
}

// GetReqID returns the current request ID
func GetReqID() string {
	return contextHookObj.ReqID
}
