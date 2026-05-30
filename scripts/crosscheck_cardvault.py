"""Cross-check our local card pool against the official cardvault.fabtcg.com DB.

Pulls every Silver-Age-legal Runeblade and Generic *print* from the official
advanced-search API, dedupes to card names, applies the same type exclusions
parsecarddb uses, and diffs against our local card.csv.
"""
import csv
import json
import sys
import time
import urllib.parse
import urllib.request

API = "https://api.cardvault.fabtcg.com/carddb/api/v1/advanced-search/"
# Subtype/type words parsecarddb excludes from the modelled pool.
EXCLUDE = ("Shadow", "Elemental", "Lightning", "Earth", "Token", "Equipment")


def fetch_class(cls):
    """Return list of (printed_name, printed_typebox) for one class, Silver Age legal."""
    out = []
    page = 1
    while True:
        params = {
            "classes": cls,
            "legal_formats": "Silver Age",
            "page_size": "250",
            "page": str(page),
            "orderby": "name",
        }
        url = API + "?" + urllib.parse.urlencode(params)
        req = urllib.request.Request(url, headers={"Accept": "application/json"})
        with urllib.request.urlopen(req, timeout=60) as r:
            d = json.load(r)
        for c in d.get("results", []):
            out.append((c.get("printed_name", ""), c.get("printed_typebox", "")))
        if not d.get("next"):
            break
        page += 1
        time.sleep(0.2)
    return out


def in_scope(typebox):
    return not any(w in typebox for w in EXCLUDE)


def official_names():
    """Map printed_name -> set of typeboxes, for Runeblade+Generic Silver Age prints."""
    by_name = {}
    raw_prints = 0
    for cls in ("Runeblade", "Generic"):
        rows = fetch_class(cls)
        raw_prints += len(rows)
        for name, typebox in rows:
            by_name.setdefault(name, set()).add(typebox)
    return by_name, raw_prints


def load_local(path):
    """All Runeblade/Generic names in our local CSV (any legality) + the SA in-scope subset."""
    rune_or_generic = set()
    sa_inscope = set()
    with open(path, encoding="utf-8") as f:
        r = csv.DictReader(f, delimiter="\t")
        for row in r:
            types = row.get("Types", "")
            name = row.get("Name", "")
            if "Runeblade" not in types and "Generic" not in types:
                continue
            rune_or_generic.add(name)
            sa = row.get("Silver Age Legal", "")
            if sa not in ("Yes", ""):
                continue
            if any(w in types for w in EXCLUDE):
                continue
            sa_inscope.add(name)
    return rune_or_generic, sa_inscope


def main():
    local_path = sys.argv[1] if len(sys.argv) > 1 else "data_sources/card.csv"
    by_name, raw_prints = official_names()

    # Official names that, after type exclusions, still have at least one in-scope printing.
    official_inscope = {n for n, tbs in by_name.items() if any(in_scope(tb) for tb in tbs)}
    official_all = set(by_name)

    local_rg, local_sa = load_local(local_path)

    print(f"official prints pulled (Runeblade+Generic, Silver Age): {raw_prints}")
    print(f"official unique names: {len(official_all)}")
    print(f"official unique names after type exclusions (in-scope): {len(official_inscope)}")
    print(f"local Runeblade/Generic names (any legality): {len(local_rg)}")
    print(f"local Silver-Age in-scope names: {len(local_sa)}")
    print()

    # The headline: in-scope official names missing entirely from our local CSV.
    missing_from_csv = sorted(official_inscope - local_rg)
    print(f"=== [A] In-scope official names ABSENT from local card.csv ({len(missing_from_csv)}) ===")
    for n in missing_from_csv:
        tbs = "; ".join(sorted(by_name[n]))
        print(f"  {n}  [{tbs}]")
    print()

    # Present in CSV but our SA-filter dropped them, yet official says in-scope+SA: filter drift.
    in_csv_not_inscope = sorted((official_inscope & local_rg) - local_sa)
    print(f"=== [B] Official in-scope+SA, in CSV, but NOT in our local in-scope set ({len(in_csv_not_inscope)}) ===")
    for n in in_csv_not_inscope:
        tbs = "; ".join(sorted(by_name[n]))
        print(f"  {n}  [{tbs}]")
    print()

    # Reverse sanity: our in-scope names the official SA query didn't return.
    local_only = sorted(local_sa - official_inscope)
    print(f"=== [C] Local in-scope names NOT in official SA Runeblade/Generic set ({len(local_only)}) ===")
    for n in local_only:
        print(f"  {n}")


if __name__ == "__main__":
    main()
