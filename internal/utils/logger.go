package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger for structured, production-grade logging to file.
type Logger struct {
	zlog    zerolog.Logger
	logFile *os.File
}

// NewLogger creates a new Logger instance with file rotation and JSON formatting.
// Logs are written to ~/.kavach/kavach.log by default.
func NewLogger(cfg *config.Config) *Logger {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Could not determine user home directory: " + err.Error())
	}
	// Clean up the log dir path (remove leading slash if present)
	logDir := strings.TrimPrefix(cfg.LogDirPath, "/")
	// Join home directory and log directory
	fullLogDir := filepath.Join(homeDir, logDir)
	// Ensure the directory exists
	if err := os.MkdirAll(fullLogDir, 0700); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}
	// Set the log file path
	logPath := filepath.Join(fullLogDir, "kavach.log")
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    cfg.LogMaxSize,    // megabytes
		MaxBackups: cfg.LogMaxBackups, // number of old files to keep
		MaxAge:     cfg.LogMaxAge,     // days
		Compress:   cfg.LogCompress,   // whether to compress old files
	}
	z := zerolog.New(lumberjackLogger).With().Timestamp().Logger()
	return &Logger{zlog: z, logFile: nil}
}

// Info logs an informational message with optional structured fields.
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.zlog.Info().Fields(fields[0]).Msg(msg)
	} else {
		l.zlog.Info().Msg(msg)
	}
}

// Warn logs a warning message with optional structured fields.
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.zlog.Warn().Fields(fields[0]).Msg(msg)
	} else {
		l.zlog.Warn().Msg(msg)
	}
}

// Error logs an error message with an error object and optional structured fields.
func (l *Logger) Error(msg string, err error, fields ...map[string]interface{}) {
	le := l.zlog.Error()
	if err != nil {
		le = le.Err(err)
	}
	if len(fields) > 0 {
		le.Fields(fields[0]).Msg(msg)
	} else {
		le.Msg(msg)
	}
}

// Debug logs a debug message with optional structured fields.
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.zlog.Debug().Fields(fields[0]).Msg(msg)
	} else {
		l.zlog.Debug().Msg(msg)
	}
}

// Fatal logs a fatal error message and exits the application.
func (l *Logger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	le := l.zlog.Fatal()
	if err != nil {
		le = le.Err(err)
	}
	if len(fields) > 0 {
		le.Fields(fields[0]).Msg(msg)
	} else {
		le.Msg(msg)
	}
	os.Exit(1)
}

// Print logs a simple info message (for compatibility with Print-style usage).
func (l *Logger) Print(msg string) {
	l.zlog.Info().Msg(msg)
}

// Printf logs a formatted info message (for compatibility with Printf-style usage).
func (l *Logger) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.zlog.Info().Msg(msg)
}

// LogRequest logs an HTTP request string for debugging/audit.
func (l *Logger) LogRequest(req string) {
	l.zlog.Info().Str("request", req).Msg("HTTP request")
}

// LogResponse logs an HTTP response string for debugging/audit.
func (l *Logger) LogResponse(resp string) {
	l.zlog.Info().Str("response", resp).Msg("HTTP response")
}

// Close closes the log file if needed (no-op for stdout).
func (l *Logger) Close() {
	if l.logFile != nil && l.logFile != os.Stdout {
		l.logFile.Close()
	}
}
