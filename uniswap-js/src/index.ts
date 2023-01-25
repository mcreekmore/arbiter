import { ethers } from 'ethers'
import * as dotenv from 'dotenv'
import Quoter from '@uniswap/v3-periphery/artifacts/contracts/lens/Quoter.sol/Quoter.json'
import IUniswapV3PoolABI from '@uniswap/v3-core/artifacts/contracts/interfaces/IUniswapV3Pool.sol/IUniswapV3Pool.json'
import rawSurfaceRates from './surface_rates.json'
import { Pool, SurfaceRate, Token } from './types.js'
import { provider } from './ethers'

dotenv.config()
require('dotenv').config({ path: '../.env' })

getDepth(1, 1)

// calculate arbitrage
function calculateArbitrage(
  amountIn: number,
  amountOut: number,
  sr: SurfaceRate
) {
  let threshold = 0 // eventually add to config
  let profitLoss = amountOut - amountIn
  let profitLossPercent = 0

  if (profitLoss > threshold) {
    profitLossPercent = (profitLoss / amountIn) * 100
    // console.log(sr)
    return profitLossPercent
  }

  return profitLossPercent
}

async function getPrice(pool: Pool, amountIn: number, tradeDirection: string) {
  // const poolContract = new ethers.Contract(
  //   pool.Id,
  //   IUniswapV3PoolABI.abi,
  //   provider
  // )

  // const tokenFee = await poolContract.fee()
  // console.log(tokenFee)
  // const decimals

  let inputTokenA: Token
  let inputTokenB: Token

  // determine which token should be input A and input B respectively
  if (tradeDirection === 'baseToQuote') {
    inputTokenA = pool.Token0
    inputTokenB = pool.Token1
  } else {
    // (tradeDirection === 'quoteToBase')
    inputTokenA = pool.Token1
    inputTokenB = pool.Token0
  }

  // reformat amount into big number
  let amountInString = ethers.utils
    .parseUnits(
      amountIn.toFixed(inputTokenA.Decimals),
      inputTokenA.Decimals.toString()
    )
    .toString()

  const quoterContract = new ethers.Contract(
    process.env.UNISWAP_V3_QUOTER_ADDRESS!,
    Quoter.abi,
    provider
  )

  let outputAmount: string

  try {
    let quotedAmountOut = await quoterContract.callStatic.quoteExactInputSingle(
      inputTokenA.Id,
      inputTokenB.Id,
      pool.FeeTier,
      amountInString,
      0
    )

    outputAmount = ethers.utils.formatUnits(
      quotedAmountOut,
      ethers.BigNumber.from(inputTokenB.Decimals.toString())
    )

    return Number(outputAmount)
  } catch (err) {
    console.log(err)
  }

  return 0
}

async function getDepth(amountIn: number, limit: number) {
  const surfaceRates: SurfaceRate[] = rawSurfaceRates

  for (const sr of surfaceRates) {
    // trade 1
    console.log('Checking new opportunity...')
    const acquiredCoinDetail1 = await getPrice(
      sr.Pool1,
      amountIn,
      sr.DirectionTrade1
    )

    if (acquiredCoinDetail1 === 0) {
      console.log('error: acquired coin 1 = 0')
      return
    }

    // trade 2
    const acquiredCoinDetail2 = await getPrice(
      sr.Pool2,
      acquiredCoinDetail1,
      sr.DirectionTrade2
    )

    if (acquiredCoinDetail2 === 0) {
      console.log('error: acquired coin 2 = 0')
      return
    }

    // trade 3
    const acquiredCoinDetail3 = await getPrice(
      sr.Pool3,
      acquiredCoinDetail2,
      sr.DirectionTrade3
    )

    let profitLossPercent = calculateArbitrage(
      amountIn,
      acquiredCoinDetail3,
      sr
    )
    if (profitLossPercent > 0) {
      console.log(`Profit Loss: ${profitLossPercent}%`)
      console.log(sr)

      // EXECUUUUUUUUUUUUUUUUTTEEEEEEEEEEEEEE
      let res = await executeArbiter(sr)
    }
  }

  return
}

async function executeArbiter(sr: SurfaceRate) {
  // TODO

  return 0
}
