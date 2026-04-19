"""Generate Go source files for Silver Age Generic cards.

Run: python scripts/gen_silverage_generics.py

Outputs:
  internal/card/generic/*.go — one per card family.
  Appends generated IDs to internal/card/id.go (generic block).
  Appends registry entries to internal/cards/index.go.

Idempotent: existing hand-written card files (dodge, evasive_leap, etc.) are skipped.
"""

import json
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
TMP = os.environ.get("TEMP", "/tmp")

IN_SCOPE_TYPES = {
    "Generic, Action",
    "Generic, Action, Attack",
    "Generic, Action, Aura",
    "Generic, Defense Reaction",
    "Generic, Block",
}
ALREADY_DONE = {
    "Dodge", "Evasive Leap", "Fate Foreseen", "Lay Low", "Put in Context",
    "Rise Above", "Sink Below", "Springboard Somersault", "Toughen Up", "Unmovable",
}

# Manual overrides: map (card name -> camel-case identifier root) for names that need tweaking
# (e.g. non-ASCII or punctuation).
IDENT_OVERRIDES = {
    "Captain's Call": "CaptainsCall",
    "Fiddler's Green": "FiddlersGreen",
    "Fyendal's Fighting Spirit": "FyendalsFightingSpirit",
    "Jack Be Nimble": "JackBeNimble",
    "Jack Be Quick": "JackBeQuick",
    "Money or Your Life?": "MoneyOrYourLife",
    "Money Where Ya Mouth Is": "MoneyWhereYaMouthIs",
    "Pick a Card, Any Card": "PickACardAnyCard",
    "Pursue to the Edge of Oblivion": "PursueToTheEdgeOfOblivion",
    "Pursue to the Pits of Despair": "PursueToThePitsOfDespair",
    "Tremor of \u00edArathael": "TremorOfIArathael",
    "Warmonger's Recital": "WarmongersRecital",
}

FILE_OVERRIDES = {
    "Captain's Call": "captains_call",
    "Fiddler's Green": "fiddlers_green",
    "Fyendal's Fighting Spirit": "fyendals_fighting_spirit",
    "Money or Your Life?": "money_or_your_life",
    "Money Where Ya Mouth Is": "money_where_ya_mouth_is",
    "Pick a Card, Any Card": "pick_a_card_any_card",
    "Tremor of \u00edArathael": "tremor_of_iarathael",
    "Warmonger's Recital": "warmongers_recital",
}

def to_identifier(name):
    if name in IDENT_OVERRIDES:
        return IDENT_OVERRIDES[name]
    # strip non-alphanumerics and title-case
    parts = re.split(r"[\s\-',!?]+", name)
    return "".join(p[0].upper() + p[1:] for p in parts if p)


def to_filename(name):
    if name in FILE_OVERRIDES:
        return FILE_OVERRIDES[name]
    s = name.lower()
    s = re.sub(r"[^a-z0-9]+", "_", s)
    s = s.strip("_")
    return s


def load_cards():
    with open(os.path.join(TMP, "cards.json"), encoding="utf-8") as f:
        cards = json.load(f)
    generics = [c for c in cards if "Generic" in c["Types"] and "Runeblade" not in c["Types"]]
    by_name = {}
    for c in generics:
        by_name.setdefault(c["Name"], []).append(c)
    return by_name


def fmt_text(text):
    """Escape text for use as Go string in docstring comments. We'll word-wrap at 100 chars."""
    if not text:
        return ""
    # Collapse runs of whitespace (incl newlines) to single space for docstring transcription.
    # Keep paragraph breaks by splitting on double-newline first.
    paragraphs = [re.sub(r"\s+", " ", p).strip() for p in text.split("\n\n")]
    return "\n\n".join(p for p in paragraphs if p)


def wrap_paragraph(text, prefix="// ", width=100):
    """Word-wrap a single paragraph into multi-line // comment. No trailing //."""
    words = text.split()
    if not words:
        return ""
    out_lines = []
    cur = prefix + words[0]
    for w in words[1:]:
        cand = cur + " " + w
        if len(cand) > width:
            out_lines.append(cur)
            cur = prefix + w
        else:
            cur = cand
    out_lines.append(cur)
    return "\n".join(out_lines)


def build_doc_comment(paragraphs):
    """Build a Go doc comment block given a list of paragraph strings. Blank paragraphs become '//'."""
    out = []
    for i, p in enumerate(paragraphs):
        if p == "":
            out.append("//")
        else:
            out.append(wrap_paragraph(p))
    return "\n".join(out)


