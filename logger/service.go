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

// Service Service
type Service struct {
	level       int
	outPath     string
	history     []string
	warnLogger  *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

func (c *Service) output(v string) {
	if c.history == nil {
		c.history = make([]string, 200)
	}
	c.history = append(c.history, v)

}

// Debug Debug
func (c *Service) Debug(v ...interface{}) {
	if c.level >= DebugLevel {
		c.debugLogger.Print(v)
	}
}

// Debugf Debugf
func (c *Service) Debugf(format string, v ...interface{}) {
	if c.level >= DebugLevel {
		c.debugLogger.Printf(format, v)
	}
}

// Info Info
func (c *Service) Info(v ...interface{}) {
	if c.level >= DebugLevel {
		c.infoLogger.Print(v)
	}
}

// Infof Infof
func (c *Service) Infof(format string, v ...interface{}) {
	if c.level >= DebugLevel {
		c.infoLogger.Printf(format, v)
	}
}

// Warn Warn
func (c *Service) Warn(v ...interface{}) {
	if c.level >= DebugLevel {
		c.warnLogger.Print(v)
	}
}

// Warnf Warnf
func (c *Service) Warnf(format string, v ...interface{}) {
	if c.level >= DebugLevel {
		c.warnLogger.Printf(format, v)
	}
}

// Error Error
func (c *Service) Error(v ...interface{}) {
	if c.level >= DebugLevel {
		c.errorLogger.Print(v)
	}
}

// Errorf Errorf
func (c *Service) Errorf(format string, v ...interface{}) {
	if c.level >= DebugLevel {
		c.errorLogger.Printf(format, v)
	}
}

// SetLevel 设置输出级别
func (c *Service) SetLevel(level int) {
	c.level = level
}

// SetOutPath 设置文件输出位置
func (c *Service) SetOutPath(path string) {
	c.outPath = path
}
