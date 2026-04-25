"""Strip the framework-restating "Stub only — marked NotImplemented..." paragraph from each
generated stub and replace it with a card-specific one-line `// not implemented: <quirk>`
comment immediately above the `NotImplemented()` method, matching the existing convention used
by hand-written stubs (e.g. internal/card/generic/back_alley_breakline.go).

Run after gen_notimpl_stubs.py and merge_notimpl_ids.py:

  python scripts/fix_notimpl_comments.py

Idempotent: if a file already lacks the boilerplate paragraph, it's left alone.
"""
import os
import re

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# One-line, card-specific summary of what's NOT modelled for each stub. Each maps card name →
# short clause (≤100 chars) inserted as `// not implemented: <clause>` above NotImplemented().
NOTES = {
    "Amulet of Assertiveness": "AR 'banish-top-of-deck-on-hit' grant; gated on hand size",
    "Amulet of Echoes": "Instant 'discard 2' against repeat-name plays",
    "Amulet of Havencall": "DR tutor for Rally the Rearguard; gated on empty hand",
    "Amulet of Ignition": "Instant 'next ability costs {r} less'",
    "Amulet of Intervention": "Instant 1 damage prevention; gated on lethal source",
    "Amulet of Oblation": "Instant 'graveyard → bottom of deck' replacement effect",
    "Arcane Polarity": "1{h} or 4/3/2{h} on arcane-damage trigger",
    "Blade Flash": "AR 'sword attack gains go again' grant",
    "Brush Off": "Instant prevent 1 damage; gated on aura/item with counter",
    "Calming Breeze": "Instant 1{h} gain",
    "Clap 'Em in Irons": "passive Pirate-target tap rider; self-destroys at upkeep",
    "Clarity Potion": "activated 'next instant costs {r} less'",
    "Count Your Blessings": "graveyard-scaled X{h} gain (also banlisted)",
    "Cracked Bauble": "draft-format pitch resource; no other effect",
    "Destructive Tendencies": "Instant remove counters from item / aura tokens",
    "Eirina's Prayer": "Instant prevent 2 damage to a non-hero target",
    "Energy Potion": "activated 'put 2 cards from graveyard on bottom, gain action'",
    "Even Bigger Than That!": "Opt + reveal-and-quicken trigger; gated on damage dealt",
    "Exposed": "AR -2{p} attacker debuff; gated on hand size",
    "Fool's Gold": "discard trigger creates Gold token",
    "Healing Potion": "activated 2{h} gain",
    "Imperial Seal of Command": "activated 'no DR this turn' lockout + Royal-only arsenal-destroy on hit",
    "Lunging Press": "AR +2{p} buff and freeze rider",
    "Memorial Ground": "Instant 'graveyard → top of deck' for low-cost attack action",
    "Nip at the Heels": "AR +1{p} buff and on-hit draw",
    "Oasis Respite": "Instant 1{h} gain to a non-hero target",
    "Peace of Mind": "Instant prevent 2 damage; gated on no attacks this turn",
    "Pilfer the Tomb": "Instant 'banish from opposing graveyard'",
    "Potion of Déjà Vu": "activated 'replay last instant from graveyard'",
    "Potion of Ironhide": "activated +2{d} on next defending card",
    "Potion of Luck": "activated peek-and-rearrange top 3",
    "Potion of Seeing": "activated reveal opponent's hand",
    "Potion of Strength": "activated +2{p} on next attack",
    "Pummel": "AR 'destroy target item with ≤2 counters' rider",
    "Razor Reflex": "modal AR +N{p} for sword/dagger or low-cost attack action",
    "Reinforce the Line": "Instant +1{d} grant to a defending card",
    "Shatter Sorcery": "Instant destroy Sigil aura or prevent 1 arcane damage",
    "Sigil of Solace": "1/2/3{h} gain (also banlisted)",
    "Talisman of Balance": "passive 'gain 1{h} when you draw 4+'",
    "Talisman of Cremation": "passive 'banish on graveyard entry'",
    "Talisman of Dousing": "passive prevent 1 arcane damage per turn",
    "Talisman of Featherfoot": "passive 'first attack each turn gains evade'",
    "Talisman of Recompense": "passive 1{h} on opponent damage",
    "Talisman of Tithes": "passive Gold-token economy on opposing damage",
    "Talisman of Warfare": "passive arsenal-wipe on a 2-damage hit",
    "Thrust": "AR +3{p} buff to a sword attack",
    "Timesnap Potion": "activated 'play next attack action from graveyard'",
    "Titanium Bauble": "Defense Reaction with 3{d}; no other effect",
    "Talishar, the Lost Prince": "rust-counter activation cost and self-destruct trigger",
}


