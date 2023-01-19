import { ethers } from 'ethers'
import Quoter from '@uniswap/v3-periphery/artifacts/contracts/lens/Quoter.sol/Quoter.json'
import * as rawSurfaceRates from './surface_rates.json'
import { SurfaceRate } from './types.js'

getDepth(1, 1)

async function getDepth(amountIn: number, limit: number) {
  console.log(amountIn, limit)
  const surfaceRates: SurfaceRate[] = rawSurfaceRates

  for (const sr of surfaceRates) {
    // trade 1
    console.log('Checking trade 1 acquired coin...')
    // const acquiredCoinDetail1 = getPrice(sr.Pool1.Id)

    // trade 2
    console.log('Checking trade 2 acquired coin...')

    // trade 3
    console.log('Checking trade 3 acquired coin...')
  }

  return
}
