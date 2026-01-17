package logger

import (
	"log"
	"os"
)

var Debug bool

func InitLogger(debug bool) {
	Debug = debug
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func Info(msg string, args ...any) {
	log.Printf("[INFO] "+msg, args...)
}

func Debugf(msg string, args ...any) {
	if Debug {
		log.Printf("[DEBUG] "+msg, args...)
	}
}

func Error(msg string, args ...any) {
	log.Printf("[ERROR] "+msg, args...)
}
