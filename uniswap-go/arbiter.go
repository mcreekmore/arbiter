package main

// run project
// go run *.go

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {
	fmt.Println("STARTING ARBITER...")

	fmt.Println("Connecting to mongo...")
	connectDB()

	rawPools := fetchPools(100)

	pools := parsePools(&rawPools)

	stp := structureTradingPairs(&pools) // structured trading pairs

	srl := []SurfaceRate{} // surface rates list

	// calculate surface rates for triangular pair
	for _, tp := range stp {
		calcTokens(&tp)

		// calculate surface rate for specific token
		for _, t := range tp.Tokens {
			calcSurfaceRateForToken(t, tp, "foreward", .1, &srl)
			calcSurfaceRateForToken(t, tp, "reverse", .1, &srl)
		}
	}

	for _, sr := range srl {
		fmt.Println()
		fmt.Println("Starting:", sr.StartingAmount, sr.Token1)
		fmt.Println("Trade 1: ", sr.Token1+" -> "+sr.Token2, sr.AcquiredCoinTrade1)
		fmt.Println("Trade 2: ", sr.Token2+" -> "+sr.Token3, sr.AcquiredCoinTrade2)
		fmt.Println("Trade 3: ", sr.Token3+" -> "+sr.Token1, sr.AcquiredCoinTrade3)
		fmt.Printf("Profit %v%%", sr.ProfitLossPercent)
		fmt.Println()
	}

	file, _ := json.MarshalIndent(srl, "", "  ")
	_ = ioutil.WriteFile("../uniswap-js/src/surface_rates.json", file, 0644)
}
