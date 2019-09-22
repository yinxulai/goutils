package logger

// globalService Config
var globalService *LoggerService

func init() {
	globalService = new(LoggerService)
}

func Debug(v ...interface{}) {
	globalService.Debug(v...)
}
func Debugf(format string, v ...interface{}) {
	globalService.Debugf(format, v)
}
func Info(v ...interface{}) {
	globalService.Info(v...)
}
func Infof(format string, v ...interface{}) {
	globalService.Infof(format, v...)
}
func Warn(v ...interface{}) {
	globalService.Warn(v...)
}
func Warnf(format string, v ...interface{}) {
	globalService.Warnf(format, v...)
}
func Error(v ...interface{}) {
	globalService.Error(v...)
}
func Errorf(format string, v ...interface{}) {
	globalService.Errorf(format, v...)
}
func SetLevel(level int) {
	globalService.SetLevel(level)
}
func SetOutPath(path string) {
	globalService.SetOutPath(path)
}