# --------------------------------------------------------------------------
# Per-card Play behaviour table. Key is card name; value is dict:
#   "play": either "attack" (return c.Attack()), "zero" (return 0), or a
#           snippet of custom code for the Play body.
#   "aura": True if the card is a Generic Action Aura — Play sets AuraCreated.
#   "simp": optional string added to the Simplification section.
# --------------------------------------------------------------------------

# Cards where we model the rider. name -> (play_body_code, simplification_note)
# The code snippet must return an int and may reference `c` (card receiver) and `s` (state).
# Helpers are defined per card family in a shared var/func above the structs.

MODELED = {
    # Yinti Yanti: aura on board → +1p (optimistic).
    "Yinti Yanti": dict(
        helper="""// yintiYantiPlay adds +1 when any aura is in play: either created this turn or played earlier.
func yintiYantiPlay(base int, s *card.TurnState) int {
\tif s != nil && (s.AuraCreated || s.HasPlayedType(card.TypeAura)) {
\t\treturn base + 1
\t}
\treturn base
}
""",
        call="yintiYantiPlay(c.Attack(), s)",
        simp="The defending-side +1{d} buff is ignored (defence is consumed before Play).",
    ),
    # Vigor Rush: if played a non-attack action this turn → go again.
    "Vigor Rush": dict(
        helper="""// vigorRushPlay grants go again when any non-attack Action has been played earlier this turn.
func vigorRushPlay(base int, s *card.TurnState) int {
\tif s == nil || s.Self == nil {
\t\treturn base
\t}
\tfor _, pl := range s.CardsPlayed {
\t\tt := pl.Types()
\t\tif t.Has(card.TypeAction) && !t.Has(card.TypeAttack) {
\t\t\ts.Self.GrantedGoAgain = true
\t\t\tbreak
\t\t}
\t}
\treturn base
}
""",
        call="vigorRushPlay(c.Attack(), s)",
        simp="",
    ),
    # Zealous Belting: any pitched card with power > base power → go again (printed says "in pitch
    # zone with {p} greater than this" — we read s.Pitched).
    "Zealous Belting": dict(
        helper="""// zealousBeltingPlay grants go again when any pitched card this turn has base power greater
// than the card's own base power.
func zealousBeltingPlay(base int, s *card.TurnState) int {
\tif s == nil || s.Self == nil {
\t\treturn base
\t}
\tfor _, p := range s.Pitched {
\t\tif p.Attack() > base {
\t\t\ts.Self.GrantedGoAgain = true
\t\t\tbreak
\t\t}
\t}
\treturn base
}
""",
        call="zealousBeltingPlay(c.Attack(), s)",
        simp="",
    ),
    # Come to Fight, Minnowism, Nimblism, Sloggism, Captain's Call, Warmonger's Recital,
    # Money Where Ya Mouth Is, Prime the Crowd: next attack card +Np. We share a helper that
    # accepts a filter predicate and a bonus.
}

# Non-attack action "next attack card +Np" with optional cost filter.
# name -> (bonus, filter name, note)
# filter names: "any", "cost<=2", "cost<=1", "cost>=2", "power<=3"
NEXT_ATTACK_BONUS = {
    "Come to Fight":         (3, "any",      ""),
    "Minnowism":             (3, "power<=3", ""),
    "Nimblism":              (3, "cost<=1",  ""),
    "Sloggism":              (6, "cost>=2",  ""),
    "Captain's Call":        (2, "cost<=2",  "Modal: we pick the +2 power mode; the alternative 'go again' mode is dropped."),
    "Warmonger's Recital":   (3, "any",      "The 'bottom of deck' rider is dropped (just credit the +3)."),
    "Money Where Ya Mouth Is": (3, "any",    "Wager Gold token rider is dropped."),
    "Prime the Crowd":       (4, "any",      "Crowd cheers/boos keywords are dropped."),
    "Scout the Periphery":   (3, "any",      "The 'play from arsenal' gate is ignored (arsenal not modelled)."),
    "Plunder Run":           (3, "any",      "Draw rider on hit is dropped; the arsenal-only +3 is credited unconditionally."),
    "Smashing Good Time":    (3, "any",      "Item-destruction rider ignored; arsenal-only +3 credited unconditionally."),
    "Force Sight":           (3, "any",      "The arsenal-gated Opt 2 is skipped."),
    "Regain Composure":      (1, "any",      "The on-hit unfreeze rider is dropped."),
    "Clearwater Elixir":     (3, "any",      "Bloodrot Pox health-gain rider dropped."),
    "Restvine Elixir":       (3, "any",      "Inertia health-gain rider dropped."),
    "Sapwood Elixir":        (3, "any",      "Frailty health-gain rider dropped."),
    "Public Bounty":         (3, "any",      "Mark isn't modelled; the +3 rider is credited unconditionally."),
}

