package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	ESC    = "\x1b[%dm"
	FORMAT = "\x1b[38;5;%dm"
	GEN    = 33
	WARN   = 215
	ERR    = 196
)

type Logger struct {
	lgr      *log.Logger
	debug    bool
	genCode  string
	errCode  string
	warnCode string
	escCode  string
	prefix   string
}

func (l *Logger) Printf(format string, val ...interface{}) {
	if !l.debug {
		return
	}
	l.lgr.Printf(format, val...)
}

func (l *Logger) Warnf(format string, val ...interface{}) {
	l.Printf(l.warnCode+format+l.escCode, val...)
}

func (l *Logger) Errorf(format string, val ...interface{}) {
	l.Printf(l.errCode+format+l.escCode, val...)
}

func (l *Logger) Infof(format string, val ...interface{}) {
	l.Printf(l.genCode+format+l.escCode, val...)
}

func (l *Logger) Fatalf(format string, val ...interface{}) {
	l.Errorf(format, val...)
	os.Exit(1)
}

func NewLogger(writer io.Writer, debug bool, prefix string) Logger {
	return Logger{
		debug:    debug,
		lgr:      log.New(writer, fmt.Sprintf("[%s] ", prefix), log.Lshortfile|log.Lmicroseconds),
		errCode:  fmt.Sprintf(FORMAT, ERR),
		genCode:  fmt.Sprintf(FORMAT, GEN),
		warnCode: fmt.Sprintf(FORMAT, WARN),
		escCode:  fmt.Sprintf(ESC, 0),
	}
}
