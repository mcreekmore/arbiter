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

async function getPrice(pool: Pool, amountIn: number, tradeDirection: string) {
  const poolContract = new ethers.Contract(
    pool.Id,
    IUniswapV3PoolABI.abi,
    provider
  )

  const tokenFee = await poolContract.fee()
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
      amountIn.toString(),
      ethers.BigNumber.from(inputTokenA.Decimals.toString())
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
      tokenFee,
      amountInString,
      0
    )

    outputAmount = ethers.utils.formatUnits(
      quotedAmountOut,
      ethers.BigNumber.from(inputTokenB.Decimals.toString())
    )

    return parseFloat(outputAmount)
  } catch (err) {
    console.log(err)
  }

  return 0
}

async function getDepth(amountIn: number, limit: number) {
  const surfaceRates: SurfaceRate[] = rawSurfaceRates

  for (const sr of surfaceRates) {
    // trade 1
    console.log('Checking trade 1 acquired coin...')
    const acquiredCoinDetail1 = await getPrice(
      sr.Pool1,
      amountIn,
      sr.DirectionTrade1
    )

    if (acquiredCoinDetail1 === 0) {
      return
    }

    // trade 2
    console.log('Checking trade 2 acquired coin...')
    const acquiredCoinDetail2 = await getPrice(
      sr.Pool2,
      acquiredCoinDetail1,
      sr.DirectionTrade2
    )

    if (acquiredCoinDetail2 === 0) {
      return
    }

    // trade 3
    console.log('Checking trade 3 acquired coin...')
    const acquiredCoinDetail3 = await getPrice(
      sr.Pool3,
      acquiredCoinDetail2,
      sr.DirectionTrade3
    )

    console.log(amountIn, acquiredCoinDetail3)
  }

  return
}
