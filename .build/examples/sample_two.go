package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	type Statistics struct {
		RunningMemory int `json:"running_memory"`
	}
	statistics := Statistics{RunningMemory: 24_000}
	fmt.Println("statistics: ", statistics)
	stats, err := json.Marshal(statistics)
	if err != nil {
		return
	}
	stats = stats
	fmt.Println("stats: ", stats)

	var name string
	fmt.Scanln(&name)
	fmt.Println("Hello", name)
}
