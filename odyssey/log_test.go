// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package odyssey

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevelConstants(t *testing.T) {
	// Test that all level constants are properly defined
	assert.Equal(t, Level(0), LevelNull)
	assert.Equal(t, Level(1), LevelError)
	assert.Equal(t, Level(2), LevelWarn)
	assert.Equal(t, Level(3), LevelInfo)
	assert.Equal(t, Level(4), LevelDebug)
}

func TestDefaultLeveledLogger(t *testing.T) {
	// Test that the default logger is properly initialized
	assert.NotNil(t, DefaultLeveledLogger)

	// Test that it implements the interface
	var _ LeveledLoggerInterface = DefaultLeveledLogger

	// Test that it has the expected default level
	if logger, ok := DefaultLeveledLogger.(*LeveledLogger); ok {
		assert.Equal(t, LevelError, logger.Level)
	}
}

func TestLeveledLogger_Debugf(t *testing.T) {
	tests := []struct {
		name           string
		level          Level
		format         string
		args           []interface{}
		expectedOutput bool
		expectedPrefix string
	}{
		{
			name:           "Debug level with debug message",
			level:          LevelDebug,
			format:         "Debug message: %s",
			args:           []interface{}{"test"},
			expectedOutput: true,
			expectedPrefix: "[DEBUG]",
		},
		{
			name:           "Info level with debug message",
			level:          LevelInfo,
			format:         "Debug message: %s",
			args:           []interface{}{"test"},
			expectedOutput: false,
			expectedPrefix: "",
		},
		{
			name:           "Warn level with debug message",
			level:          LevelWarn,
			format:         "Debug message: %s",
			args:           []interface{}{"test"},
			expectedOutput: false,
			expectedPrefix: "",
		},
		{
			name:           "Error level with debug message",
			level:          LevelError,
			format:         "Debug message: %s",
			args:           []interface{}{"test"},
			expectedOutput: false,
			expectedPrefix: "",
		},
		{
			name:           "Null level with debug message",
			level:          LevelNull,
			format:         "Debug message: %s",
			args:           []interface{}{"test"},
			expectedOutput: false,
			expectedPrefix: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := &LeveledLogger{
				Level:          tt.level,
				stdoutOverride: &buf,
			}

			logger.Debugf(tt.format, tt.args...)
			output := buf.String()

			if tt.expectedOutput {
				assert.Contains(t, output, tt.expectedPrefix)
				assert.Contains(t, output, "Debug message: test")
			} else {
				assert.Empty(t, output)
			}
		})
	}
}

func TestLeveledLogger_Errorf(t *testing.T) {
	tests := []struct {
		name           string
		level          Level
		format         string
		args           []interface{}
		expectedOutput bool
		expectedPrefix string
	}{
		{
			name:           "Error level with error message",
			level:          LevelError,
			format:         "Error message: %s",
			args:           []interface{}{"test"},
			expectedOutput: true,
			expectedPrefix: "[ERROR]",
		},
		{
			name:           "Warn level with error message",
			level:          LevelWarn,
			format:         "Error message: %s",
			args:           []interface{}{"test"},
			expectedOutput: true,
			expectedPrefix: "[ERROR]",
		},
		{
			name:           "Info level with error message",
			level:          LevelInfo,
			format:         "Error message: %s",
			args:           []interface{}{"test"},
			expectedOutput: true,
			expectedPrefix: "[ERROR]",
		},
		{
			name:           "Debug level with error message",
			level:          LevelDebug,
			format:         "Error message: %s",
			args:           []interface{}{"test"},
			expectedOutput: true,
			expectedPrefix: "[ERROR]",
		},
		{
			name:           "Null level with error message",
			level:          LevelNull,
			format:         "Error message: %s",
			args:           []interface{}{"test"},
			expectedOutput: false,
			expectedPrefix: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := &LeveledLogger{
				Level:          tt.level,
				stderrOverride: &buf,
			}

			logger.Errorf(tt.format, tt.args...)
			output := buf.String()

			if tt.expectedOutput {
				assert.Contains(t, output, tt.expectedPrefix)
				assert.Contains(t, output, "Error message: test")
			} else {
				assert.Empty(t, output)
			}
		})
	}
}