# "Next attack gets Go Again" non-attack actions. name -> (filter, note)
NEXT_ATTACK_GRANT_GA = {
    "Flying High":     ("any", "The '+1{p} if red' rider is modelled: when the granted target has pitch 1, +1 is also credited."),
    "Trot Along":      ("power<=3", ""),
}

# Plain non-attack actions whose riders we skip (Play returns 0, card still a resource).
PLAIN_NONATTACK_ACTIONS = {
    "Healing Balm":          "Gain 3{h} is dropped (hero health isn't tracked).",
    "Starting Stake":        "Gold-token economy isn't tracked.",
    "Sift":                  "Hand cycling isn't modelled.",
    "Strategic Planning":    "Graveyard recovery and end-phase draw aren't modelled.",
    "Whisper of the Oracle": "Opt isn't modelled.",
    "Pick a Card, Any Card": "Opponent hand inspection and Silver-token economy aren't modelled.",
    "Sun Kiss":              "Health gain and the Moon Wish synergy aren't modelled.",
    "Cash In":               "Activated Gold/Silver/Copper economy and draws aren't modelled.",
    "Lead the Charge":       "Action-point tracking isn't modelled.",
    "High Striker":          "Copper-token economy isn't modelled.",
    "Ransack and Raze":      "Landmarks and Gold tokens aren't modelled; X cost treated as 0.",
    "Relentless Pursuit":    "Marks and 'attacked them this turn' tracking aren't modelled.",
    "Tit for Tat":           "Freeze/unfreeze (tap/untap) isn't modelled.",
    "Visit the Blacksmith":  "Next-sword-attack bonuses aren't applied (weapon chain isn't peeked).",
    "On a Knife Edge":       "Next-sword-attack go-again isn't applied (weapon chain isn't peeked).",
}

# Aura cards — Play sets AuraCreated and returns 0. Key -> simplification note.
AURAS = {
    "Enchanting Melody":   "Damage-prevention trigger and end-phase destruction clause aren't modelled; value is just the aura-created flag (read by Yinti Yanti, Runerager Swarm, etc.).",
    "Sigil of Cycles":     "At-action-phase self-destroy and leaves-arena discard/draw aren't modelled; only the aura-created flag is credited.",
    "Sigil of Fyendal":    "At-action-phase self-destroy and the 1{h} gain on leave are dropped; only the aura-created flag is credited.",
    "Sigil of Protection": "Ward N isn't modelled (opponent's damage prevention); only the aura-created flag is credited.",
}

# Block cards — behave as DR with printed defense only.
BLOCKS_NO_RIDER = {
    "Fiddler's Green":  "The 'gain 3{h} on entering graveyard' trigger isn't modelled.",
    "On the Horizon":   "The deck-peek trigger isn't modelled.",
    "Test of Strength": "Clash with the attacking hero and the Gold-token winner rider aren't modelled.",
}

# Plain Defense Reactions (no riders we model).
DR_SKIPPED = {
    "Drag Down":        "The -3{p} attacker debuff isn't modelled (solver doesn't expose defender-side power reductions).",
}


def classify_attack(name, text):
    """For Action-Attack cards, decide whether rider is modelable. Returns tuple
    (play_mode, simplification_note). play_mode:
      'attack' — just return c.Attack().
      'custom:<key>' — use MODELED[key].
    """
    if name in MODELED:
        return "custom", MODELED[name].get("simp", "")
    # Everything else: attack only, with a description-specific simplification noting what's ignored.
    return "attack", attack_simp_for(name, text)


