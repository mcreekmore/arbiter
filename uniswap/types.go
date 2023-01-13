package main

type Pool struct {
	Id                  string
	Token0Price         float64
	Token1Price         float64
	Token0              Token
	Token1              Token
	TotalValueLockedETH float64
}

type RawPool struct {
	Id                  string
	Token0Price         string
	Token1Price         string
	Token0              RawToken
	Token1              RawToken
	TotalValueLockedETH string
}

type Token struct {
	Id       string
	Name     string
	Symbol   string
	Decimals int64
}

type RawToken struct {
	Id       string
	Name     string
	Symbol   string
	Decimals string
}

type TriPool struct {
	PoolA Pool
	PoolB Pool
	PoolC Pool
}

type PoolRes struct {
	Data RawPools
}

type RawPools struct {
	Pools []RawPool
}
