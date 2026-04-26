"""Generate NotImplemented{} stubs for Silver Age cards that are still missing entirely.

Run: python scripts/gen_notimpl_stubs.py

Produces:
  internal/card/generic/<name>.go for each missing card family (one file per name with
  per-colour structs).
  internal/weapon/talishar.go for the Talishar weapon stub.
  /tmp/notimpl_id_block.txt — the new ID constants to insert into internal/card/id.go.
  /tmp/notimpl_registry_block.txt — the new entries to insert into internal/cards/index.go.
  /tmp/notimpl_weapon_id.txt + /tmp/notimpl_weapon_registry.txt — Talishar bits.

Every generated card has NotImplemented(){} so the optimizer skips it from random / mutation pools.
Stats (Cost / Pitch / Power / Defense / GoAgain) come straight from card.csv via parsecarddb's
JSON dump (/tmp/all_cards.json). Play returns 0; the card's printed bonus is never modelled.
"""

import json
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
TMP = os.environ.get("TEMP", "/tmp")

# Card names from TODO.md's "Unimplemented generic Silver-Age cards" section. These are the
# cards we currently have NO file for at all. We stub them with NotImplemented so the optimizer
# stops introducing them via mutation / random generation but they still resolve as static stats
# if a hand-authored deck includes one.
ITEMS = [
    "Amulet of Assertiveness", "Amulet of Echoes", "Amulet of Havencall",
    "Amulet of Ignition", "Amulet of Intervention", "Amulet of Oblation",
    "Clap 'Em in Irons", "Clarity Potion", "Energy Potion", "Healing Potion",
    "Imperial Seal of Command", "Potion of Déjà Vu", "Potion of Ironhide",
    "Potion of Luck", "Potion of Seeing", "Potion of Strength",
    "Talisman of Balance", "Talisman of Cremation", "Talisman of Dousing",
    "Talisman of Featherfoot", "Talisman of Recompense", "Talisman of Tithes",
    "Talisman of Warfare", "Timesnap Potion",
]

INSTANTS = [
    "Arcane Polarity", "Brush Off", "Calming Breeze", "Count Your Blessings",
    "Destructive Tendencies", "Eirina's Prayer", "Even Bigger Than That!",
    "Memorial Ground", "Oasis Respite", "Peace of Mind", "Pilfer the Tomb",
    "Reinforce the Line", "Shatter Sorcery", "Sigil of Solace",
]

ATTACK_REACTIONS = [
    "Blade Flash", "Exposed", "Lunging Press", "Nip at the Heels", "Pummel",
    "Razor Reflex", "Thrust",
]

RESOURCES = ["Cracked Bauble", "Fool's Gold", "Titanium Bauble"]

WEAPONS = ["Talishar, the Lost Prince"]

# Map names that need explicit identifier handling (apostrophes, accents, punctuation).
IDENT_OVERRIDES = {
    "Clap 'Em in Irons":          "ClapEmInIrons",
    "Eirina's Prayer":            "EirinasPrayer",
    "Even Bigger Than That!":     "EvenBiggerThanThat",
    "Fool's Gold":                "FoolsGold",
    "Nip at the Heels":           "NipAtTheHeels",
    "Potion of Déjà Vu": "PotionOfDejaVu",
    "Talishar, the Lost Prince":  "Talishar",
}

FILE_OVERRIDES = {
    "Clap 'Em in Irons":          "clap_em_in_irons",
    "Eirina's Prayer":            "eirinas_prayer",
    "Even Bigger Than That!":     "even_bigger_than_that",
    "Fool's Gold":                "fools_gold",
    "Potion of Déjà Vu": "potion_of_deja_vu",
    "Talishar, the Lost Prince":  "talishar",
}


def to_identifier(name):
    if name in IDENT_OVERRIDES:
        return IDENT_OVERRIDES[name]
    parts = re.split(r"[\s\-',!?]+", name)
    return "".join(p[0].upper() + p[1:] for p in parts if p)


def to_filename(name):
    if name in FILE_OVERRIDES:
        return FILE_OVERRIDES[name]
    s = name.lower()
    s = re.sub(r"[^a-z0-9]+", "_", s)
    return s.strip("_")


