import csv

quarters = dict()

with open("./data/market_cap.csv", "r") as f:
    reader = csv.DictReader(f)
    for entry in reader:
        if entry["Quarter"] not in quarters:
            quarters[entry["Quarter"]] = []
        quarters[entry["Quarter"]].append(entry)

tracked_symbols = set()

for q in quarters:
    quarters[q].sort(key=lambda x: int(x["MarketCap"]), reverse=True)
    i = 0
    for entry in quarters[q]:
        if i == 250:
            break
        tracked_symbols.add(entry["Symbol"])
        i += 1

ordered_symbols = sorted(list(tracked_symbols))

with open("./data/candidate_stocks.txt", "w+") as output:
    output.writelines(ordered_symbols)

time_series = "Symbol"
rows = {}
for symbol in ordered_symbols:
    rows[symbol] = symbol

for q in reversed(quarters):
    time_series += "," + q
    existing_symbols = set()
    for entry in quarters[q]:
        if entry["Symbol"] not in rows:
            continue
        existing_symbols.add(entry["Symbol"])
        rows[entry["Symbol"]] += "," + str(entry['MarketCap'])
    missing_symbols = tracked_symbols - existing_symbols
    for symbol in missing_symbols:
        rows[symbol] += ","

with open("./data/market_cap_table.csv", "w+") as output:
    output.write(time_series + "\n")
    for symbol in ordered_symbols:
        output.write(rows[symbol] + "\n")

print(len(tracked_symbols))
