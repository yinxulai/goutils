package logger

import (
	"log"
)

// 日志级别
const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// Logger 接口类型
type Logger interface {
	SetLevel(level int)
	SetOutPath(path string)
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

// LoggerService LoggerService
type LoggerService struct {
	level       int
	outPath     string
	history     []string
	warnLogger  *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

func (c *LoggerService) output(v string) {
	if c.history == nil {
		c.history = make([]string, 200)
	}
	c.history = append(c.history, v)
}

// Debug Debug
func (c *LoggerService) Debug(v ...interface{}) {
	if c.level >= DebugLevel {
		c.debugLogger.Print(v...)
	}
}

// Debugf Debugf
func (c *LoggerService) Debugf(format string, v ...interface{}) {
	if c.level >= DebugLevel {
		c.debugLogger.Printf(format, v...)
	}
}

// Info Info
func (c *LoggerService) Info(v ...interface{}) {
	if c.level >= DebugLevel {
		c.infoLogger.Print(v...)
	}
}

// Infof Infof
func (c *LoggerService) Infof(format string, v ...interface{}) {
	if c.level >= DebugLevel {
		c.infoLogger.Printf(format, v...)
	}
}

// Warn Warn
func (c *LoggerService) Warn(v ...interface{}) {
	if c.level >= DebugLevel {
		c.warnLogger.Print(v...)
	}
}

// Warnf Warnf
func (c *LoggerService) Warnf(format string, v ...interface{}) {
	if c.level >= DebugLevel {
		c.warnLogger.Printf(format, v...)
	}
}

// Error Error
func (c *LoggerService) Error(v ...interface{}) {
	if c.level >= DebugLevel {
		c.errorLogger.Print(v...)
	}
}

// Errorf Errorf
func (c *LoggerService) Errorf(format string, v ...interface{}) {
	if c.level >= DebugLevel {
		c.errorLogger.Printf(format, v...)
	}
}

// SetLevel 设置输出级别
func (c *LoggerService) SetLevel(level int) {
	c.level = level
}

// SetOutPath 设置文件输出位置
func (c *LoggerService) SetOutPath(path string) {
	c.outPath = path
}
