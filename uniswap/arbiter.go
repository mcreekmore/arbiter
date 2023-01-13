package main

// run project
// go run *.go

import (
	"fmt"
)

func main() {
	fmt.Println("STARTING ARBITER...")

	pools := fetchPools()

	triPools := structureTradingPairs(pools)

	for _, tp := range triPools {
		fmt.Println(tp.PoolA.Token0.Symbol, tp.PoolA.Token1.Symbol)
		fmt.Println(tp.PoolB.Token0.Symbol, tp.PoolB.Token1.Symbol)
		fmt.Println(tp.PoolC.Token0.Symbol, tp.PoolC.Token1.Symbol)
	}
}