def load_cards():
    """Index cards.json by name → list of printings (one entry per colour)."""
    with open(os.path.join(TMP, "all_cards.json"), encoding="utf-8") as f:
        cards = json.load(f)
    by_name = {}
    for c in cards:
        by_name.setdefault(c["Name"], []).append(c)
    return by_name


def go_string(s):
    return '"' + s.replace("\\", "\\\\").replace('"', '\\"') + '"'


def types_expr_for(types_field):
    """Translate the CSV Types column (e.g. 'Generic, Action, Item') into a Go NewTypeSet expr.

    Resources have no dedicated TypeSet — we represent them with TypeGeneric only since the CSV
    keyword 'Resource' isn't a CardType in this codebase (and these cards aren't actionable in
    the sim — they're pitch fodder).
    """
    parts = [p.strip() for p in types_field.split(",")]
    mapping = {
        "Generic":          "card.TypeGeneric",
        "Action":           "card.TypeAction",
        "Attack":           "card.TypeAttack",
        "Attack Reaction":  "card.TypeAttackReaction",
        "Aura":             "card.TypeAura",
        "Defense Reaction": "card.TypeDefenseReaction",
        "Instant":          "card.TypeInstant",
        "Item":             "card.TypeItem",
        "Sword":            "card.TypeSword",
        "2H":               "card.TypeTwoHand",
        "1H":               "card.TypeOneHand",
        "Weapon":           "card.TypeWeapon",
    }
    consts = []
    for p in parts:
        if p == "Resource":
            # No CardType for Resource; the card is still tagged Generic so it routes through the
            # generic-pitch path. NotImplemented keeps it out of random pools regardless.
            continue
        if p in mapping:
            consts.append(mapping[p])
    if not consts:
        consts.append("card.TypeGeneric")
    # de-dupe preserving order
    seen, out = set(), []
    for c in consts:
        if c not in seen:
            seen.add(c)
            out.append(c)
    return "card.NewTypeSet(" + ", ".join(out) + ")"


def go_int(field, default=0):
    v = (field or "").strip()
    if v == "" or v == "X":
        return default
    try:
        return int(v)
    except ValueError:
        return default


def render_card_file(name, printings):
    """Generate a .go file for one card name with all of its colour printings."""
    ident_root = to_identifier(name)
    file_root = to_filename(name)
    types_var = ident_root[0].lower() + ident_root[1:] + "Types"
    types_expr = types_expr_for(printings[0]["Types"])

    order = {"Red": 0, "Yellow": 1, "Blue": 2}
    printings = sorted(printings, key=lambda c: order.get(c.get("Color", ""), 9))

    type_text = printings[0].get("TypeText", "") or printings[0]["Types"]

    # Variant summary for the header docstring.
    pitches = [(p.get("Color", ""), go_int(p.get("Pitch"))) for p in printings]
    powers = [(p.get("Color", ""), go_int(p.get("Power"))) for p in printings]
    defs = [(p.get("Color", ""), go_int(p.get("Defense"))) for p in printings]
    cost = go_int(printings[0].get("Cost"))

    bits = [f"{name} — {type_text}.", f"Cost {cost}."]
    if any(v != 0 for _, v in powers):
        bits.append("Printed power: " + ", ".join(f"{c} {v}" for c, v in powers if c) + ".")
    bits.append("Printed pitch variants: " + ", ".join(f"{c} {v}" for c, v in pitches if c) + ".")
    if any(v != 0 for _, v in defs):
        if len({v for _, v in defs}) == 1:
            bits.append(f"Defense {defs[0][1]}.")
        else:
            bits.append("Printed defense: " + ", ".join(f"{c} {v}" for c, v in defs if c) + ".")
    summary = " ".join(bits)

    text = re.sub(r"\s+", " ", (printings[0].get("FunctionalText") or "")).strip()

    paragraphs = [summary]
    if text:
        paragraphs.append("")
        paragraphs.append(f"Text: \"{text}\"")
    paragraphs.append("")
    paragraphs.append(
        "Stub only — marked NotImplemented so the optimizer skips it. The printed effect "
        "isn't modelled; Play returns 0."
    )

    doc = build_doc_comment(paragraphs)

    structs = []
    for p in printings:
        color = p.get("Color", "")
        s = f"{ident_root}{color}"
        pitch = go_int(p.get("Pitch"))
        pcost = go_int(p.get("Cost"))
        power = go_int(p.get("Power"))
        defense = go_int(p.get("Defense"))
        kw = p.get("CardKeywords", "") or ""
        go_again = "true" if "Go again" in kw else "false"

        structs.append("\n".join([
            f"type {s} struct{{}}",
            "",
            f"func ({s}) ID() card.ID                            {{ return card.{s} }}",
            f"func ({s}) Name() string                           {{ return {go_string(name + ' (' + color + ')')} }}",
            f"func ({s}) Cost(*card.TurnState) int               {{ return {pcost} }}",
            f"func ({s}) Pitch() int                             {{ return {pitch} }}",
            f"func ({s}) Attack() int                            {{ return {power} }}",
            f"func ({s}) Defense() int                           {{ return {defense} }}",
            f"func ({s}) Types() card.TypeSet                    {{ return {types_var} }}",
            f"func ({s}) GoAgain() bool                          {{ return {go_again} }}",
            f"func ({s}) NotImplemented()                        {{}}",
            f"func ({s}) Play(*card.TurnState, *card.CardState) int {{ return 0 }}",
        ]))

    parts = [doc.rstrip(), "", "package generic", "",
             'import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"', "",
             f"var {types_var} = {types_expr}", "",
             "\n\n".join(structs), ""]
    return file_root, "\n".join(parts), [(ident_root, p.get("Color", "")) for p in printings]


