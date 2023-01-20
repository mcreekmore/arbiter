export type SurfaceRate = {
  Token1: string
  Token2: string
  Token3: string
  Pool1: Pool
  Pool2: Pool
  Pool3: Pool
  ProfitLoss: number
  ProfitLossPercent: number
  StartingAmount: number
  AcquiredCoinTrade1: number
  AcquiredCoinTrade2: number
  AcquiredCoinTrade3: number
  DirectionTrade1: string
  DirectionTrade2: string
  DirectionTrade3: string
}

export type Pool = {
  Id: string
  Token0Price: number
  Token1Price: number
  Token0: Token
  Token1: Token
  TotalValueLockedETH: number
  FeeTier: number
}

export type Token = {
  Id: string
  Name: string
  Symbol: string
  Decimals: number
}

export type QuotedDepth = {
  SurfaceRate: SurfaceRate
  ProfitLoss: number
  ProfitLossPercent: number
}
