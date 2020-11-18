// +build !jsonpack_debug

package logger

// Trace outputs nothing in null loagger, it emits any log output in production environment.
func Trace(args ...interface{}) {}

// Tracef outputs nothing in null loagger, it emits any log output in production environment.
func Tracef(format string, args ...interface{}) {}

// Debug outputs nothing in null loagger, it emits any log output in production environment.
func Debug(args ...interface{}) {}

// Debugf outputs nothing in null loagger, it emits any log output in production environment.
func Debugf(format string, args ...interface{}) {}

// Info outputs nothing in null loagger, it emits any log output in production environment.
func Info(args ...interface{}) {}

// Infof outputs nothing in null loagger, it emits any log output in production environment.
func Infof(format string, args ...interface{}) {}

// Warn outputs nothing in null loagger, it emits any log output in production environment.
func Warn(args ...interface{}) {}

// Warnf outputs nothing in null loagger, it emits any log output in production environment.
func Warnf(format string, args ...interface{}) {}

// Error outputs nothing in null loagger, it emits any log output in production environment.
func Error(args ...interface{}) {}

// Errorf outputs nothing in null loagger, it emits any log output in production environment.
func Errorf(format string, args ...interface{}) {}

// Fatal outputs nothing in null loagger, it emits any log output in production environment.
func Fatal(args ...interface{}) {}

// Fatalf outputs nothing in null loagger, it emits any log output in production environment.
func Fatalf(format string, args ...interface{}) {}

// Panic outputs nothing in null loagger, it emits any log output in production environment.
func Panic(args ...interface{}) {}

// Panicf outputs nothing in null loagger, it emits any log output in production environment.
func Panicf(format string, args ...interface{}) {}
