package logger

// defaultService Config
var defaultService *Service

func init() {
	defaultService = new(Service)
}

func Debug(v ...interface{}) {
	defaultService.Debug(v...)
}
func Debugf(format string, v ...interface{}) {
	defaultService.Debugf(format, v)
}
func Info(v ...interface{}) {
	defaultService.Info(v...)
}
func Infof(format string, v ...interface{}) {
	defaultService.Infof(format, v...)
}
func Warn(v ...interface{}) {
	defaultService.Warn(v...)
}
func Warnf(format string, v ...interface{}) {
	defaultService.Warnf(format, v...)
}
func Error(v ...interface{}) {
	defaultService.Error(v...)
}
func Errorf(format string, v ...interface{}) {
	defaultService.Errorf(format, v...)
}
func SetLevel(level int) {
	defaultService.SetLevel(level)
}
func SetOutPath(path string) {
	defaultService.SetOutPath(path)
}
