package logs

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_configureLog(t *testing.T) {
	type args struct {
		logger *logrus.Entry
		level  string
		format string
	}
	tests := []struct {
		name string
		args args
	}{
		{"1", args{logger: logrus.WithField("a", "b")}},
		{"2", args{logger: logrus.WithField("a", "b"), format: "text"}},
		{"3", args{logger: logrus.WithField("a", "b"), format: "json"}},
		{"4", args{logger: logrus.WithField("a", "b"), level: "trace"}},
		{"5", args{logger: logrus.WithField("a", "b"), level: "info"}},
		{"6", args{logger: logrus.WithField("a", "b"), level: "debug"}},
		{"7", args{logger: logrus.WithField("a", "b"), level: "warn"}},
		{"8", args{logger: logrus.WithField("a", "b"), level: "panic"}},
		{"9", args{logger: logrus.WithField("a", "b"), level: "fatal"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configureLog(tt.args.logger, tt.args.level, tt.args.format)
		})
	}
}

func TestNewLogrus(t *testing.T) {
	logger := NewLogrus("TEST")
	if logger == nil {
		t.Fail()
	}

	logger.Upsert("a", "b")
	logger.Upsert("c", "d")
	logger.Update()
	child := logger.Child("Some")
	logger.Error("parent will logging all")
	child.Info("child will inform info")

	WhenError("logging", errors.New("test"), logger)
	WhenError("no logging", nil, logger)
	HandlerWhenError(errors.New("test"), logger)

}

func TestHooks(t *testing.T) {
	logger := NewLogrus("TEST")
	if logger == nil {
		t.Fail()
	}

	customBeforeUpdate := func(bag map[string]interface{}) map[string]interface{} {
		delete(bag, "a")
		return bag
	}

	logger.Upsert("a", "b")
	logger.Upsert("c", "d")
	logger.BeforeUpdate(customBeforeUpdate)
	logger.Update()
	child := logger.Child("Some")
	logger.Error("parent will logging all")
	child.Info("child will inform info")

	WhenError("logging", errors.New("test"), logger)
	WhenError("no logging", nil, logger)
	HandlerWhenError(errors.New("test"), logger)
}
