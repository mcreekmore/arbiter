package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func calcSurfaceRateForToken(token1 string, tp TriPool, d string, msr float64) {
	sa := 1.0
	var pool1 Pool
	var pool2 Pool
	var pool3 Pool
	var dt1 string // direction trade 1
	var dt2 string // direction trade 2
	var dt3 string // direction trade 3
	var token2 string
	var token3 string
	// var token3 string
	var swapRate1 float64
	var swapRate2 float64
	var swapRate3 float64
	var act1 float64 // acquired coin transaction 1
	var act2 float64 // acquired coin transaction 2
	var act3 float64 // acquired coin transaction 3

	// forward
	// find first pool that has token 1
	if d == "foreward" {
		if tp.PoolA.Token0.Symbol == token1 || tp.PoolA.Token1.Symbol == token1 {
			pool1 = tp.PoolA
			if tp.PoolA.Token0.Symbol == token1 {
				dt1 = "baseToQuote"
				token2 = tp.PoolA.Token1.Symbol
				swapRate1 = tp.PoolA.Token1Price
			} else {
				dt1 = "quoteToBase"
				token2 = tp.PoolA.Token0.Symbol
				swapRate1 = tp.PoolA.Token0Price
			}
		} else if tp.PoolB.Token0.Symbol == token1 || tp.PoolB.Token1.Symbol == token1 {
			pool1 = tp.PoolB
			if tp.PoolB.Token0.Symbol == token1 {
				dt1 = "baseToQuote"
				token2 = tp.PoolB.Token1.Symbol
				swapRate1 = tp.PoolB.Token1Price
			} else {
				dt1 = "quoteToBase"
				token2 = tp.PoolB.Token0.Symbol
				swapRate1 = tp.PoolB.Token0Price
			}
		} else if tp.PoolC.Token0.Symbol == token1 || tp.PoolC.Token1.Symbol == token1 {
			pool1 = tp.PoolC
			if tp.PoolC.Token0.Symbol == token1 {
				dt1 = "baseToQuote"
				token2 = tp.PoolC.Token1.Symbol
				swapRate1 = tp.PoolC.Token1Price
			} else {
				dt1 = "quoteToBase"
				token2 = tp.PoolC.Token0.Symbol
				swapRate1 = tp.PoolC.Token0Price
			}
		}
	}

	// reverse
	// find first pool that has token 1
	if d == "reverse" {
		if tp.PoolC.Token0.Symbol == token1 || tp.PoolC.Token1.Symbol == token1 {
			pool1 = tp.PoolC
			if tp.PoolC.Token0.Symbol == token1 {
				dt1 = "baseToQuote"
				token2 = tp.PoolC.Token1.Symbol
				swapRate1 = tp.PoolC.Token1Price
			} else {
				dt1 = "quoteToBase"
				token2 = tp.PoolC.Token0.Symbol
				swapRate1 = tp.PoolC.Token0Price
			}
		} else if tp.PoolB.Token0.Symbol == token1 || tp.PoolB.Token1.Symbol == token1 {
			pool1 = tp.PoolB
			if tp.PoolB.Token0.Symbol == token1 {
				dt1 = "baseToQuote"
				token2 = tp.PoolB.Token1.Symbol
				swapRate1 = tp.PoolB.Token1Price
			} else {
				dt1 = "quoteToBase"
				token2 = tp.PoolB.Token0.Symbol
				swapRate1 = tp.PoolB.Token0Price
			}
		} else if tp.PoolA.Token0.Symbol == token1 || tp.PoolA.Token1.Symbol == token1 {
			pool1 = tp.PoolA
			if tp.PoolA.Token0.Symbol == token1 {
				dt1 = "baseToQuote"
				token2 = tp.PoolA.Token1.Symbol
				swapRate1 = tp.PoolA.Token1Price
			} else {
				dt1 = "quoteToBase"
				token2 = tp.PoolA.Token0.Symbol
				swapRate1 = tp.PoolA.Token0Price
			}
		}
	}

	// execute first trade
	act1 = float64(sa) * swapRate1

	// find next pool to use for token 2
	if tp.PoolA != pool1 && (tp.PoolA.Token0.Symbol == token2 || tp.PoolA.Token1.Symbol == token2) {
		pool2 = tp.PoolA
		if tp.PoolA.Token0.Symbol == token2 {
			dt2 = "baseToQuote"
			// token2 = pool2.Token1.Symbol
			token3 = pool2.Token1.Symbol
			swapRate2 = pool2.Token1Price
		} else {
			dt2 = "quoteToBase"
			// token2 = pool2.Token0.Symbol
			token3 = pool2.Token0.Symbol
			swapRate2 = pool2.Token0Price
		}
	} else if tp.PoolB != pool1 && (tp.PoolB.Token0.Symbol == token2 || tp.PoolB.Token1.Symbol == token2) {
		pool2 = tp.PoolB
		if tp.PoolB.Token0.Symbol == token2 {
			dt2 = "baseToQuote"
			// token2 = pool2.Token1.Symbol
			token3 = pool2.Token1.Symbol
			swapRate2 = pool2.Token1Price
		} else {
			dt2 = "quoteToBase"
			// token2 = pool2.Token0.Symbol
			token3 = pool2.Token0.Symbol
			swapRate2 = pool2.Token0Price
		}
	} else if tp.PoolC != pool1 && (tp.PoolC.Token0.Symbol == token2 || tp.PoolC.Token1.Symbol == token2) {
		pool2 = tp.PoolC
		if tp.PoolC.Token0.Symbol == token2 {
			dt2 = "baseToQuote"
			// token2 = pool2.Token1.Symbol
			token3 = pool2.Token1.Symbol
			swapRate2 = pool2.Token1Price
		} else {
			dt2 = "quoteToBase"
			// token2 = pool2.Token0.Symbol
			token3 = pool2.Token0.Symbol
			swapRate2 = pool2.Token0Price
		}
	}

	// execute second trade
	act2 = act1 * swapRate2

	// find remaining pool
	if tp.PoolA != pool1 && tp.PoolA != pool2 {
		pool3 = tp.PoolA
		if tp.PoolA.Token0.Symbol == token3 {
			dt3 = "baseToQuote"
			swapRate3 = tp.PoolA.Token1Price
		} else {
			dt3 = "quoteToBase"
			swapRate3 = tp.PoolA.Token0Price
		}
	} else if tp.PoolB != pool1 && tp.PoolB != pool2 {
		pool3 = tp.PoolB
		if tp.PoolB.Token0.Symbol == token3 {
			dt3 = "baseToQuote"
			swapRate3 = tp.PoolB.Token1Price
		} else {
			dt3 = "quoteToBase"
			swapRate3 = tp.PoolB.Token0Price
		}
	} else if tp.PoolC != pool1 && tp.PoolC != pool2 {
		pool3 = tp.PoolC
		if tp.PoolC.Token0.Symbol == token3 {
			dt3 = "baseToQuote"
			swapRate3 = tp.PoolC.Token1Price
		} else {
			dt3 = "quoteToBase"
			swapRate3 = tp.PoolC.Token0Price
		}
	}

	// execute third trade
	act3 = act2 * swapRate3

	// remove for production
	Use(dt1, dt2, dt3, pool3, act2, act3)

	// calculate profit/loss
	pl := act3 - sa        // profit/loss
	plp := (pl / sa) * 100 // profit/loss percent

	if plp > msr { // msr = minimum surface rate
		fmt.Println()
		fmt.Println("Starting:", sa, token1)
		fmt.Println("Trade 1: ", token1+" -> "+token2, act1)
		fmt.Println("Trade 2: ", token2+" -> "+token3, act2)
		fmt.Println("Trade 3: ", token3+" -> "+token1, act3)
		fmt.Println("Profit", plp, "%")
	}

}

