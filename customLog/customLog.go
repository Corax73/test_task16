// customLog wrapper for using the standard log library.
package customLog

import (
	"log"
	"os"
)

// LogInit using the passed path string, creates, if missing, a log file and assigns errors to be output to it.
func LogInit(path string) {
	defaultPath := "./logs/app.log"
	if path == "" {
		path = defaultPath
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)
}

// Logging writes an indefinite number of transmitted non-empty errors to the log.
func Logging(errors ...error) {
	for _, err := range errors {
		if err != nil {
			log.Println(err)
		}
	}
}