func TestLeveledLogger_Infof(t *testing.T) {
	tests := []struct {
		name           string
		level          Level
		format         string
		args           []interface{}
		expectedOutput bool
		expectedPrefix string
	}{
		{
			name:           "Info level with info message",
			level:          LevelInfo,
			format:         "Info message: %s",
			args:           []interface{}{"test"},
			expectedOutput: true,
			expectedPrefix: "[INFO]",
		},
		{
			name:           "Debug level with info message",
			level:          LevelDebug,
			format:         "Info message: %s",
			args:           []interface{}{"test"},
			expectedOutput: true,
			expectedPrefix: "[INFO]",
		},
		{
			name:           "Warn level with info message",
			level:          LevelWarn,
			format:         "Info message: %s",
			args:           []interface{}{"test"},
			expectedOutput: false,
			expectedPrefix: "",
		},
		{
			name:           "Error level with info message",
			level:          LevelError,
			format:         "Info message: %s",
			args:           []interface{}{"test"},
			expectedOutput: false,
			expectedPrefix: "",
		},
		{
			name:           "Null level with info message",
			level:          LevelNull,
			format:         "Info message: %s",
			args:           []interface{}{"test"},
			expectedOutput: false,
			expectedPrefix: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := &LeveledLogger{
				Level:          tt.level,
				stdoutOverride: &buf,
			}

			logger.Infof(tt.format, tt.args...)
			output := buf.String()

			if tt.expectedOutput {
				assert.Contains(t, output, tt.expectedPrefix)
				assert.Contains(t, output, "Info message: test")
			} else {
				assert.Empty(t, output)
			}
		})
	}
}

func TestLeveledLogger_Warnf(t *testing.T) {
	tests := []struct {
		name           string
		level          Level
		format         string
		args           []interface{}
		expectedOutput bool
		expectedPrefix string
	}{
		{
			name:           "Warn level with warn message",
			level:          LevelWarn,
			format:         "Warn message: %s",
			args:           []interface{}{"test"},
			expectedOutput: true,
			expectedPrefix: "[WARN]",
		},
		{
			name:           "Info level with warn message",
			level:          LevelInfo,
			format:         "Warn message: %s",
			args:           []interface{}{"test"},
			expectedOutput: true,
			expectedPrefix: "[WARN]",
		},
		{
			name:           "Debug level with warn message",
			level:          LevelDebug,
			format:         "Warn message: %s",
			args:           []interface{}{"test"},
			expectedOutput: true,
			expectedPrefix: "[WARN]",
		},
		{
			name:           "Error level with warn message",
			level:          LevelError,
			format:         "Warn message: %s",
			args:           []interface{}{"test"},
			expectedOutput: false,
			expectedPrefix: "",
		},
		{
			name:           "Null level with warn message",
			level:          LevelNull,
			format:         "Warn message: %s",
			args:           []interface{}{"test"},
			expectedOutput: false,
			expectedPrefix: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := &LeveledLogger{
				Level:          tt.level,
				stderrOverride: &buf,
			}

			logger.Warnf(tt.format, tt.args...)
			output := buf.String()

			if tt.expectedOutput {
				assert.Contains(t, output, tt.expectedPrefix)
				assert.Contains(t, output, "Warn message: test")
			} else {
				assert.Empty(t, output)
			}
		})
	}
}

func TestLeveledLogger_stderr(t *testing.T) {
	t.Run("with stderr override", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &LeveledLogger{
			stderrOverride: &buf,
		}

		writer := logger.stderr()
		assert.Equal(t, &buf, writer)
	})

	t.Run("without stderr override", func(t *testing.T) {
		logger := &LeveledLogger{}

		writer := logger.stderr()
		assert.Equal(t, os.Stderr, writer)
	})
}

