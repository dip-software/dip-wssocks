package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// Fields type, used similarly to logrus.Fields
type Fields map[string]interface{}

// Logger is a wrapper for slog.Logger to provide a similar API to logrus
type Logger struct {
	logger *slog.Logger
}

// Default logger instance
var defaultLogger = &Logger{
	logger: slog.Default(),
}

// New creates a new logger instance
func New() *Logger {
	return &Logger{
		logger: slog.Default(),
	}
}

// SetLevel sets the log level
func SetLevel(level slog.Level) {
	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	defaultLogger.logger = slog.New(handler)
	slog.SetDefault(defaultLogger.logger)
}

// WithFields returns a new Logger with the given fields added to the context
func WithFields(fields Fields) *Logger {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}
	
	return &Logger{
		logger: defaultLogger.logger.With(attrs...),
	}
}

// WithField returns a new Logger with the given field added to the context
func WithField(key string, value interface{}) *Logger {
	return &Logger{
		logger: defaultLogger.logger.With(key, value),
	}
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		logger: l.logger.With(key, value),
	}
}

// WithFields adds fields to the logger
func (l *Logger) WithFields(fields Fields) *Logger {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}
	
	return &Logger{
		logger: l.logger.With(attrs...),
	}
}

// Trace logs a message at level Trace
func (l *Logger) Trace(args ...interface{}) {
	l.logger.Debug("TRACE", "msg", args)
}

// Debug logs a message at level Debug
func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(formatArgs(args...))
}

// Info logs a message at level Info
func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(formatArgs(args...))
}

// Warn logs a message at level Warn
func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(formatArgs(args...))
}

// Error logs a message at level Error
func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(formatArgs(args...))
}

// Fatal logs a message at level Fatal then calls os.Exit(1)
func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Error(formatArgs(args...), slog.String("level", "FATAL"))
	os.Exit(1)
}

// Tracef logs a formatted message at level Trace
func (l *Logger) Tracef(format string, args ...interface{}) {
	// slog doesn't have Trace level, so we use Debug with a TRACE prefix
	ctx := context.Background()
	l.logger.DebugContext(ctx, "TRACE: "+format, args...)
}

// Debugf logs a formatted message at level Debug
func (l *Logger) Debugf(format string, args ...interface{}) {
	ctx := context.Background()
	l.logger.DebugContext(ctx, format, args...)
}

// Infof logs a formatted message at level Info
func (l *Logger) Infof(format string, args ...interface{}) {
	ctx := context.Background()
	l.logger.InfoContext(ctx, format, args...)
}

// Warnf logs a formatted message at level Warn
func (l *Logger) Warnf(format string, args ...interface{}) {
	ctx := context.Background()
	l.logger.WarnContext(ctx, format, args...)
}

// Errorf logs a formatted message at level Error
func (l *Logger) Errorf(format string, args ...interface{}) {
	ctx := context.Background()
	l.logger.ErrorContext(ctx, format, args...)
}

// Fatalf logs a formatted message at level Fatal then calls os.Exit(1)
func (l *Logger) Fatalf(format string, args ...interface{}) {
	ctx := context.Background()
	l.logger.ErrorContext(ctx, format, slog.String("level", "FATAL"), slog.Any("args", args))
	os.Exit(1)
}

// Helper function to format interface args into a string
func formatArgs(args ...interface{}) string {
	if len(args) == 1 {
		if str, ok := args[0].(string); ok {
			return str
		}
	}
	
	// Just return the first arg as a string if possible
	if len(args) > 0 {
		if str, ok := args[0].(string); ok {
			return str
		}
		return "message"
	}
	
	return "message"
}

// Global functions that mimic the logrus global API

// Trace logs a message at level Trace on the default logger
func Trace(args ...interface{}) {
	defaultLogger.Trace(args...)
}

// Debug logs a message at level Debug on the default logger
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Info logs a message at level Info on the default logger
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Warn logs a message at level Warn on the default logger
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Error logs a message at level Error on the default logger
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Fatal logs a message at level Fatal on the default logger then calls os.Exit(1)
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Tracef logs a formatted message at level Trace on the default logger
func Tracef(format string, args ...interface{}) {
	defaultLogger.Tracef(format, args...)
}

// Debugf logs a formatted message at level Debug on the default logger
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Infof logs a formatted message at level Info on the default logger
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Warnf logs a formatted message at level Warn on the default logger
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Errorf logs a formatted message at level Error on the default logger
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Infoln logs a message at level Info with a newline
func (l *Logger) Infoln(args ...interface{}) {
	l.logger.Info(formatArgs(args...) + "\n")
}

// Warnln logs a message at level Warn with a newline
func (l *Logger) Warnln(args ...interface{}) {
	l.logger.Warn(formatArgs(args...) + "\n")
}

// Warningln logs a message at level Warn with a newline (alias for Warnln)
func (l *Logger) Warningln(args ...interface{}) {
	l.Warnln(args...)
}

// Warning logs a message at level Warn (alias for Warn)
func (l *Logger) Warning(args ...interface{}) {
	l.Warn(args...)
}

// Traceln logs a message at level Trace with a newline
func (l *Logger) Traceln(args ...interface{}) {
	ctx := context.Background()
	l.logger.DebugContext(ctx, "TRACE: "+formatArgs(args...)+"\n")
}

// Errorln logs a message at level Error with a newline
func (l *Logger) Errorln(args ...interface{}) {
	l.logger.Error(formatArgs(args...) + "\n")
}

// Println logs a message at level Info with a newline
func (l *Logger) Println(args ...interface{}) {
	l.logger.Info(formatArgs(args...) + "\n")
}

// SetOutput sets the output destination for the logger
func SetOutput(w io.Writer) {
	// Create a new handler with the provided writer
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(w, opts)
	defaultLogger.logger = slog.New(handler)
	slog.SetDefault(defaultLogger.logger)
}

// Global functions for ln variants

// Infoln logs a message at level Info with a newline on the default logger
func Infoln(args ...interface{}) {
	defaultLogger.Infoln(args...)
}

// Warnln logs a message at level Warn with a newline on the default logger
func Warnln(args ...interface{}) {
	defaultLogger.Warnln(args...)
}

// Warningln logs a message at level Warn with a newline on the default logger (alias for Warnln)
func Warningln(args ...interface{}) {
	defaultLogger.Warnln(args...)
}

// Warning logs a message at level Warn on the default logger (alias for Warn)
func Warning(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Traceln logs a message at level Trace with a newline on the default logger
func Traceln(args ...interface{}) {
	defaultLogger.Traceln(args...)
}

// Errorln logs a message at level Error with a newline on the default logger
func Errorln(args ...interface{}) {
	defaultLogger.Errorln(args...)
}

// Println logs a message at level Info with a newline on the default logger
func Println(args ...interface{}) {
	defaultLogger.Println(args...)
}
