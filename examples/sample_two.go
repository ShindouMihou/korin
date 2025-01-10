package main

import (
	"encoding/json"
	"fmt"
	"github.com/ShindouMihou/korin/examples/test"
)

func main() {
	statistics := test.Statistics{RunningMemory: 24_000} // +k:println
	stats, _ := json.Marshal(statistics)                 // +k:float
	stats = stats                                        // +k:println

	var name string
	fmt.Scanln(&name)
	fmt.Println("Hello", name)
}
