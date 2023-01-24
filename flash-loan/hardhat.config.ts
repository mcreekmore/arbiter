import { HardhatUserConfig } from 'hardhat/config'
import '@nomicfoundation/hardhat-toolbox'
import * as dotenv from 'dotenv'
dotenv.config()

require('dotenv').config({ path: '../.env' })

const config: HardhatUserConfig = {
  solidity: '0.8.17',
  networks: {
    goerli: {
      url:
        process.env.INFURA_GOERLI_ENDPOINT! +
        process.env.INFURA_FLASH_LOAN_KEY!,
      accounts: [process.env.GOERLI_PRIVATE_KEY!],
    },
  },
}

export default config
