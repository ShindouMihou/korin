package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	type Statistics struct {
		RunningMemory int // +k:named(camelCase,json)
	}
	statistics := Statistics{RunningMemory: 24_000} // +k:println
	stats, _ := json.Marshal(statistics)            // +k:float
	stats = stats                                   // +k:println

	var name string
	fmt.Scanln(&name)
	fmt.Println("Hello", name)
}