def process_file(path, name):
    src = open(path, encoding="utf-8").read()
    note = NOTES.get(name)
    if note is None:
        print(f"NO NOTE for {name}; skipping {path}")
        return False

    # 1. Strip the "Stub only — ..." paragraph from the doc comment block. The paragraph runs
    #    from the line that begins "// Stub only —" through the last consecutive `//`-prefixed
    #    line. The blank `//` separator above it is removed too.
    new_src = re.sub(
        r"//\n// Stub only —[^\n]*(?:\n//[^\n]*)*\n",
        "",
        src,
        count=1,
    )

    # 2. Insert a `// not implemented: <note>` line immediately above each NotImplemented()
    #    method. Idempotent: skip if a "// not implemented:" line already sits directly above.
    pattern = re.compile(
        r"(^[ \t]*)func \(([A-Za-z0-9]+)\) NotImplemented\(\)",
        re.MULTILINE,
    )

    def repl(m):
        indent = m.group(1)
        # Look one line back: if the previous line is already a "not implemented:" comment, skip.
        start = m.start()
        # Find the previous newline.
        prev_nl = new_src.rfind("\n", 0, start - 1)
        prev_line = new_src[prev_nl + 1 : start - 1] if prev_nl != -1 else ""
        if "not implemented:" in prev_line:
            return m.group(0)
        return f"{indent}// not implemented: {note}\n{m.group(0)}"

    new_src = pattern.sub(repl, new_src)

    if new_src == src:
        return False
    with open(path, "w", encoding="utf-8", newline="\n") as f:
        f.write(new_src)
    return True


CARD_FILES = {
    "amulet_of_assertiveness.go": "Amulet of Assertiveness",
    "amulet_of_echoes.go": "Amulet of Echoes",
    "amulet_of_havencall.go": "Amulet of Havencall",
    "amulet_of_ignition.go": "Amulet of Ignition",
    "amulet_of_intervention.go": "Amulet of Intervention",
    "amulet_of_oblation.go": "Amulet of Oblation",
    "arcane_polarity.go": "Arcane Polarity",
    "blade_flash.go": "Blade Flash",
    "brush_off.go": "Brush Off",
    "calming_breeze.go": "Calming Breeze",
    "clap_em_in_irons.go": "Clap 'Em in Irons",
    "clarity_potion.go": "Clarity Potion",
    "count_your_blessings.go": "Count Your Blessings",
    "cracked_bauble.go": "Cracked Bauble",
    "destructive_tendencies.go": "Destructive Tendencies",
    "eirinas_prayer.go": "Eirina's Prayer",
    "energy_potion.go": "Energy Potion",
    "even_bigger_than_that.go": "Even Bigger Than That!",
    "exposed.go": "Exposed",
    "fools_gold.go": "Fool's Gold",
    "healing_potion.go": "Healing Potion",
    "imperial_seal_of_command.go": "Imperial Seal of Command",
    "lunging_press.go": "Lunging Press",
    "memorial_ground.go": "Memorial Ground",
    "nip_at_the_heels.go": "Nip at the Heels",
    "oasis_respite.go": "Oasis Respite",
    "peace_of_mind.go": "Peace of Mind",
    "pilfer_the_tomb.go": "Pilfer the Tomb",
    "potion_of_deja_vu.go": "Potion of Déjà Vu",
    "potion_of_ironhide.go": "Potion of Ironhide",
    "potion_of_luck.go": "Potion of Luck",
    "potion_of_seeing.go": "Potion of Seeing",
    "potion_of_strength.go": "Potion of Strength",
    "pummel.go": "Pummel",
    "razor_reflex.go": "Razor Reflex",
    "reinforce_the_line.go": "Reinforce the Line",
    "shatter_sorcery.go": "Shatter Sorcery",
    "sigil_of_solace.go": "Sigil of Solace",
    "talisman_of_balance.go": "Talisman of Balance",
    "talisman_of_cremation.go": "Talisman of Cremation",
    "talisman_of_dousing.go": "Talisman of Dousing",
    "talisman_of_featherfoot.go": "Talisman of Featherfoot",
    "talisman_of_recompense.go": "Talisman of Recompense",
    "talisman_of_tithes.go": "Talisman of Tithes",
    "talisman_of_warfare.go": "Talisman of Warfare",
    "thrust.go": "Thrust",
    "timesnap_potion.go": "Timesnap Potion",
    "titanium_bauble.go": "Titanium Bauble",
}


def main():
    gen_dir = os.path.join(ROOT, "internal", "card", "generic")
    weapon_dir = os.path.join(ROOT, "internal", "weapon")

    changed = 0
    for fname, card_name in CARD_FILES.items():
        path = os.path.join(gen_dir, fname)
        if process_file(path, card_name):
            changed += 1
            print("rewrote", path)

    talishar_path = os.path.join(weapon_dir, "talishar.go")
    if process_file(talishar_path, "Talishar, the Lost Prince"):
        changed += 1
        print("rewrote", talishar_path)

    print(f"changed {changed} files")


if __name__ == "__main__":
    main()
