"""Rewrite the generic block of internal/registry/ids/card_ids.go and
internal/registry/card_registry.go to include newly generated card IDs, merged
alphabetically with the existing hand-written ones.
"""
import os
import re

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
TMP = os.environ.get("TEMP", "/tmp")


def family(n):
    for suf in ("Blue", "Yellow", "Red"):
        if n.endswith(suf):
            return n[: -len(suf)]
    return n


def color(n):
    for suf in ("Red", "Yellow", "Blue"):
        if n.endswith(suf):
            return suf
    return ""


COLOR_ORDER = {"Red": 0, "Yellow": 1, "Blue": 2, "": 3}


def read_existing_generic_ids(content):
    # The generic-IDs block runs from "// Generic card IDs." down to the test-only fakes
    # block (weapon IDs now live in registry/ids/weapon_ids.go, not in this file).
    gen_start = content.index("\t// Generic card IDs.")
    gen_end = content.index("\t// Test-only synthetic")
    gen_block = content[gen_start:gen_end]
    ids = []
    for line in gen_block.splitlines():
        m = re.match(r"^\t([A-Z][A-Za-z0-9]+)\s*$", line)
        if m:
            ids.append(m.group(1))
    return gen_start, gen_end, ids


def build_ordered_ids(existing, new_ids):
    all_ids = existing + new_ids
    fam_map = {}
    for i in all_ids:
        fam_map.setdefault(family(i), []).append(i)
    ordered = []
    for fam in sorted(fam_map):
        within = sorted(fam_map[fam], key=lambda n: COLOR_ORDER[color(n)])
        ordered.append(within)
    return ordered


def render_id_block(ordered):
    lines = ["\t// Generic card IDs. Ordered alphabetically by card name, Red → Yellow → Blue within each family."]
    for fam_ids in ordered:
        for i in fam_ids:
            lines.append("\t" + i)
    return "\n".join(lines) + "\n"


def update_id_go(new_ids):
    path = os.path.join(ROOT, "internal", "registry", "ids", "card_ids.go")
    with open(path, encoding="utf-8") as f:
        content = f.read()
    gen_start, gen_end, existing = read_existing_generic_ids(content)
    ordered = build_ordered_ids(existing, new_ids)
    new_block = render_id_block(ordered)
    new_content = content[:gen_start] + new_block + "\n" + content[gen_end:]
    with open(path, "w", encoding="utf-8", newline="\n") as f:
        f.write(new_content)
    print("updated", path, "with", sum(len(x) for x in ordered), "generic IDs")


def read_existing_registry_entries(content):
    # The generic registry block runs from the first `ids.AdrenalineRushRed:` line
    # to the closing `}` of the cardsByID map.
    idx = content.index("\tids.AdrenalineRushRed:")
    end = content.index("\n}\n", idx)
    existing_block = content[idx:end]
    # Each entry is like: ids.DodgeBlue: cards.DodgeBlue{},
    pat = re.compile(r"^\tids\.([A-Za-z0-9]+):\s+cards\.[A-Za-z0-9]+\{\},\s*$", re.M)
    ids = pat.findall(existing_block)
    return idx, end, ids


def render_registry_block(ordered_families, max_name_len):
    lines = []
    for fam_ids in ordered_families:
        for i in fam_ids:
            lines.append(f"\tids.{i}: cards.{i}{{}},")
        lines.append("")
    # trim trailing blank
    while lines and lines[-1] == "":
        lines.pop()
    return "\n".join(lines) + "\n"


def update_index_go(new_ids):
    path = os.path.join(ROOT, "internal", "registry", "card_registry.go")
    with open(path, encoding="utf-8") as f:
        content = f.read()
    idx, end, existing = read_existing_registry_entries(content)
    ordered = build_ordered_ids(existing, new_ids)
    new_block = render_registry_block(ordered, 0)
    new_content = content[:idx] + new_block + content[end:]
    with open(path, "w", encoding="utf-8", newline="\n") as f:
        f.write(new_content)
    print("updated", path, "with", sum(len(x) for x in ordered), "generic entries")


def main():
    # Load newly generated IDs from the snippet file.
    snippet_path = os.path.join(TMP, "silverage_id_snippet.txt")
    with open(snippet_path, encoding="utf-8") as f:
        snippet = f.read()
    new_ids = [line.strip() for line in snippet.splitlines() if line.strip()]
    update_id_go(new_ids)
    update_index_go(new_ids)


if __name__ == "__main__":
    main()
