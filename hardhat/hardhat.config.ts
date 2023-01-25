import { HardhatUserConfig } from 'hardhat/config'
import '@nomicfoundation/hardhat-toolbox'
import * as dotenv from 'dotenv'
import { version } from 'hardhat'
dotenv.config()

require('dotenv').config({ path: '../.env' })

// deployed flash loan (goerli): 0xB18bB0220e2Af838cb4f7f1AB6dEA70882c4203a
// deployed single swap (goerli): 0x5F92d80E8e99C3818bD912a1b49802Ce324D00Cd

const config: HardhatUserConfig = {
  solidity: {
    compilers: [{ version: '0.8.10' }, { version: '0.7.6' }],
  },
  networks: {
    goerli: {
      url:
        process.env.INFURA_GOERLI_ENDPOINT! +
        process.env.INFURA_FLASH_LOAN_KEY!,
      accounts: [process.env.METAMASK_PRIVATE_KEY!],
    },
    polygon: {
      url:
        process.env.INFURA_POLYGON_MAINNET! +
        process.env.INFURA_FLASH_LOAN_KEY!,
      accounts: [process.env.METAMASK_PRIVATE_KEY!],
    },
  },
}

export default config
