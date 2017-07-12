package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	// open a file
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
}

// SetupLogFile takes a log name, a level, and a pointer to a formatter, and configures the logger appropriately. SetupLogFile returns a pointer to an os.File object, which should be closed when done with.
func SetupLogFile(logname string, level log.Level, formatterPtr log.Formatter) *os.File {
	f, err := os.OpenFile(logname, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}

	log.SetLevel(level)
	log.SetOutput(f)
	log.SetFormatter(formatterPtr)
	return f
}