func TestLeveledLogger_stdout(t *testing.T) {
	t.Run("with stdout override", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &LeveledLogger{
			stdoutOverride: &buf,
		}

		writer := logger.stdout()
		assert.Equal(t, &buf, writer)
	})

	t.Run("without stdout override", func(t *testing.T) {
		logger := &LeveledLogger{}

		writer := logger.stdout()
		assert.Equal(t, os.Stdout, writer)
	})
}

func TestLeveledLogger_InterfaceCompliance(t *testing.T) {
	// Test that LeveledLogger implements LeveledLoggerInterface
	var _ LeveledLoggerInterface = &LeveledLogger{}

	// Test that we can call all interface methods
	logger := &LeveledLogger{Level: LevelDebug}

	// These should not panic
	logger.Debugf("test")
	logger.Infof("test")
	logger.Warnf("test")
	logger.Errorf("test")
}

func TestLeveledLogger_OutputDestinations(t *testing.T) {
	t.Run("Debug and Info go to stdout", func(t *testing.T) {
		var stdoutBuf, stderrBuf bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelDebug,
			stdoutOverride: &stdoutBuf,
			stderrOverride: &stderrBuf,
		}

		logger.Debugf("debug message")
		logger.Infof("info message")

		assert.Contains(t, stdoutBuf.String(), "debug message")
		assert.Contains(t, stdoutBuf.String(), "info message")
		assert.Empty(t, stderrBuf.String())
	})

	t.Run("Warn and Error go to stderr", func(t *testing.T) {
		var stdoutBuf, stderrBuf bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelDebug,
			stdoutOverride: &stdoutBuf,
			stderrOverride: &stderrBuf,
		}

		logger.Warnf("warn message")
		logger.Errorf("error message")

		assert.Contains(t, stderrBuf.String(), "warn message")
		assert.Contains(t, stderrBuf.String(), "error message")
		assert.Empty(t, stdoutBuf.String())
	})
}

func TestLeveledLogger_Formatting(t *testing.T) {
	t.Run("Simple formatting", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelDebug,
			stdoutOverride: &buf,
		}

		logger.Debugf("Simple message")
		output := buf.String()

		assert.Contains(t, output, "[DEBUG]")
		assert.Contains(t, output, "Simple message")
		assert.True(t, strings.HasSuffix(output, "\n"))
	})

	t.Run("Complex formatting", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelInfo,
			stdoutOverride: &buf,
		}

		logger.Infof("User %s has %d items", "john", 42)
		output := buf.String()

		assert.Contains(t, output, "[INFO]")
		assert.Contains(t, output, "User john has 42 items")
		assert.True(t, strings.HasSuffix(output, "\n"))
	})

	t.Run("Multiple arguments", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelWarn,
			stderrOverride: &buf,
		}

		logger.Warnf("Values: %s, %d, %t, %f", "test", 123, true, 3.14)
		output := buf.String()

		assert.Contains(t, output, "[WARN]")
		assert.Contains(t, output, "Values: test, 123, true, 3.140000")
		assert.True(t, strings.HasSuffix(output, "\n"))
	})
}

func TestLeveledLogger_EdgeCases(t *testing.T) {
	t.Run("Empty format string", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelDebug,
			stdoutOverride: &buf,
		}

		logger.Debugf("")
		output := buf.String()

		assert.Contains(t, output, "[DEBUG]")
		assert.True(t, strings.HasSuffix(output, "\n"))
	})

	t.Run("No arguments", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelInfo,
			stdoutOverride: &buf,
		}

		logger.Infof("No arguments")
		output := buf.String()

		assert.Contains(t, output, "[INFO]")
		assert.Contains(t, output, "No arguments")
	})

	t.Run("Nil arguments", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelError,
			stderrOverride: &buf,
		}

		logger.Errorf("Nil: %v", nil)
		output := buf.String()

		assert.Contains(t, output, "[ERROR]")
		assert.Contains(t, output, "Nil: <nil>")
	})
}

