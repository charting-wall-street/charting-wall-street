import sys

import requests

if len(sys.argv) < 2:
    print("missing args: python fetch_tickers.py [api_key]")
    exit(0)

api_token = sys.argv[1]
url = f"https://eodhistoricaldata.com/api/exchange-symbol-list/US?api_token={api_token}"
result = requests.get(url)

with open("data/tickers_raw.csv", "w+") as f:
    f.write(result.text)
