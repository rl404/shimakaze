package utils

import (
	_log "github.com/rl404/fairy/log"
	"github.com/rl404/fairy/log/chain"
	"github.com/rl404/shimakaze/pkg/log"
)

var l _log.Logger
var ls []_log.Logger

// InitLog to init global logger.
func InitLog(lvl log.LogLevel, json, color bool) (err error) {
	l, err = log.New(log.Config{
		Type:       log.Zerolog,
		Level:      lvl,
		JsonFormat: json,
		Color:      color,
	})
	if err != nil {
		return err
	}
	ls = []_log.Logger{l}
	return nil
}

// AddLog to add logger chain.
func AddLog(l1 _log.Logger) {
	l = chain.New(l, l1)
	ls = append(ls, l1)
}

// GetLogger to get logger.
func GetLogger(i ...int) _log.Logger {
	if len(i) > 0 {
		if len(ls) <= i[0] {
			tmp, _ := log.New(log.Config{Type: log.NOP})
			return tmp
		}
		return ls[i[0]]
	}
	return l
}

// Log to log with custom field.
func Log(field map[string]interface{}) {
	GetLogger().Log(field)
}

// Info to log info.
func Info(str string, args ...interface{}) {
	GetLogger().Info(str, args...)
}

// Error to print error.
func Error(str string, args ...interface{}) {
	GetLogger().Error(str, args...)
}

// Fatal to log fatal.
func Fatal(str string, args ...interface{}) {
	GetLogger().Fatal(str, args...)
}
