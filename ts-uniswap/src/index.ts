import { gql } from 'graphql-request'
import { fetchUniswapInfo } from './graphql.js'

async function main() {
  const data = fetchUniswapInfo()

  console.log(data)
}

main()
