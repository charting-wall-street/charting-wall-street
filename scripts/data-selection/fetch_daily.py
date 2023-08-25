import csv
import sys

import requests

if len(sys.argv) < 2:
    print("missing args: python fetch_tickers.py [api_key]")
    exit(0)

api_token = sys.argv[1]

with open("../fundamentals-fetcher/data/tickers.txt", "r") as f:
    reader = csv.DictReader(f)
    for line in reader:
        symbol = line["Code"]
        print(f"fetching {symbol} candles")
        url = f"https://eodhistoricaldata.com/api/eod/{symbol}.US?api_token={api_token}&period=d"
        result = requests.get(url)
        with open(f"candles/daily_{symbol}.csv", "w+") as f:
            f.write(result.text)
