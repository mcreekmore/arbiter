import requests
import func_arbitrage


# Step 0: Get list of tradeable pairs
def step_0():
    # get asset pairs
    res = requests.get(
        'https://api.kraken.com/0/public/AssetPairs').json()

    asset_pairs = res['result']

    # get tradeable pairs
    tradeable_pairs = func_arbitrage.get_tradeable_pairs(asset_pairs)

    return tradeable_pairs


# Step 1: Structuring Triangular Pairs
def step_1(pairs):
    structured_list = func_arbitrage.structure_arbitrage_pairs(pairs)


if __name__ == "__main__":
    tradeable_pairs = step_0()
    step_1(tradeable_pairs)
