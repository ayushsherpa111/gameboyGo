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

type logger struct {
	lgr      *log.Logger
	debug    bool
	genCode  string
	errCode  string
	warnCode string
	escCode  string
}

func (l *logger) Printf(format string, val ...interface{}) {
	if !l.debug {
		return
	}
	l.lgr.Printf(format, val...)
}

func (l *logger) Warnf(format string, val ...interface{}) {
	l.Printf(l.warnCode+format+l.escCode, val...)
}

func (l *logger) Errorf(format string, val ...interface{}) {
	l.Printf(l.errCode+format+l.escCode, val...)
}

func (l *logger) Infof(format string, val ...interface{}) {
	l.Printf(l.genCode+format+l.escCode, val...)
}

func (l *logger) Fatalf(format string, val ...interface{}) {
	l.Errorf(format, val...)
	os.Exit(1)
}

func NewLogger(writer io.Writer, debug bool) logger {
	return logger{
		debug:    debug,
		lgr:      log.New(writer, "[GB] ", log.Lshortfile|log.Lmicroseconds),
		errCode:  fmt.Sprintf(FORMAT, ERR),
		genCode:  fmt.Sprintf(FORMAT, GEN),
		warnCode: fmt.Sprintf(FORMAT, WARN),
		escCode:  fmt.Sprintf(ESC, 0),
	}
}
