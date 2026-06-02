#!/usr/bin/env bash
# Fetch the English card data CSVs from the-fab-cube/flesh-and-blood-cards.
# Run from anywhere; writes to the directory containing this script.
#
#   card.csv          - one row per card (the main database the tools read)
#   card-printing.csv - one row per physical printing; carries the per-printing Rarity, joined
#                       to card.csv by the card's Unique ID. Needed for the Silver Age rarity
#                       rule (a card is legal only if it has a Basic/Common/Rare printing).
#   rarity.csv        - rarity shorthand (C/R/M/…) -> full name.
set -euo pipefail

dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# The omens-of-the-third-age branch carries the OMN (Omens of the Third Age) set, the latest
# release, which hasn't been merged to develop yet. It's a superset of the older
# compendium-of-rathe branch (it adds OMN and drops nothing).
base="https://raw.githubusercontent.com/the-fab-cube/flesh-and-blood-cards/omens-of-the-third-age/csvs/english"

for f in card.csv card-printing.csv rarity.csv; do
	echo "Fetching $base/$f"
	curl -sSLf "$base/$f" -o "$dir/$f"
	echo "Wrote $dir/$f ($(wc -l < "$dir/$f") lines)"
done
