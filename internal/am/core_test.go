package am

import (
	"testing"
)

func TestBaseCoreDefaultLogger(t *testing.T) {
	core := &BaseCore{name: "test"}
	
	logger := core.Log()
	if logger == nil {
		t.Error("Expected default logger, got nil")
	}
	
	if _, ok := logger.(*BaseLogger); !ok {
		t.Error("Expected BaseLogger instance")
	}
	
	logger.Info("Test message from default logger")
}

func TestBaseCoreDefaultLoggerSingleton(t *testing.T) {
	core1 := &BaseCore{name: "test1"}
	core2 := &BaseCore{name: "test2"}
	
	logger1 := core1.Log()
	logger2 := core2.Log()
	
	if logger1 != logger2 {
		t.Error("Expected same default logger instance (singleton)")
	}
}

func TestBaseCoreCustomLogger(t *testing.T) {
	customLogger := NewLogger("debug")
	core := &BaseCore{name: "test"}
	core.SetLog(customLogger)
	
	returnedLogger := core.Log()
	if returnedLogger != customLogger {
		t.Error("Expected custom logger, got different instance")
	}
}