# --------------------------------------------------------------------------
# Attack rider simplification notes (what we skip for each attack card).
# If empty, use a generic note. These are free-form prose describing the rider
# that isn't being modelled.
# --------------------------------------------------------------------------
ATTACK_SIMPS = {
    "Adrenaline Rush":       "'Less life than opposing hero' health comparison isn't modelled; the +3{p} rider never fires.",
    "Back Alley Breakline":  "Face-up-from-deck action-point grant isn't modelled.",
    "Barraging Brawnhide":   "Defended-by-<2-non-equipment condition isn't modelled; the +1{p} rider never applies.",
    "Battlefront Bastion":   "Defense-prevention rider isn't modelled.",
    "Belittle":              "Additional-cost reveal and deck search for Minnowism aren't modelled.",
    "Blanch":                "Opponent 'lose all colors' debuff isn't modelled.",
    "Blow for a Blow":       "Health comparison for go-again and on-hit 1 damage aren't modelled.",
    "Bluster Buff":          "Pay {r} or lose 1{p} — we keep base power; players who can't afford should be rare and the loss is 1.",
    "Brandish":              "Next weapon attack +1{p} isn't modelled (weapons aren't scanned in CardsRemaining).",
    "Brothers in Arms":      "Pay-to-buff-defence rider isn't modelled (defence-side costs aren't solved).",
    "Brutal Assault":        "",
    "Cadaverous Contraband": "Graveyard → top-of-deck rider isn't modelled.",
    "Chest Puff":            "Pay {r} or lose 1{p} — base power is kept.",
    "Crash Down the Gates":  "Deck-reveal comparison and top-of-deck destruction aren't modelled.",
    "Critical Strike":       "",
    "Cut Down to Size":      "Opponent discard rider isn't modelled.",
    "Demolition Crew":       "Additional reveal cost isn't modelled; Dominate keyword is held but unused.",
    "Destructive Deliberation": "Ponder token creation isn't modelled.",
    "Down But Not Out":      "Health/equipment/token comparison isn't modelled; none of the riders fire.",
    "Drone of Brutality":    "Graveyard-replacement-to-deck trigger isn't modelled.",
    "Emissary of Moon":      "Hand-cycle draw rider isn't modelled.",
    "Emissary of Tides":     "Hand-cycle-for-+2{p} rider isn't modelled.",
    "Emissary of Wind":      "Hand-cycle-for-go-again rider isn't modelled.",
    "Fact-Finding Mission":  "Peeking opponent arsenal/equipment isn't modelled.",
    "Feisty Locals":         "The 'defended by action card' +2{p} rider isn't modelled.",
    "Fervent Forerunner":    "On-hit Opt 2 and arsenal-only go-again aren't modelled.",
    "Flex":                  "Pay-{r}{r}-for-+2{p} rider isn't modelled.",
    "Flock of the Feather Walkers": "Additional reveal cost and Quicken token creation aren't modelled.",
    "Freewheeling Renegades": "The 'defended by action card' -2{p} rider isn't modelled.",
    "Frontline Scout":       "Hand-peek and arsenal-only go-again aren't modelled.",
    "Fyendal's Fighting Spirit": "Conditional health-gain rider isn't modelled.",
    "Gravekeeping":          "Graveyard banish rider isn't modelled.",
    "Hand Behind the Pen":   "Arsenal-manipulation rider isn't modelled.",
    "Humble":                "Hero-ability suppression rider isn't modelled.",
    "Infectious Host":       "Frailty/Inertia/Bloodrot Pox token creation isn't modelled.",
    "Jack Be Nimble":        "Graveyard banish for +1{p}/go-again and on-hit steal aren't modelled.",
    "Jack Be Quick":         "Graveyard banish for +1{p}/go-again and on-hit steal aren't modelled.",
    "Life for a Life":       "Health comparison for go-again and on-hit 1{h} gain aren't modelled.",
    "Life of the Party":     "Crazy Brew substitute cost and random-mode selection aren't modelled.",
    "Look Tuff":             "Pay {r} or lose 1{p} — base power is kept.",
    "Looking for a Scrap":   "Graveyard-banish additional cost and bonus rider aren't modelled.",
    "Money or Your Life?":   "Gold-token exchange rider isn't modelled.",
    "Moon Wish":             "Alternative hand-on-top cost and Sun Kiss search rider aren't modelled.",
    "Muscle Mutt":           "",
    "Nimble Strike":         "Graveyard-banish additional cost and bonus rider aren't modelled.",
    "Nimby":                 "Deck search for Nimblism isn't modelled.",
    "Outed":                 "Marked-hero checks aren't modelled; the +1{p} rider never applies.",
    "Out Muscle":            "Defended-by-equal-or-greater-power gate isn't modelled; the printed Go again keyword is kept.",
    "Overload":              "On-hit go-again and Dominate aren't modelled (keyword held but solver ignores hit-gated grants).",
    "Performance Bonus":     "Gold-token creation and arsenal-only go-again aren't modelled.",
    "Pound for Pound":       "Health comparison for Dominate isn't modelled; Dominate keyword held but unused.",
    "Promise of Plenty":     "Arsenal-placement rider and arsenal-only go-again aren't modelled.",
    "Punch Above Your Weight": "Pay-{r}{r}{r}-for-+5{p} rider isn't modelled.",
    "Pursue to the Edge of Oblivion": "Mark on hit isn't modelled.",
    "Pursue to the Pits of Despair":  "Mark on hit isn't modelled.",
    "Push the Point":        "Chain-history +2{p} rider isn't modelled.",
    "Raging Onslaught":      "",
    "Rally the Coast Guard": "Defense-time instant activated ability isn't modelled.",
    "Rally the Rearguard":   "Defense-time instant activated ability isn't modelled.",
    "Ravenous Rabble":       "Deck-reveal -X{p} rider isn't modelled; base power is returned.",
    "Right Behind You":      "Defend-together +1{d} and deck-bottom rider aren't modelled.",
    "Rifting":               "On-hit instant-casting rider isn't modelled.",
    "Scar for a Scar":       "Health comparison for go-again isn't modelled.",
    "Scour the Battlescape": "Hand-cycle and arsenal-only go-again aren't modelled.",
    "Seek Horizon":          "Hand-on-top additional cost and conditional go-again aren't modelled.",
    "Sirens of Safe Harbor": "Graveyard-trigger 1{h} gain isn't modelled.",
    "Smash Up":              "Arsenal-manipulation rider isn't modelled.",
    "Snatch":                "Draw on hit isn't modelled.",
    "Sound the Alarm":       "Deck search rider isn't modelled.",
    "Spring Load":           "'No cards in hand' +3{p} rider isn't modelled.",
    "Stony Woottonhog":      "Defended-by-<2-non-equipment condition isn't modelled.",
    "Strike Gold":           "Gold-token creation isn't modelled.",
    "Surging Militia":       "Defended-by +N{p} rider isn't modelled.",
    "Tip-Off":               "Instant discard activation isn't modelled.",
    "Tongue Tied":           "Arsenal-manipulation rider isn't modelled.",
    "Trade In":              "Discard-to-draw and arsenal-only go-again aren't modelled.",
    "Tremor of \u00edArathael": "Banished-zone tracking isn't modelled; +2{p} rider never fires.",
    "Wage Gold":             "Universal keyword and Gold-token wager aren't modelled.",
    "Walk the Plank":        "Pirate-specific target-freezing rider isn't modelled.",
    "Water the Seeds":       "Chain-bonus rider for next low-power attack isn't modelled.",
    "Wounded Bull":          "Health comparison isn't modelled; +1{p} rider never fires.",
    "Wounding Blow":         "",
    "Wreck Havoc":           "Defense-reaction lockout and arsenal-banish aren't modelled.",
}