def render_weapon_file(name, printings):
    """Single-colour weapon stub (Talishar)."""
    p = printings[0]
    ident = to_identifier(name)
    file_root = to_filename(name)
    text = re.sub(r"\s+", " ", (p.get("FunctionalText") or "")).strip()
    type_text = p.get("TypeText", "") or p["Types"]
    power = go_int(p.get("Power"))

    paragraphs = [
        f"{name} — {type_text}. Power {power}.",
        "",
        f"Text: \"{text}\"",
        "",
        "Stub only — marked NotImplemented so the optimizer skips it. The activation cost "
        "and rust-counter destruction clause aren't modelled; Play returns the printed power "
        "and the weapon never destroys itself.",
    ]
    doc = build_doc_comment(paragraphs)

    types_expr = types_expr_for(p["Types"])
    types_var = ident[0].lower() + ident[1:] + "Types"
    id_const = ident + "ID"

    body = "\n".join([
        f"type {ident} struct{{}}",
        "",
        f"func ({ident}) ID() card.ID                            {{ return card.{id_const} }}",
        f"func ({ident}) Name() string                           {{ return {go_string(name)} }}",
        f"func ({ident}) Cost(*card.TurnState) int               {{ return 0 }}",
        f"func ({ident}) Pitch() int                             {{ return 0 }}",
        f"func ({ident}) Attack() int                            {{ return {power} }}",
        f"func ({ident}) Defense() int                           {{ return 0 }}",
        f"func ({ident}) Types() card.TypeSet                    {{ return {types_var} }}",
        f"func ({ident}) GoAgain() bool                          {{ return false }}",
        f"func ({ident}) Hands() int                             {{ return 2 }}",
        f"func ({ident}) NotImplemented()                        {{}}",
        f"func (c {ident}) Play(*card.TurnState, *card.CardState) int {{ return c.Attack() }}",
    ])

    parts = [doc.rstrip(), "", "package weapon", "",
             'import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"', "",
             f"var {types_var} = {types_expr}", "",
             body, ""]
    return file_root, "\n".join(parts), ident, id_const


def wrap_paragraph(text, prefix="// ", width=100):
    words = text.split()
    if not words:
        return prefix.rstrip()
    out, cur = [], prefix + words[0]
    for w in words[1:]:
        cand = cur + " " + w
        if len(cand) > width:
            out.append(cur)
            cur = prefix + w
        else:
            cur = cand
    out.append(cur)
    return "\n".join(out)