func TestLeveledLogger_LevelComparison(t *testing.T) {
	// Test basic level comparisons
	assert.True(t, LevelDebug >= LevelInfo)
	assert.True(t, LevelInfo >= LevelWarn)
	assert.True(t, LevelWarn >= LevelError)
	assert.True(t, LevelError >= LevelNull)

	assert.False(t, LevelNull >= LevelError)
	assert.False(t, LevelError >= LevelWarn)
	assert.False(t, LevelWarn >= LevelInfo)
	assert.False(t, LevelInfo >= LevelDebug)

	// Test that levels are properly ordered
	assert.Equal(t, Level(0), LevelNull)
	assert.Equal(t, Level(1), LevelError)
	assert.Equal(t, Level(2), LevelWarn)
	assert.Equal(t, Level(3), LevelInfo)
	assert.Equal(t, Level(4), LevelDebug)
}

func TestLeveledLogger_CustomWriter(t *testing.T) {
	t.Run("Custom writer for stdout", func(t *testing.T) {
		var customWriter bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelDebug,
			stdoutOverride: &customWriter,
		}

		logger.Debugf("Custom stdout test")
		output := customWriter.String()

		assert.Contains(t, output, "[DEBUG]")
		assert.Contains(t, output, "Custom stdout test")
	})

	t.Run("Custom writer for stderr", func(t *testing.T) {
		var customWriter bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelError,
			stderrOverride: &customWriter,
		}

		logger.Errorf("Custom stderr test")
		output := customWriter.String()

		assert.Contains(t, output, "[ERROR]")
		assert.Contains(t, output, "Custom stderr test")
	})
}

func TestLeveledLogger_ConcurrentAccess(t *testing.T) {
	// Test that the logger can handle concurrent access
	logger := &LeveledLogger{Level: LevelDebug}

	// Create multiple goroutines that log simultaneously
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// These should not panic
			logger.Debugf("Concurrent message %d", id)
			logger.Infof("Concurrent info %d", id)
			logger.Warnf("Concurrent warn %d", id)
			logger.Errorf("Concurrent error %d", id)
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestLeveledLogger_InterfaceMethods(t *testing.T) {
	// Test that all interface methods work correctly
	var logger LeveledLoggerInterface = &LeveledLogger{Level: LevelDebug}

	// These should not panic and should work correctly
	logger.Debugf("Interface debug")
	logger.Infof("Interface info")
	logger.Warnf("Interface warn")
	logger.Errorf("Interface error")
}

func TestLeveledLogger_LevelBoundaries(t *testing.T) {
	t.Run("Exact level match", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelInfo,
			stdoutOverride: &buf,
		}

		logger.Infof("Exact level test")
		output := buf.String()

		assert.Contains(t, output, "[INFO]")
		assert.Contains(t, output, "Exact level test")
	})

	t.Run("Level above threshold", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelWarn,
			stdoutOverride: &buf,
		}

		logger.Infof("Above threshold test")
		output := buf.String()

		// Should not output because Info level is below Warn level
		assert.Empty(t, output)
	})
}

func TestLeveledLogger_NewlineHandling(t *testing.T) {
	t.Run("Message with newline", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelDebug,
			stdoutOverride: &buf,
		}

		logger.Debugf("Message with\nnewline")
		output := buf.String()

		assert.Contains(t, output, "[DEBUG]")
		assert.Contains(t, output, "Message with")
		assert.Contains(t, output, "newline")
		assert.True(t, strings.HasSuffix(output, "\n"))
	})

	t.Run("Message without newline", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &LeveledLogger{
			Level:          LevelDebug,
			stdoutOverride: &buf,
		}

		logger.Debugf("Message without newline")
		output := buf.String()

		assert.Contains(t, output, "[DEBUG]")
		assert.Contains(t, output, "Message without newline")
		assert.True(t, strings.HasSuffix(output, "\n"))
	})
}