def attack_simp_for(name, text):
    return ATTACK_SIMPS.get(name, "Riders described above aren't modelled; Play returns base power.")


def build_file(name, kind, printings, text):
    """Return (file_name_without_ext, file_content)."""
    ident_root = to_identifier(name)
    file_root = to_filename(name)

    # Sort printings Red → Yellow → Blue
    order = {"Red": 0, "Yellow": 1, "Blue": 2}
    printings = sorted(printings, key=lambda c: order.get(c["Color"], 9))

    colors = [p["Color"] for p in printings]
    struct_names = [f"{ident_root}{c}" for c in colors]

    # Header docstring.
    typestring_map = {
        "Generic, Action, Attack":  "Generic Action - Attack",
        "Generic, Action":          "Generic Action",
        "Generic, Action, Aura":    "Generic Action - Aura",
        "Generic, Defense Reaction": "Generic Defense Reaction",
        "Generic, Block":           "Generic Block",
    }
    type_label = typestring_map[kind]

    # Base stats summary for the first printing (Red) plus variant table.
    p0 = printings[0]
    cost = p0.get("Cost", "0") or "0"
    df = p0.get("Defense", "0") or "0"
    # If multiple printings vary power, we'll list variants.
    powers = [(p["Color"], p.get("Power", "0") or "0") for p in printings]
    pitches = [(p["Color"], p.get("Pitch", "0") or "0") for p in printings]
    defenses = [(p["Color"], p.get("Defense", "0") or "0") for p in printings]

    # Build the header summary line(s) — everything up to the Text: block.
    if len(printings) == 1:
        p = printings[0]
        parts = [f"Cost {cost}"]
        parts.append(f"Pitch {p.get('Pitch','0') or '0'}")
        if kind == "Generic, Action, Attack":
            parts.append(f"Power {p.get('Power','0') or '0'}")
        parts.append(f"Defense {p.get('Defense','0') or '0'}")
        color_note = f" Only printed in {p['Color']}." if p.get('Color') else ""
        summary = f"{name} \u2014 {type_label}. " + ", ".join(parts) + "." + color_note
    else:
        summary_bits = [f"{name} \u2014 {type_label}."]
        summary_bits.append(f"Cost {cost}.")
        if kind == "Generic, Action, Attack":
            summary_bits.append("Printed power: " + ", ".join(f"{c} {v}" for c, v in powers) + ".")
        summary_bits.append("Printed pitch variants: " + ", ".join(f"{c} {v}" for c, v in pitches) + ".")
        if len(set(v for _, v in defenses)) > 1:
            summary_bits.append("Printed defense: " + ", ".join(f"{c} {v}" for c, v in defenses) + ".")
        else:
            summary_bits.append(f"Defense {defenses[0][1]}.")
        summary = " ".join(summary_bits)

    # Paragraphs that will go into the doc comment. "" means a blank separator line (//).
    paragraphs = [summary]
    text_para = fmt_text(text)
    if text_para:
        paragraphs.append("")
        paragraphs.append(f"Text: \"{text_para}\"")

    # Pick the Play mode + build the simplification note.
    aura = kind == "Generic, Action, Aura"
    block = kind == "Generic, Block"
    is_attack = kind == "Generic, Action, Attack"
    is_nonattack = kind == "Generic, Action"
    is_dr = kind == "Generic, Defense Reaction"

    helper_code = ""
    play_call = None   # if None, auto-generate default for kind
    simp_note = ""

    if is_attack:
        mode, s = classify_attack(name, text)
        simp_note = s
        if mode == "custom":
            m = MODELED[name]
            helper_code = m["helper"]
            play_call = m["call"]
        else:
            play_call = "c.Attack()"
    elif is_nonattack:
        if name in NEXT_ATTACK_BONUS:
            bonus, filt, extra = NEXT_ATTACK_BONUS[name]
            helper_code, play_call = next_attack_bonus_helper(ident_root, bonus, filt)
            simp_note = build_next_attack_simp(extra)
        elif name in NEXT_ATTACK_GRANT_GA:
            filt, extra = NEXT_ATTACK_GRANT_GA[name]
            helper_code, play_call = next_attack_grant_ga_helper(name, ident_root, filt)
            simp_note = extra
        elif name in PLAIN_NONATTACK_ACTIONS:
            simp_note = PLAIN_NONATTACK_ACTIONS[name]
            play_call = "0"
        else:
            simp_note = "Rider isn't modelled."
            play_call = "0"
    elif aura:
        simp_note = AURAS.get(name, "Rider isn't modelled.")
        play_call = "auraPlay(s)"
        # auraPlay is a shared helper defined once in the package. We'll install it into dodge.go's
        # var block? No — we'll define an internal helper in this file, or share via _shared.go.
        helper_code = f"""// {ident_root.lower()[0]}{ident_root[1:] if False else ''}
"""
        # We'll actually just inline the aura set directly; define a tiny local func.
        helper_code = """"""
        play_call = "setAuraCreated(s)"
    elif block:
        simp_note = BLOCKS_NO_RIDER.get(name, "Rider isn't modelled.")
        play_call = "0"
    elif is_dr:
        simp_note = DR_SKIPPED.get(name, "Rider isn't modelled.")
        play_call = "0"

    if simp_note:
        paragraphs.append("")
        paragraphs.append(f"Simplification: {simp_note}")
    paragraphs.append("")
    paragraphs.append("Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).")
    doc_comment = build_doc_comment(paragraphs) + "\n"

    # Types bitfield + var.
    if is_attack:
        types_expr = "card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)"
        types_var = f"{ident_root[0].lower()}{ident_root[1:]}Types"
    elif is_nonattack:
        types_expr = "card.NewTypeSet(card.TypeGeneric, card.TypeAction)"
        types_var = f"{ident_root[0].lower()}{ident_root[1:]}Types"
    elif aura:
        types_expr = "card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAura)"
        types_var = f"{ident_root[0].lower()}{ident_root[1:]}Types"
    elif block or is_dr:
        # model blocks as DR
        types_expr = None
        types_var = "defenseReactionTypes"  # shared from dodge.go
    else:
        raise RuntimeError(f"unknown kind {kind} for {name}")

    # Build the struct block for each printing.
    struct_blocks = []
    for p in printings:
        color = p["Color"]
        sname = f"{ident_root}{color}"
        id_name = f"{ident_root}{color}"
        pitch = int(p.get("Pitch") or "0")
        pcost_raw = p.get("Cost") or "0"
        pcost = 0 if pcost_raw in ("", "X") else int(pcost_raw)
        power = int(p.get("Power") or "0")
        defense = int(p.get("Defense") or "0")
        kw = p.get("CardKeywords", "")
        go_again = "true" if "Go again" in kw else "false"

        # Need to decide whether Play gets a receiver (c) or not.
        use_receiver = play_call and ("c." in play_call)
        recv = "c" if use_receiver else ""
        if recv:
            play_sig = f"func ({recv} {sname}) Play(s *card.TurnState) int"
        else:
            play_sig = f"func ({sname}) Play(s *card.TurnState) int"

        # If play doesn't use s, avoid "unused var"?
        # Actually in Go s is just a param; unused params are OK.

        block_lines = [
            f"type {sname} struct{{}}",
            "",
            f"func ({sname}) ID() card.ID                 {{ return card.{id_name} }}",
            f"func ({sname}) Name() string                {{ return {json_str(name, color)} }}",
            f"func ({sname}) Cost(*card.TurnState) int    {{ return {pcost} }}",
            f"func ({sname}) Pitch() int                  {{ return {pitch} }}",
            f"func ({sname}) Attack() int                 {{ return {power} }}",
            f"func ({sname}) Defense() int                {{ return {defense} }}",
            f"func ({sname}) Types() card.TypeSet         {{ return {types_var} }}",
            f"func ({sname}) GoAgain() bool               {{ return {go_again} }}",
            f"{play_sig} {{ return {play_call} }}",
        ]
        struct_blocks.append("\n".join(block_lines))

    # Assemble file.
    parts = [doc_comment.rstrip() + "\n", "package generic", ""]
    parts.append('import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"')
    parts.append("")
    if types_expr:
        parts.append(f"var {types_var} = {types_expr}")
        parts.append("")
    if helper_code.strip():
        parts.append(helper_code.rstrip())
        parts.append("")
    parts.append("\n\n".join(struct_blocks))
    parts.append("")
    return file_root, "\n".join(parts)


