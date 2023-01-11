import { GraphQLClient, gql } from 'graphql-request'

const API_KEY = ''

export const client = new GraphQLClient(
  `https://gateway.thegraph.com/api/${API_KEY}/subgraphs/id/ELUcwgpm14LKPLrBRuVvPvNKHQ9HvwmtKgKSH6123cr7`,
  { headers: {} }
)

export async function fetchUniswapInfo() {
  const POOLS_QUERY = gql`
    {
      liquidityPools(
        orderDirection: desc
        first: 100
        orderBy: totalValueLockedUSD
      ) {
        id
        name
        symbol
      }
    }
  `

  return await client.request(POOLS_QUERY)
}
