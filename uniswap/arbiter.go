package main

// run project
// go run *.go

import (
	"fmt"
)

func main() {
	fmt.Println("STARTING ARBITER...")

	pools := fetchPools()

	structureTradingPairs(pools)
}
