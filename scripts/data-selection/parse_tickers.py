import csv

# Most common exchanges
included_exchanges = ["NYSE", "AMEX", "NYSE MKT", "NASDAQ"]

with open("data/tickers_raw.csv", "r") as f:
    rows = list(csv.DictReader(f))

candidates = []
for row in rows:
    # Exclude funds
    if row["Code"][-1] == "X":
        continue
    # Exclude all other types of stocks than regular stocks
    if row["Type"] != "Common Stock":
        continue
    # Make sure exchange is a relevant US exchange
    if row["Exchange"] not in included_exchanges:
        continue
    # Make sure stock is publicly traded and not a derivative
    if row["Isin"] == "":
        continue
    # Exclude unit stocks
    if len(row["Code"]) > 4 and row["Code"][-1] == "U":
        continue
    # Exclude odd symbols containing dashes
    if "-" in row["Code"]:
        continue
    # Exclude symbols that are not the standard 4 length
    if len(row["Code"]) > 4:
        continue

    candidates.append(row)

print(f"found {len(candidates)} stocks")

with open("data/tickers.txt", "w+") as f:
    writer = csv.DictWriter(f, fieldnames=dict.keys(candidates[0]))
    writer.writeheader()
    writer.writerows(candidates)
