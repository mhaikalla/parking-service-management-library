package logs

import (
	"strings"
	"sync"
	"time"

	"parking-service/pkg/condutils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// LogContext contains all thing to logging.
type LogrusContext struct {
	name           string
	ID             string
	baggage        map[string]interface{}
	parentID       string
	childs         []string
	beforeUpdateFn func(baggage map[string]interface{}) map[string]interface{}
	lock           sync.Mutex
	*logrus.Entry
}

// BeforeUpdateFn func called before updating log fields.
// Receive baggage of context and return the updated baggage.
// Log's child will inheritance parent's func.
type BeforeUpdateFn func(baggage map[string]interface{}) map[string]interface{}

// defaultBeforeUpdatefn default func.
func defaultBeforeUpdatefn(baggage map[string]interface{}) map[string]interface{} {
	return baggage
}

// ILog interface to exposed log method
type ILog interface {
	Upsert(key string, value interface{})
	Update()
	Child(name string) ILog
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	BeforeUpdate(fn BeforeUpdateFn)
}

// Upsert add key value pairs to a context baggage for populate log field.
// After adding some value, the logger will not updated without calling `Update`.
func (lc *LogrusContext) Upsert(key string, value interface{}) {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	lc.baggage[key] = value
}

// Update add key value pairs from context to log field.
// This not mean to start logging, need to call on of log method eg: `Info`, `Warn`, etc.
func (lc *LogrusContext) Update() {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	trans := lc.beforeUpdateFn(lc.baggage)
	lc.Entry = lc.WithFields(trans).WithField(LogChilds, lc.childs)
}

// BeforeUpdate add a hook before updating/commit key-pairs of data via `logs#LogrusContext`.
func (lc *LogrusContext) BeforeUpdate(fn BeforeUpdateFn) {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	lc.beforeUpdateFn = fn
}

// Child create a new logger with all fields own by this logger.
// If `name` is empty string, use parent `name`.
func (lc *LogrusContext) Child(name string) ILog {
	lc.Update()
	lc.lock.Lock()
	defer lc.lock.Unlock()
	id := uuid.New().String()

	lc.childs = append(lc.childs, id)

	logger := lc.WithFields(lc.baggage).WithFields(map[string]interface{}{
		LogID:     id,
		LogName:   condutils.Or(name, lc.name),
		LogType:   condutils.Or(name, lc.name),
		LogParent: lc.ID,
		LogChilds: []string{},
	})

	logCtx := &LogrusContext{
		ID:             id,
		parentID:       lc.ID,
		name:           lc.name,
		Entry:          logger,
		baggage:        map[string]interface{}{LogCreatedAt: time.Now().Local().Format(time.RFC3339), SystemName: ServiceName},
		childs:         []string{},
		beforeUpdateFn: lc.beforeUpdateFn,
	}

	return logCtx
}

// configureLog configure logger log level and output format.
func configureLog(logger *logrus.Entry, level, format string) {
	level = strings.ToLower(level)
	format = strings.ToLower(format)
	switch level {
	default:
		logger.Logger.SetLevel(logrus.ErrorLevel)
	case "trace":
		logger.Logger.SetLevel(logrus.TraceLevel)
	case "debug":
		logger.Logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.Logger.SetLevel(logrus.WarnLevel)
	case "fatal":
		logger.Logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.Logger.SetLevel(logrus.PanicLevel)
	}

	if format == "json" {
		logger.Logger.SetFormatter(&logrus.JSONFormatter{})
		return
	}
	logger.Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		QuoteEmptyFields: true,
	})
}

// NewLogrus create a new logrus wrapper implemented `ILog`.
func NewLogrus(name string) ILog {
	id := uuid.New().String()
	logger := logrus.New().WithFields(map[string]interface{}{
		LogID:   id,
		LogName: name,
		LogType: name,
	})

	configureLog(logger, LogLevel, LogFormat)

	return &LogrusContext{
		ID:             id,
		name:           name,
		Entry:          logger,
		baggage:        map[string]interface{}{LogCreatedAt: time.Now().Local().Format(time.RFC3339), SystemName: ServiceName},
		beforeUpdateFn: defaultBeforeUpdatefn,
	}
}
