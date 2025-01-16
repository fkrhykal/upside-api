package log

import (
	"testing"
)

type TestLogger struct {
	t testing.TB
}

func (log *TestLogger) Infof(format string, args ...any) {
	log.t.Logf("INFO: "+format, args...)
}

func (log *TestLogger) Debugf(format string, args ...any) {
	log.t.Logf("DEBUG: "+format, args...)
}

func (log *TestLogger) Warnf(format string, args ...any) {
	log.t.Logf("WARN: "+format, args...)
}

func (log *TestLogger) Errorf(format string, args ...any) {
	log.t.Logf("ERROR: "+format, args...)
}

func (log *TestLogger) Info(args ...any) {
	log.t.Logf("INFO: %+v", args...)
}

func (log *TestLogger) Debug(args ...any) {
	log.t.Logf("DEBUG: %+v", args...)
}

func NewTestLogger(t testing.TB) Logger {
	return &TestLogger{
		t: t,
	}
}
