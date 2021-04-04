package logger

import (
	"fmt"
	"strings"
	"time"
)

type LOG_LEVEL int

const (
	LOG_LEVEL_DEBUG LOG_LEVEL = iota
	LOG_LEVEL_INFO
	LOG_LEVEL_ERROR
)

type Logger interface {
	AddExtraTag(tag string) Logger
	Log(lvl LOG_LEVEL, args ...interface{})
	LogDebug(args ...interface{})
	LogInfo(args ...interface{})
	LogError(args ...interface{})
}

type logger struct {
	tag       string
	lvl       LOG_LEVEL
	extraTags []string
}

func NewLogger(tag string, lvl LOG_LEVEL) Logger {
	return &logger{
		tag: tag,
		lvl: lvl,
	}
}

func (l *logger) AddExtraTag(tag string) Logger {
	l.extraTags = append(l.extraTags, tag)
	return l
}

func (l *logger) Log(lvl LOG_LEVEL, args ...interface{}) {
	if l.lvl <= lvl {
		fmt.Print(time.Now().Format("2006-01-02 15:04:05"), "\t", translateLogLevel(lvl), "\t", l.tag, "\t")
		if len(l.extraTags) > 0 {
			fmt.Print("[", strings.Join(l.extraTags, ", "), "] ")
		}
		fmt.Println(args...)
	}
}

func (l *logger) LogDebug(args ...interface{}) {
	l.Log(LOG_LEVEL_DEBUG, args...)
}

func (l *logger) LogInfo(args ...interface{}) {
	l.Log(LOG_LEVEL_INFO, args...)
}

func (l *logger) LogError(args ...interface{}) {
	l.Log(LOG_LEVEL_ERROR, args...)
}

func translateLogLevel(lvl LOG_LEVEL) (strlvl string) {
	strlvl = "unknown"
	switch lvl {
	case LOG_LEVEL_DEBUG:
		strlvl = "debug"
	case LOG_LEVEL_INFO:
		strlvl = "info"
	case LOG_LEVEL_ERROR:
		strlvl = "error"
	}
	return strlvl
}
