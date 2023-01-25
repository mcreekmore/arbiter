package main

type Pool struct {
	Id                  string
	Token0Price         float64
	Token1Price         float64
	Token0              Token
	Token1              Token
	TotalValueLockedETH float64
	FeeTier             int64
}

type RawPool struct {
	Id                  string
	Token0Price         string
	Token1Price         string
	Token0              RawToken
	Token1              RawToken
	TotalValueLockedETH string
	FeeTier             string
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
	PoolA  Pool
	PoolB  Pool
	PoolC  Pool
	Tokens []string
}

type PoolRes struct {
	Data RawPools
}

type RawPools struct {
	Pools []RawPool
}

type SurfaceRate struct {
	Token1             string
	Token2             string
	Token3             string
	Pool1              Pool
	Pool2              Pool
	Pool3              Pool
	ProfitLoss         float64
	ProfitLossPercent  float64
	StartingAmount     float64
	AcquiredCoinTrade1 float64
	AcquiredCoinTrade2 float64
	AcquiredCoinTrade3 float64
	DirectionTrade1    string
	DirectionTrade2    string
	DirectionTrade3    string
}
