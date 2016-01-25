package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
)

/*

Not the test for main.go. This is the entry point for testing. Package level setup, including fixtures, go here.

Reminder - to test this package, use
go test -tags test
*/

// TestMain provides a mechanism for setting
func TestMain(m *testing.M) {
	// env variable defined in constants.go
	// LogTestsToFileEnvVar="SG-RESTFUL_LOG_TO_FILE"
	logTestsFlag := logTestsEnv()
	var f *os.File
	if logTestsFlag {
		nowtime := time.Now()
		timeStr := nowtime.Format("20060102150405")
		name := fmt.Sprintf("sg-restful-test.%s.log", timeStr)
		f = SetupLogFile(name,
			log.DebugLevel,
			&log.TextFormatter{})
	}
	log.Info("TestMain setup")

	res := m.Run()
	if f != nil {
		f.Close()
	}
	os.Exit(res)
}

func logTestsEnv() bool {
	if lt := os.Getenv(LogTestsToFileEnvVar); strings.ToLower(lt) == "true" ||
		strings.ToLower(lt) == "yes" {
		return true
	}
	return false
}