func calcTokens(tp *TriPool) {
	if !contains(&tp.Tokens, &tp.PoolA.Token0.Symbol) {
		tp.Tokens = append(tp.Tokens, tp.PoolA.Token0.Symbol)
	}
	if !contains(&tp.Tokens, &tp.PoolA.Token1.Symbol) {
		tp.Tokens = append(tp.Tokens, tp.PoolA.Token1.Symbol)
	}
	if !contains(&tp.Tokens, &tp.PoolB.Token0.Symbol) {
		tp.Tokens = append(tp.Tokens, tp.PoolB.Token0.Symbol)
	}
	if !contains(&tp.Tokens, &tp.PoolB.Token1.Symbol) {
		tp.Tokens = append(tp.Tokens, tp.PoolB.Token1.Symbol)
	}
	if !contains(&tp.Tokens, &tp.PoolC.Token0.Symbol) {
		tp.Tokens = append(tp.Tokens, tp.PoolC.Token0.Symbol)
	}
	if !contains(&tp.Tokens, &tp.PoolC.Token1.Symbol) {
		tp.Tokens = append(tp.Tokens, tp.PoolC.Token1.Symbol)
	}
}

// structures trading pair groups
func structureTradingPairs(pools *[]Pool) []TriPool {
	tPoolsList := []TriPool{}
	removeDuplicatePoolsList := []string{}

	// loop through each coin to find potential matches
	for _, poolA := range *pools {
		aPair := []string{poolA.Token0.Symbol, poolA.Token1.Symbol}

		// get second pair (B)
		for _, poolB := range *pools {
			if poolB != poolA {
				if poolB.Token0.Symbol == aPair[0] || poolB.Token0.Symbol == aPair[1] ||
					poolB.Token1.Symbol == aPair[0] || poolB.Token1.Symbol == aPair[1] {

					// get third pair (C)
					for _, poolC := range *pools {
						if poolC != poolA && poolC != poolB {

							// count number of C items
							// there should be two of everything
							tokenList := []string{
								poolA.Token0.Symbol,
								poolA.Token1.Symbol,
								poolB.Token0.Symbol,
								poolB.Token1.Symbol,
								poolC.Token0.Symbol,
								poolC.Token1.Symbol,
							}

							cBaseCount := 0
							for _, t := range tokenList {
								if t == poolC.Token0.Symbol {
									cBaseCount++
								}
							}

							cQuoteCount := 0
							for _, t := range tokenList {
								if t == poolC.Token1.Symbol {
									cQuoteCount++
								}
							}

							if cBaseCount == 2 && cQuoteCount == 2 && poolC.Token0.Symbol != poolC.Token1.Symbol {
								sort.Strings(tokenList)
								uniqueString := strings.Join(tokenList, " ")

								// output pair
								if !contains(&removeDuplicatePoolsList, &uniqueString) {
									triPool := TriPool{
										PoolA:  poolA,
										PoolB:  poolB,
										PoolC:  poolC,
										Tokens: []string{},
									}

									tPoolsList = append(tPoolsList, triPool)
									removeDuplicatePoolsList = append(removeDuplicatePoolsList, uniqueString)
								}
							}
						}
					}
				}
			}
		}
	}

	return tPoolsList
}

