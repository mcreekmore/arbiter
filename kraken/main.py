import requests

resp = requests.get(
    'https://api.kraken.com/0/public/AssetPairs')

asset_pairs = resp.json()

for pair in asset_pairs['result']:
    if 'XMR' in pair:
        print(asset_pairs['result'][pair]['base'] +
              '/' + asset_pairs['result'][pair]['quote'])
        print(pair)
        print(asset_pairs['result'][pair])
