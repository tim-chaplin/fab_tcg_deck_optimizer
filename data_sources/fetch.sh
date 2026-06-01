#!/usr/bin/env bash
# Fetch the latest English card.csv from the-fab-cube/flesh-and-blood-cards.
# Run from anywhere; writes to the directory containing this script.
set -euo pipefail

dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# The omens-of-the-third-age branch carries the OMN (Omens of the Third Age) set, the latest
# release, which hasn't been merged to develop yet. It's a superset of the older
# compendium-of-rathe branch (it adds OMN and drops nothing).
url="https://raw.githubusercontent.com/the-fab-cube/flesh-and-blood-cards/omens-of-the-third-age/csvs/english/card.csv"

echo "Fetching $url"
curl -sSLf "$url" -o "$dir/card.csv"
echo "Wrote $dir/card.csv ($(wc -l < "$dir/card.csv") lines)"
