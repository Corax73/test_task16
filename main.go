package main

import (
	"timeTracker/customDb"
	"timeTracker/customLog"
	"timeTracker/router"

	_ "github.com/lib/pq"
)

func main() {
	customLog.LogInit("./logs/app.log")
	result := customDb.Init()
	if result {
		router.RunRouter()
	}
}
