import csv
import sys

from pymongo import MongoClient

if len(sys.argv) < 2:
    print("missing args: python calculate_market_cap.py [db_ip_address]")
    exit(0)

client = MongoClient(f"mongodb://{sys.argv[1]}/marketDB")

fundamentals_coll = client.get_database().get_collection("fundamentals")

raw_fundamentals = fundamentals_coll.find()


def format_date_number(n):
    s = str(n)
    if len(s) == 1:
        return "0" + s
    return s


def increment_date(d):
    d[2] += 1
    if d[2] > 31:
        d[2] = 1
        d[1] += 1
    if d[1] > 12:
        d[1] = 1
        d[0] += 1


with open("./data/market_cap.csv", "w+") as output:
    output.write("Symbol,Quarter,MarketCap\n")
    for funda in raw_fundamentals:
        current_market_cap = funda["Highlights"]["MarketCapitalizationMln"]
        if current_market_cap is None or current_market_cap < 2000.0:
            continue
        quarterly_stats = funda["outstandingShares"]["quarterly"]
        if len(quarterly_stats) < 80:
            continue
        symbol = funda["General"]["Code"]
        print(f"parsing {symbol}")
        day_map = dict()
        with open(f"candles/daily_{symbol}.csv") as f:
            daily_candles = csv.reader(f)
            is_header = True
            for row in daily_candles:
                if is_header:
                    is_header = False
                    continue
                day_map[row[0]] = row

        # Parse all quarters, excluding the first one
        temp_output = []
        first_entry = True
        for key in quarterly_stats:
            if first_entry:
                first_entry = False
                continue
            stats = quarterly_stats[key]
            # retrieve price from the day of earnings or the next available day
            date_split = [int(a) for a in stats["dateFormatted"].split("-")]
            if date_split[0] >= 2022:
                continue
            if date_split[0] < 1986:
                continue
            eod_candle = None
            attempts = 0
            while eod_candle is None:
                try:
                    eod_candle = day_map["-".join([format_date_number(a) for a in date_split])]
                    break
                except KeyError:
                    pass
                attempts += 1
                if attempts > 30:
                    break
                increment_date(date_split)
            if attempts > 30:
                continue

            closing_price = eod_candle[5]
            shares = stats["shares"]
            market_cap = int(float(closing_price) * shares)
            quarter = stats["date"]
            temp_output.append(f"{symbol},{quarter},{market_cap}\n")

        if len(temp_output) >= 80:
            output.writelines(temp_output)
