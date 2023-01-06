import json
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
    # find arbitrage pairs
    structured_list = func_arbitrage.structure_arbitrage_pairs(pairs)

    # save list to file
    with open('structured_triangular_pairs.json', "w") as fp:
        json.dump(structured_list, fp)


# Step 2: Calculate surface arbitrage opportunities
def step_2():
    # get structured pairs
    with open('structured_triangular_pairs.json', 'r') as json_file:
        structured_pairs = json.load(json_file)

    # get latest surface prices
    for t_pair in structured_pairs:
        prices_dict = func_arbitrage.get_price_for_t_pair(t_pair)

    # requests.get('https://api.kraken.com/0/public/Depth').json()


if __name__ == "__main__":
    # tradeable_pairs = step_0()
    # step_1(tradeable_pairs)
    step_2()