// graphql returns all values as strings
// parse to floats/ints when appropriate and error check
func parsePools(rp *[]RawPool) []Pool {
	pp := []Pool{}

	for _, p := range *rp {
		t0p, err := strconv.ParseFloat(p.Token0Price, 64)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		t1p, err := strconv.ParseFloat(p.Token1Price, 64)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		tvle, err := strconv.ParseFloat(p.TotalValueLockedETH, 64)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		t0d, err := strconv.ParseInt(p.Token0.Decimals, 10, 64)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		t1d, err := strconv.ParseInt(p.Token1.Decimals, 10, 64)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		np := Pool{
			Id:          p.Id,
			Token0Price: t0p,
			Token1Price: t1p,
			Token0: Token{
				Id:       p.Token0.Id,
				Name:     p.Token0.Name,
				Symbol:   p.Token0.Symbol,
				Decimals: t0d,
			},
			Token1: Token{
				Id:       p.Token1.Id,
				Name:     p.Token1.Name,
				Symbol:   p.Token1.Symbol,
				Decimals: t1d,
			},
			TotalValueLockedETH: tvle,
		}

		pp = append(pp, np)
	}

	return pp
}

// queries uniswap graph for pools
func fetchPools(n int) []RawPool {
	q := fmt.Sprintf(`
	{
		pools(first: %d, orderBy: totalValueLockedETH, orderDirection: desc) {
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
	`, n)

	b := query(&q)

	var pools PoolRes
	json.Unmarshal(b, &pools)

	return pools.Data.Pools
}

