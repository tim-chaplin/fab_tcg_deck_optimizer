#!/usr/bin/env bash
# Fetch the latest English card.csv from the-fab-cube/flesh-and-blood-cards.
# Run from anywhere; writes to the directory containing this script.
set -euo pipefail

dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
url="https://raw.githubusercontent.com/the-fab-cube/flesh-and-blood-cards/develop/csvs/english/card.csv"

echo "Fetching $url"
curl -sSLf "$url" -o "$dir/card.csv"
echo "Wrote $dir/card.csv ($(wc -l < "$dir/card.csv") lines)"
