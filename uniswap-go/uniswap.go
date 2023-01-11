package main

// run project
// go run *.go

import (
	"fmt"
)

func main() {
	fmt.Println("testing")

	pools := fetchPools()

	for _, p := range pools {
		fmt.Println(p)
	}
}
