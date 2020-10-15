package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	limit := 300
	afterTenSeconds := start.Add(time.Duration(-limit) * time.Second)

	fmt.Printf("start = %v\n", start)
	fmt.Printf("afterTenSeconds = %v\n", afterTenSeconds)
	fmt.Printf("after = %v\n", start.After(afterTenSeconds))
	fmt.Printf("before = %v\n", start.Before(afterTenSeconds))
}
