package utils

import (
	"fmt"
	"runtime"
)

func PrintMemoryAndGC() {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	fmt.Println(stat.Alloc / 1024)
	runtime.GC()
}
