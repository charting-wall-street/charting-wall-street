import requests
import json
import csv
import time
import sys
from pymongo import MongoClient

if len(sys.argv) < 3:
    print("missing args: python fetch_fundamentals.py [db_ip_address] [api_key]")
    exit(0)

api_token = sys.argv[2]
exchange_symbol = "US"
url = "https://eodhistoricaldata.com/api/fundamentals"
db_url = f"mongodb://{sys.argv[1]}"

client = MongoClient(f"{db_url}/marketDB")
collection = client.get_database().get_collection("fundamentals")

with open("data/tickers.txt", "r") as f:
    rows = list(csv.DictReader(f))

for row in rows:
    attempts = 0
    ticker = row["Code"]
    print(f"fetching {ticker}")
    while True:
        if attempts > 10:
            print("too many fails")
            exit(1)
        try:
            result = requests.get(f"{url}/{ticker}.US?fmt=json&api_token={api_token}")
            data = json.loads(result.text)
            collection.insert_one(data)
            break
        except:
            print("fetch failed, retrying")
            time.sleep(attempts * 10)
            attempts += 1
