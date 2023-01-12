package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

// structures trading pair groups
func structureTradingPairs(pools []Pool) {
	// triangluarPairsList := []Pool{}
	// removeDuplicatesList := []Pool{}
	// pairsList := []Pool{}

	// loop through each coin to find potential matches
	for _, poolA := range pools {
		aPair := []string{poolA.Token0.Symbol, poolA.Token1.Symbol}
		// fmt.Println(aBox)

		// get second pair (B)
		for _, poolB := range pools {
			if poolB != poolA {
				if poolB.Token0.Symbol == aPair[0] || poolB.Token0.Symbol == aPair[1] || poolB.Token1.Symbol == aPair[0] || poolB.Token1.Symbol == aPair[1] {
					bPair := []string{poolB.Token0.Symbol, poolB.Token1.Symbol}

					// get third pair (C)
					for _, poolC := range pools {
						if poolC != poolA && poolC != poolB {

							// count number of C items
							// there should be two of everything
							pairBox := []string{poolA.Token0.Symbol, poolA.Token1.Symbol, poolB.Token0.Symbol, poolB.Token1.Symbol, poolC.Token0.Symbol, poolC.Token1.Symbol}

							cBaseCount := 0
							for _, t := range pairBox {
								if t == poolC.Token0.Symbol {
									cBaseCount++
								}
							}

							cQuoteCount := 0
							for _, t := range pairBox {
								if t == poolC.Token1.Symbol {
									cQuoteCount++
								}
							}

							if cBaseCount == 2 && cQuoteCount == 2 && poolC.Token0.Symbol != poolC.Token1.Symbol {
								cPair := []string{poolC.Token0.Symbol, poolC.Token1.Symbol}
								sort.Strings(aPair)
								sort.Strings(bPair)
								sort.Strings(cPair)
								combined := aPair + bPair + cPair
								fmt.Println(combined)
								// fmt.Println(pairBox)
							}
						}
					}
				}
			}
		}
	}
}

// queries uniswap graph for pools
func fetchPools() []Pool {
	q := `
	{
		pools(first: 20, orderBy: totalValueLockedETH, orderDirection: desc) {
			id
			token0Price
			token1Price
			token0 {
				id
				name
				symbol
				decimals
			}
			token1 {
				id
				name
				symbol
				decimals
			}
			totalValueLockedETH
		}
	}
	`

	b := query(q)

	var pools PoolRes
	json.Unmarshal(b, &pools)

	return pools.Data.Pools
}

// makes a query to uniswap graph and converts response body to byte slice
func query(q string) []byte {
	jsonData := map[string]string{
		"query": q,
	}

	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	defer response.Body.Close()
	data, _ := io.ReadAll(response.Body)

	return data
}
