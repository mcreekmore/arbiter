package main

// run project
// go run *.go

import (
	"fmt"
)

func main() {
	fmt.Println("STARTING ARBITER...")

	rawPools := fetchPools()

	pools := parsePools(&rawPools)

	triPools := structureTradingPairs(&pools)

	// calculate surface rates
	for _, tp := range triPools {
		sr := calcSurfaceRate(&tp)
		fmt.Println(sr)
	}
}