// makes a query to uniswap graph and converts response body to byte slice
func query(q *string) []byte {
	jsonData := map[string]string{
		"query": *q,
	}

	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)

	defer response.Body.Close()

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	data, _ := io.ReadAll(response.Body)

	return data
}

// checks if a string is found in array
func contains(s *[]string, str *string) bool {
	for _, v := range *s {
		if v == *str {
			return true
		}
	}
	return false
}

// gets rid of pesky "declared but not used"
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

// calculates surface rates for a triangular pair of pools
// func calcSurfaceRate(tp *TriPool) float32 {
// sa := 1.0 // starting amount
// // msr := 1.5 				// minimum surface rate
// var s1 string    // swap 1
// var s2 string    // swap 2
// var s3 string    // swap 3
// var s1r float64  // swap 1 rate
// var s2r float64  // swap 2 rate
// var s3r float64  // swap 3 rate
// var dt1 string   // direction trade 1
// var dt2 string   // direction trade 2
// var dt3 string   // direction trade 3
// var c1 string    // contract 1
// var c2 string    // contract 2
// var c3 string    // contract 3
// var act1 float64 // acquired coin trade 1
// var act2 float64 // acquired coin trade 2
// var act3 float64 // acquired coin trade 3
// c := false       // calculated

// dl := []string{"forward", "reverse"}

// for _, d := range dl {
// 	// assume start with aBase if forward
// 	if d == "forward" {
// 		s1 = tp.PoolA.Token0.Symbol
// 		s2 = tp.PoolA.Token1.Symbol
// 		s1r = tp.PoolA.Token1Price
// 		dt1 = "baseToQuote"
// 	}

// 	// assume start with aQuote if forward
// 	if d == "reverse" {
// 		s1 = tp.PoolA.Token1.Symbol
// 		s2 = tp.PoolA.Token0.Symbol
// 		s1r = tp.PoolA.Token0Price
// 		dt1 = "quoteToBase"
// 	}

// 	// place first trade
// 	c1 = tp.PoolA.Token0.Symbol + "_" + tp.PoolA.Token1.Symbol
// 	act1 = sa * s1r

// 	// forward: check if aQuote (acquired coin) matches bQuote
// 	if d == "forward" {
// 		if tp.PoolA.Token0.Id == tp.PoolB.Token0.Id && c == false {
// 			s2r = tp.PoolB.Token0Price
// 			act2 = act1 * s2r
// 			dt2 = "quoteToBase"
// 			c2 = tp.PoolB.Token0.Symbol + "_" + tp.PoolB.Token1.Symbol

// 			// forward: check if bBase (acquired coin) matches cBase
// 			if tp.PoolB.Token0.Symbol == tp.PoolC.Token0.Symbol {
// 				s3 = tp.PoolC.Token0.Symbol
// 				s3r = tp.PoolC.Token1Price
// 				dt3 = "baseToQuote"
// 				c3 = tp.PoolC.Token0.Symbol + "_" + tp.PoolC.Token1.Symbol
// 			}

// 			// forward: check if bBase (acquired coin) matches cQuote
// 			if tp.PoolB.Token0.Symbol == tp.PoolC.Token1.Symbol {
// 				s3 = tp.PoolC.Token1.Symbol
// 				s3r = tp.PoolC.Token0Price
// 				dt3 = "quoteToBase"
// 				c3 = tp.PoolC.Token0.Symbol + "_" + tp.PoolC.Token1.Symbol
// 			}
// 		}
// 	}
// }

// fmt.Println(tp.PoolA.Token0.Symbol + "_" + tp.PoolA.Token1.Symbol + " / " + tp.PoolB.Token0.Symbol + "_" + tp.PoolB.Token1.Symbol + " / " + tp.PoolC.Token0.Symbol + "_" + tp.PoolC.Token1.Symbol)
// 	return 0
// }