def build_doc_comment(paragraphs):
    out = []
    for p in paragraphs:
        if p == "":
            out.append("//")
        else:
            out.append(wrap_paragraph(p))
    return "\n".join(out) + "\n"


def main():
    by_name = load_cards()

    all_names = ITEMS + INSTANTS + ATTACK_REACTIONS + RESOURCES
    out_dir = os.path.join(ROOT, "internal", "card", "generic")
    weapon_dir = os.path.join(ROOT, "internal", "weapon")

    # Collect (name, ident_root, color) tuples in alphabetical order for ID and registry blocks.
    id_entries = []  # list of (ident_root_with_color, comment_name)
    reg_entries = []  # list of (ident_full, struct_full, name)

    missing = []
    for name in sorted(all_names):
        if name not in by_name:
            missing.append(name)
            continue
        printings = by_name[name]
        # Filter to legal-or-unknown rows (the CSV has a few duplicate entries from old printings).
        # SilverAgeLegal is blank for many; treat blank or "Yes" as legal.
        printings = [p for p in printings if p.get("SilverAgeLegal") in ("Yes", "")]
        if not printings:
            missing.append(name)
            continue
        # De-dup by Color (some cards appear twice with the same printing).
        seen_color = set()
        deduped = []
        for p in printings:
            c = p.get("Color", "")
            if c in seen_color:
                continue
            seen_color.add(c)
            deduped.append(p)
        printings = deduped

        file_root, content, structs = render_card_file(name, printings)
        path = os.path.join(out_dir, file_root + ".go")
        if os.path.exists(path):
            print(f"SKIP existing card file: {path}", file=sys.stderr)
            continue
        with open(path, "w", encoding="utf-8", newline="\n") as f:
            f.write(content)
        print("wrote", path)
        for ident_root, color in structs:
            id_entries.append((f"{ident_root}{color}", name))
            reg_entries.append((f"{ident_root}{color}", f"generic.{ident_root}{color}{{}}", name))

    # Talishar weapon
    weapon_id_lines = []
    weapon_reg_lines = []
    for name in WEAPONS:
        if name not in by_name:
            missing.append(name)
            continue
        printings = by_name[name]
        file_root, content, ident, id_const = render_weapon_file(name, printings)
        path = os.path.join(weapon_dir, file_root + ".go")
        if os.path.exists(path):
            print(f"SKIP existing weapon file: {path}", file=sys.stderr)
        else:
            with open(path, "w", encoding="utf-8", newline="\n") as f:
                f.write(content)
            print("wrote", path)
        weapon_id_lines.append(f"\t{id_const}")
        weapon_reg_lines.append(f"\t{ident}{{}},")

    # Snippet for id.go: list of new IDs to insert in alphabetical order in the generic block.
    id_path = os.path.join(TMP, "notimpl_id_block.txt")
    with open(id_path, "w", encoding="utf-8", newline="\n") as f:
        last = None
        for ident, name in id_entries:
            if last is not None and last != name:
                f.write("\n")
            f.write(f"\t{ident}\n")
            last = name
    print("ID snippet:", id_path)

    reg_path = os.path.join(TMP, "notimpl_registry_block.txt")
    with open(reg_path, "w", encoding="utf-8", newline="\n") as f:
        last = None
        for ident, struct_name, name in reg_entries:
            if last is not None and last != name:
                f.write("\n")
            f.write(f"\tcard.{ident}: {struct_name},\n")
            last = name
    print("Registry snippet:", reg_path)

    weapon_id_path = os.path.join(TMP, "notimpl_weapon_id.txt")
    with open(weapon_id_path, "w", encoding="utf-8", newline="\n") as f:
        f.write("\n".join(weapon_id_lines) + "\n")
    weapon_reg_path = os.path.join(TMP, "notimpl_weapon_registry.txt")
    with open(weapon_reg_path, "w", encoding="utf-8", newline="\n") as f:
        f.write("\n".join(weapon_reg_lines) + "\n")
    print("Weapon snippets:", weapon_id_path, weapon_reg_path)

    if missing:
        print("MISSING from CSV (or filtered out):", missing, file=sys.stderr)


if __name__ == "__main__":
    main()
