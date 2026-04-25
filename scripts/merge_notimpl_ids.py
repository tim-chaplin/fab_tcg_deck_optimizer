"""Merge generated NotImplemented stub IDs / registrations into id.go and index.go.

Run after gen_notimpl_stubs.py:

  python scripts/merge_notimpl_ids.py

Reads /tmp/notimpl_id_block.txt and /tmp/notimpl_registry_block.txt and rewrites:
  internal/card/id.go         — interleaves new generic IDs alphabetically with existing ones,
                                adds TalisharID into the weapon block.
  internal/cards/index.go     — interleaves new generic registrations alphabetically.
  internal/weapon/registry.go — adds Talishar{} to weapon.All in alphabetical order.
"""
import os
import re

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
TMP = os.environ.get("TEMP", "/tmp")


def write_text(path, text):
    with open(path, "w", encoding="utf-8", newline="\n") as f:
        f.write(text)


def parse_blocks(path):
    """Snippet files use blank lines to delimit per-card families. Return [[id, ...], ...]."""
    text = open(path, encoding="utf-8").read()
    blocks = []
    cur = []
    for ln in text.splitlines():
        if not ln.strip():
            if cur:
                blocks.append(cur)
                cur = []
            continue
        cur.append(ln)
    if cur:
        blocks.append(cur)
    return blocks


def family_key(line):
    """Strip leading-tab + trailing colour suffix and any trailing struct text. Returns
    the family root used for alphabetisation."""
    m = re.search(r"\b([A-Za-z0-9]+?)(Red|Yellow|Blue)\b", line)
    if m:
        return m.group(1).lower()
    # Fallback: strip non-letters from the line
    return re.sub(r"[^A-Za-z]", "", line).lower()


def merge_id_go():
    path = os.path.join(ROOT, "internal", "card", "id.go")
    src = open(path, encoding="utf-8").read()

    # --- Generic block: between "// Generic card IDs." and the blank line preceding "// Weapon"
    gen_re = re.compile(
        r"(\t// Generic card IDs\.[^\n]*\n)"  # 1: header
        r"((?:\t[A-Za-z0-9]+\n)+)"             # 2: existing IDs (one per line)
        r"(\n\t// Weapon IDs\.)",              # 3: weapon header (kept verbatim)
        re.DOTALL,
    )
    m = gen_re.search(src)
    assert m, "couldn't find generic block in id.go"
    header, body, tail = m.group(1), m.group(2), m.group(3)

    existing_ids = [ln.strip() for ln in body.splitlines() if ln.strip()]

    # Group existing IDs by family root.
    families = {}
    for fid in existing_ids:
        k = family_key(fid)
        families.setdefault(k, []).append(fid)

    # Append new families.
    for fam in parse_blocks(os.path.join(TMP, "notimpl_id_block.txt")):
        ids = [ln.strip() for ln in fam if ln.strip()]
        if not ids:
            continue
        k = family_key(ids[0])
        # Skip if any of these already exist.
        existing_set = set(existing_ids)
        ids = [i for i in ids if i not in existing_set]
        if not ids:
            continue
        families.setdefault(k, []).extend(ids)

    # Sort families alphabetically by key, preserve R/Y/B order within.
    color_order = {"Red": 0, "Yellow": 1, "Blue": 2}

    def color_rank(idn):
        for c, r in color_order.items():
            if idn.endswith(c):
                return r
        return 9

    new_body_lines = []
    for k in sorted(families):
        ids = sorted(families[k], key=color_rank)
        new_body_lines.extend(f"\t{i}" for i in ids)
    new_body = "\n".join(new_body_lines) + "\n"

    new_src = src[: m.start()] + header + new_body + tail + src[m.end():]

    # --- Weapon block: insert TalisharID alphabetically.
    weapon_re = re.compile(
        r"(\t// Weapon IDs\.[^\n]*\n(?:\t//[^\n]*\n)?)"
        r"((?:\t[A-Za-z0-9]+\n)+)"
        r"(\n\t// Test-only)",
        re.DOTALL,
    )
    wm = weapon_re.search(new_src)
    assert wm, "couldn't find weapon block in id.go"
    wbody = wm.group(2)
    wids = [ln.strip() for ln in wbody.splitlines() if ln.strip()]
    if "TalisharID" not in wids:
        wids.append("TalisharID")
    wids.sort()
    new_wbody = "\n".join(f"\t{i}" for i in wids) + "\n"
    new_src = new_src[: wm.start(2)] + new_wbody + new_src[wm.end(2):]

    write_text(path, new_src)
    print("rewrote", path)


def merge_index_go():
    path = os.path.join(ROOT, "internal", "cards", "index.go")
    src = open(path, encoding="utf-8").read()

    # Generic registry block sits between the runeblade entries and the fake/test entries.
    gen_start = src.index("\tcard.AdrenalineRushRed: generic.AdrenalineRushRed{},")
    fake_idx = src.index("\tcard.FakeRedAttack:")
    gen_block = src[gen_start:fake_idx]

    # Parse: list of (family_root_lower, [lines])
    families = {}
    for fam in re.split(r"\n\n+", gen_block.strip("\n")):
        lines = [ln.rstrip() for ln in fam.split("\n") if ln.strip()]
        if not lines:
            continue
        m = re.match(r"\tcard\.([A-Za-z0-9]+?)(Red|Yellow|Blue): ", lines[0])
        if not m:
            # Single-printing or unusual: take everything up to the colour.
            m = re.match(r"\tcard\.([A-Za-z0-9]+?)(Red|Yellow|Blue|Blue):", lines[0])
        assert m, lines[0]
        families[m.group(1).lower()] = (m.group(1), lines)

    for fam in parse_blocks(os.path.join(TMP, "notimpl_registry_block.txt")):
        if not fam:
            continue
        m = re.match(r"\tcard\.([A-Za-z0-9]+?)(Red|Yellow|Blue): ", fam[0])
        assert m, fam[0]
        key = m.group(1).lower()
        if key in families:
            existing_lines = set(families[key][1])
            for ln in fam:
                if ln not in existing_lines:
                    families[key][1].append(ln)
        else:
            families[key] = (m.group(1), list(fam))

    color_order = {"Red": 0, "Yellow": 1, "Blue": 2}

    def color_rank(line):
        for c, r in color_order.items():
            if f"{c}:" in line.split("generic.")[0]:
                return r
        return 9

    blocks_out = []
    for k in sorted(families):
        _, lines = families[k]
        # Sort lines within family by colour.
        lines = sorted(lines, key=color_rank)
        blocks_out.append("\n".join(lines))

    new_block = "\n\n".join(blocks_out) + "\n\n"
    new_src = src[:gen_start] + new_block + src[fake_idx:]
    write_text(path, new_src)
    print("rewrote", path)


def merge_weapon_registry():
    path = os.path.join(ROOT, "internal", "weapon", "registry.go")
    src = open(path, encoding="utf-8").read()
    if "Talishar{}" in src:
        return
    m = re.search(r"(var All = \[\]Weapon\{\n)(.*?)(\n\})", src, re.DOTALL)
    assert m
    body = m.group(2)
    entries = [ln.rstrip() for ln in body.split("\n") if ln.strip()]
    entries.append("\tTalishar{},")
    entries.sort(key=lambda s: s.strip())
    new_body = "\n".join(entries)
    new_src = src[: m.start(2)] + new_body + src[m.end(2):]
    write_text(path, new_src)
    print("rewrote", path)


def main():
    merge_id_go()
    merge_index_go()
    merge_weapon_registry()


if __name__ == "__main__":
    main()
