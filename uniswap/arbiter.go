package main

// run project
// go run *.go

import (
	"fmt"
)

func main() {
	fmt.Println("STARTING ARBITER...")

	// fmt.Println("Connecting to mongo...")
	// connectDB()

	rawPools := fetchPools(30)

	pools := parsePools(&rawPools)

	triPools := structureTradingPairs(&pools)

	// calculate surface rates for triangular pair
	for _, tp := range triPools {
		calcTokens(&tp)

		// calculate surface rate for specific token
		for _, t := range tp.Tokens {
			calcSurfaceRateForToken(t, tp, "foreward", .2)
			calcSurfaceRateForToken(t, tp, "reverse", .2)
		}

		// sr := calcSurfaceRate(&tp)
		// fmt.Println(sr)
	}
}
