package main

import log "github.com/Sirupsen/logrus"

func init() {
	// open a file

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
}
