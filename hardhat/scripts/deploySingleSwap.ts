import { ethers } from 'hardhat'

async function main() {
  const SingleSwap = await ethers.getContractFactory('SingleSwap')
  const singleSwap = await SingleSwap.deploy()

  await singleSwap.deployed()

  console.log('Single swap contract deployed: ', singleSwap.address)
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
