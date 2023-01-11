package main

type PoolRes struct {
	Data Pools
}

type Pools struct {
	Pools []Pool
}

type Pool struct {
	Id                  string
	Token0Price         string
	Token1Price         string
	Token0              Token
	Token1              Token
	TotalValueLockedETH string
}

type Token struct {
	Id     string
	Name   string
	Symbol string
}