def json_str(name, color):
    # return a Go string literal for "Name (Color)"
    s = f"{name} ({color})"
    return "\"" + s.replace("\\", "\\\\").replace("\"", "\\\"") + "\""


def next_attack_bonus_helper(ident_root, bonus, filt):
    """Return (helper_code, call_expr) for a 'next attack +N' non-attack action."""
    var = f"{ident_root[0].lower()}{ident_root[1:]}"
    fn = f"{var}Play"
    pred_map = {
        "cost<=2":  "pc.Card.Cost(s) <= 2",
        "cost<=1":  "pc.Card.Cost(s) <= 1",
        "cost>=2":  "pc.Card.Cost(s) >= 2",
        "power<=3": "pc.Card.Attack() <= 3",
    }
    if filt == "any":
        body = f"return {bonus}"
    else:
        pred = pred_map[filt]
        body = f"if {pred} {{\n\t\t\treturn {bonus}\n\t\t}}\n\t\tcontinue"
    code = f"""// {fn} returns {bonus} when a matching attack action card is scheduled later this turn.
func {fn}(s *card.TurnState) int {{
\tfor _, pc := range s.CardsRemaining {{
\t\tt := pc.Card.Types()
\t\tif !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {{
\t\t\tcontinue
\t\t}}
\t\t{body}
\t}}
\treturn 0
}}
"""
    return code, f"{fn}(s)"


