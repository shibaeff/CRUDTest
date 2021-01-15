package functions

import (
	"log"
	"os"
	"time"
)

type TimeLogger interface {
	Start()
	End()
}

type timelogger struct {
	id int
	// logFile *os.File
	logger *log.Logger
}

func (t *timelogger) Start() {
	t.logger.Printf("start %d %v\n", t.id, time.Now())
}

func (t *timelogger) End() {
	t.logger.Printf("end %d %v\n", t.id, time.Now())
}

func NewTimeLogger(path, prefix string) (tl TimeLogger) {
	t := &timelogger{}
	t.id = 0
	// t.logFile, _ = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	t.logger = log.New(os.Stdout, prefix, log.LstdFlags)
	// defer t.logFile.Close()
	return t
}
