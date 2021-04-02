package logger

import (
	"fmt"
	"time"
)

type Logger interface {
	Log(args ...interface{})
}

type logger struct {
	tag string
}

func NewLogger(tag string) Logger {
	return &logger{tag}
}

func (l *logger) Log(args ...interface{}) {
	fmt.Print(time.Now().Format("2006-01-02 15:04:05"), "\t", l.tag, "\t")
	fmt.Println(args...)
}
