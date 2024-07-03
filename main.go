package main

import (
	"timeTracker/customDb"
	"timeTracker/router"

	_ "github.com/lib/pq"
)

func main() {
	result := customDb.Init()
	if result {
		router.RunRouter()
	}
}
