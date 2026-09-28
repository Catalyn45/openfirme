package common

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {	
	infoLogger *log.Logger
	warningLogger *log.Logger
	errorLogger *log.Logger
	criticalLogger *log.Logger
}

func (this *Logger) Info(v ...any) {
	this.infoLogger.Output(2, fmt.Sprintln(v...))
}

func (this *Logger) Warning(v ...any) {
	this.warningLogger.Output(2, fmt.Sprintln(v...))
}

func (this *Logger) Error(v ...any) {
	this.errorLogger.Output(2, fmt.Sprintln(v...))
}

func (this *Logger) Critical(v ...any) {
	this.criticalLogger.Output(2, fmt.Sprintln(v...))
}

var logger *Logger
func init() {
	logger = &Logger {
		infoLogger: log.New(
			os.Stdout,
			"[INFO] ",
			log.Ldate | log.Ltime | log.Lshortfile  | log.Lmsgprefix,
		),
		warningLogger: log.New(
			os.Stdout,
			"[WARNING] ",
			log.Ldate | log.Ltime | log.Lshortfile  | log.Lmsgprefix,
		),
		errorLogger: log.New(
			os.Stderr,
			"[ERROR] ",
			log.Ldate | log.Ltime | log.Lshortfile  | log.Lmsgprefix,
		),
		criticalLogger: log.New(
			os.Stderr,
			"[CRITICAL] ",
			log.Ldate | log.Ltime | log.Lshortfile  | log.Lmsgprefix,
		),
	}
}