def next_attack_grant_ga_helper(name, ident_root, filt):
    var = f"{ident_root[0].lower()}{ident_root[1:]}"
    fn = f"{var}Play"
    pred_map = {
        "cost<=2":  "pc.Card.Cost(s) <= 2",
        "cost<=1":  "pc.Card.Cost(s) <= 1",
        "power<=3": "pc.Card.Attack() <= 3",
    }
    if name == "Flying High":
        # Grant +1 if target has pitch 1 (red).
        body = (
            "bonus := 0\n"
            "\t\tif pc.Card.Pitch() == 1 {\n"
            "\t\t\tbonus = 1\n"
            "\t\t}\n"
            "\t\tpc.GrantedGoAgain = true\n"
            "\t\treturn bonus"
        )
        doc = (f"// {fn} grants go again to the next attack action card scheduled later this turn. For "
               "Flying High, if\n"
               f"// that target is red (pitch 1) we also credit +1 power as a bonus.")
    else:
        body = "pc.GrantedGoAgain = true\n\t\treturn 0"
        doc = f"// {fn} grants go again to the next qualifying attack action card scheduled later this turn."
    if filt == "any":
        wrap = body
    else:
        pred = pred_map[filt]
        wrap = f"if {pred} {{\n\t\t\t" + body.replace("\n\t\t", "\n\t\t\t") + "\n\t\t}"
    code = f"""{doc}
func {fn}(s *card.TurnState) int {{
\tfor _, pc := range s.CardsRemaining {{
\t\tt := pc.Card.Types()
\t\tif !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {{
\t\t\tcontinue
\t\t}}
\t\t{wrap}
\t}}
\treturn 0
}}
"""
    return code, f"{fn}(s)"


