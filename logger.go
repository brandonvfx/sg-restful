package main

import (
	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
}
