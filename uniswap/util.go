package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func fetchPools() []Pool {
	q := `
	{
		pools(first: 5, orderBy: totalValueLockedETH) {
			id
			token0Price
			token1Price
			token1 {
				id
				name
				symbol
			}
			token0 {
				id
				name
				symbol
			}
			totalValueLockedETH
		}
	}
	`

	b := request(q)

	var pools PoolRes
	json.Unmarshal(b, &pools)

	return pools.Data.Pools
}

func request(q string) []byte {
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
	// fmt.Println(string(data))
	return data

}
