import { ethers } from 'ethers'
import * as dotenv from 'dotenv'

dotenv.config()

require('dotenv').config({ path: '../.env' })

export const provider = new ethers.providers.JsonRpcProvider(
  process.env.INFURA_ETHEREUM_MAINNET
)
