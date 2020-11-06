// +build !jsonpack_debug

package logger

func Trace(args ...interface{})                 {}
func Tracef(format string, args ...interface{}) {}
func Debug(args ...interface{})                 {}
func Debugf(format string, args ...interface{}) {}
func Info(args ...interface{})                  {}
func Infof(format string, args ...interface{})  {}
func Warn(args ...interface{})                  {}
func Warnf(format string, args ...interface{})  {}
func Error(args ...interface{})                 {}
func Errorf(format string, args ...interface{}) {}
func Fatal(args ...interface{})                 {}
func Fatalf(format string, args ...interface{}) {}
func Panic(args ...interface{})                 {}
func Panicf(format string, args ...interface{}) {}
