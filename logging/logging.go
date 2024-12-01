package logging

import (
	"log"
	"os"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

// Init allways is executed before the main function
func init() {
	//logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	logOutput := os.Stdout
	/*
		if err != nil {
			log.Fatalf("error creating/opening the log file: %v", err)
		}
	*/

	infoLogger = log.New(logOutput, "INFO: ", log.Ldate|log.Ltime)
	warningLogger = log.New(logOutput, "WARNING: ", log.Ldate|log.Ltime)
	errorLogger = log.New(logOutput, "ERROR: ", log.Ldate|log.Ltime)
}

func LogInfo(textFormat string, args ...interface{}) {
	infoLogger.Printf(textFormat, args...)
}

func LogWarning(textFormat string, args ...interface{}) {
	warningLogger.Printf(textFormat, args...)
}

func LogError(textFormat string, args ...interface{}) {
	errorLogger.Printf(textFormat, args...)
}

func LogFatal(textFormat string, args ...interface{}) {
	errorLogger.Fatalf(textFormat, args...)
}
