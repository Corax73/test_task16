package utils

import (
	"fmt"
	"runtime"
	"timeTracker/customLog"

	"github.com/joho/godotenv"
)

// GetConfFromEnvFile receives data for the database from the environment file. If successful, returns a non-empty map.
func GetConfFromEnvFile() map[string]string {
	resp := make(map[string]string)
	envFile, err := godotenv.Read(".env")
	if err == nil {
		resp = envFile
	} else {
		customLog.Logging(err)
	}
	return resp
}

// GCRunAndPrintMemory runs a garbage collection and if setting the APP_ENV environment variable as "dev" prints currently allocated number of bytes on the heap.
func GCRunAndPrintMemory() {
	debugSet := false
	settings := GetConfFromEnvFile()
	if val, ok := settings["APP_ENV"]; ok && val == "dev" {
		debugSet = true
	}
	if debugSet {
		var stat runtime.MemStats
		runtime.ReadMemStats(&stat)
		fmt.Println(stat.Alloc / 1024)
	}
	if val, ok := settings["GC_MANUAL_RUN"]; ok && val == "true" {
		runtime.GC()
	}
}