def build_next_attack_simp(extra):
    base = ("Scans TurnState.CardsRemaining for the first matching attack action card and credits the "
            "bonus assuming it will be played; if none is scheduled after this card, the bonus fizzles.")
    if extra:
        return extra + " " + base
    return base


def main():
    by_name = load_cards()
    # Decide which files to emit.
    to_emit = []
    for name, printings in sorted(by_name.items()):
        t = printings[0]["Types"]
        if t not in IN_SCOPE_TYPES:
            continue
        if name in ALREADY_DONE:
            continue
        to_emit.append((name, t, printings, printings[0]["FunctionalText"]))

    gen_dir = os.path.join(ROOT, "internal", "card", "generic")
    ids_and_structs = []   # list of (ident_root, color) for ID and byID registration
    for name, t, printings, text in to_emit:
        file_root, content = build_file(name, t, printings, text)
        path = os.path.join(gen_dir, file_root + ".go")
        if os.path.exists(path):
            # don't clobber existing (hand-written) files
            print("SKIP existing:", path, file=sys.stderr)
            continue
        with open(path, "w", encoding="utf-8", newline="\n") as f:
            f.write(content)
        print("wrote", path)
        ident = to_identifier(name)
        for p in sorted(printings, key=lambda c: {"Red":0,"Yellow":1,"Blue":2}.get(c["Color"],9)):
            ids_and_structs.append((name, ident, p["Color"]))

    # Emit shared helper (setAuraCreated) if any aura emitted.
    emitted_aura = any(t == "Generic, Action, Aura" for (_, t, _, _) in to_emit)
    if emitted_aura:
        helper_path = os.path.join(gen_dir, "aura_helper.go")
        helper_content = (
            "// Shared aura-creation helper for Generic Action - Aura cards.\n"
            "//\n"
            "// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).\n"
            "\n"
            "package generic\n"
            "\n"
            "import \"github.com/tim-chaplin/fab-deck-optimizer/internal/card\"\n"
            "\n"
            "// setAuraCreated marks the turn state so cards that read AuraCreated (e.g. Yinti Yanti,\n"
            "// Runerager Swarm) see the aura entering play. Returns 0 — the aura itself contributes\n"
            "// no direct damage; its value is in the flag it leaves behind.\n"
            "func setAuraCreated(s *card.TurnState) int {\n"
            "\ts.AuraCreated = true\n"
            "\treturn 0\n"
            "}\n"
        )
        if not os.path.exists(helper_path):
            with open(helper_path, "w", encoding="utf-8", newline="\n") as f:
                f.write(helper_content)
            print("wrote", helper_path)

    # Emit a separate file listing generated IDs + registry insertions to integrate manually.
    id_snippet_path = os.path.join(TMP, "silverage_id_snippet.txt")
    reg_snippet_path = os.path.join(TMP, "silverage_registry_snippet.txt")
    with open(id_snippet_path, "w", encoding="utf-8", newline="\n") as f:
        last_name = None
        for name, ident, color in ids_and_structs:
            if name != last_name and last_name is not None:
                f.write("\n")  # blank line between families
            f.write(f"\t{ident}{color}\n")
            last_name = name
    with open(reg_snippet_path, "w", encoding="utf-8", newline="\n") as f:
        last_name = None
        for name, ident, color in ids_and_structs:
            if name != last_name and last_name is not None:
                f.write("\n")
            f.write(f"\tcard.{ident}{color}: generic.{ident}{color}{{}},\n")
            last_name = name
    print("ID snippet:", id_snippet_path)
    print("Registry snippet:", reg_snippet_path)


if __name__ == "__main__":
    main()
