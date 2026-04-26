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
    "Amulet of Assertiveness": "AR grant: target attack 'banish top of deck on hit'; gated on 4+ cards in hand",
    "Amulet of Echoes": "Instant 'opposing hero discards 2'; gated on a repeat-name play this turn",
    "Amulet of Havencall": "DR tutor for Rally the Rearguard; gated on empty hand",
    "Amulet of Ignition": "Instant 'next activated ability costs {r} less'",
    "Amulet of Intervention": "Instant 1 damage prevention; gated on a lethal incoming source",
    "Amulet of Oblation": "Instant 'graveyard → bottom of deck' replacement; gated on graveyard entry",
    "Arcane Polarity": "1{h} gain (4/3/2{h} if dealt arcane damage this turn)",
    "Blade Flash": "AR 'target sword attack gains go again'",
    "Brush Off": "Instant 'prevent next damage of N or less' (3 / 2 / 1 by colour)",
    "Calming Breeze": "Instant 'prevent 1 of each of the next 3 damage events'",
    "Clap 'Em in Irons": "passive tap-target Pirate; can't unfreeze; self-destroys at start of turn",
    "Clarity Potion": "activated Opt 2",
    "Count Your Blessings": "graveyard-scaled X{h} gain (also banlisted)",
    "Cracked Bauble": "draft-format pitch resource; no other effect",
    "Destructive Tendencies": "Instant 'remove counters from target item / aura token'",
    "Eirina's Prayer": "Instant prevent X arcane to your hero; X scaled by revealed top-card pitch",
    "Energy Potion": "activated 'gain {r}{r}'",
    "Even Bigger Than That!": "Opt + reveal-and-Quicken trigger; gated on damage dealt this turn",
    "Exposed": "AR +1{p}; gated on attacker not being marked",
    "Fool's Gold": "discard trigger creates a Gold token",
    "Healing Potion": "activated 2{h} gain",
    "Imperial Seal of Command": "activated 'no DR this turn' + Royal-only arsenal-wipe on hit",
    "Lunging Press": "AR +1{p} buff to a target attack action card",
    "Memorial Ground": "Instant 'graveyard → top of deck' for low-cost attack action card",
    "Nip at the Heels": "AR +1{p} buff to a target attack with ≤3 base {p}",
    "Oasis Respite": "Instant 'prevent N damage from chosen source to target hero'; conditional 1{h}",
    "Peace of Mind": "Instant 'prevent 4 of next {p}-damage hit'; creates a Ponder token",
    "Pilfer the Tomb": "Instant banish from an opposing graveyard / aura",
    "Potion of Déjà Vu": "activated 'put pitch zone on top of deck in any order'",
    "Potion of Ironhide": "activated +1{d} buff on all your attack actions this turn",
    "Potion of Luck": "activated 'shuffle hand+arsenal into deck, draw that many'",
    "Potion of Seeing": "activated reveal opposing hero's hand",
    "Potion of Strength": "activated +2{p} on next attack",
    "Pummel": "modal AR +4{p}: club/hammer weapon attack OR cost-2+ attack action (on-hit discard)",
    "Razor Reflex": "modal AR +N{p}: dagger/sword weapon attack OR cost ≤1 attack action (on-hit go again)",
    "Reinforce the Line": "Instant +N{d} grant to a defending attack action card",
    "Shatter Sorcery": "Instant: destroy a Sigil aura, and/or prevent 1 arcane damage",
    "Sigil of Solace": "3/2/1{h} gain (also banlisted)",
    "Talisman of Balance": "end-phase arsenal-fill from top of deck if behind on arsenal count",
    "Talisman of Cremation": "self-destroys on play-from-banished → banish a named card from opposing graveyards",
    "Talisman of Dousing": "passive Spellvoid 1",
    "Talisman of Featherfoot": "self-destroys when an attack gains exactly +1{p} in the reaction step → grants go again",
    "Talisman of Recompense": "self-destroys on pitching a 1-resource card → gain {r}{r}{r} instead",
    "Talisman of Tithes": "self-destroys on an opposing draw during your action phase → opponent draws minus 1",
    "Talisman of Warfare": "self-destroys + wipes all arsenals on a 2-damage hit",
    "Thrust": "AR +3{p} buff to a target sword attack",
    "Timesnap Potion": "activated 'gain 2 action points'",
    "Titanium Bauble": "pitch-3 resource with 3{d}; no other effect",
    "Talishar, the Lost Prince": "rust-counter activation cost and end-phase self-destruct at 3+ counters",
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

    # 2. Drop any existing `// not implemented: ...` line that sits directly above a
    #    NotImplemented() method, then insert a fresh one with the current note. This makes the
    #    rewrite idempotent and lets us correct previous notes by editing NOTES and re-running.
    new_src = re.sub(
        r"^[ \t]*// not implemented:[^\n]*\n(?=[ \t]*func \([A-Za-z0-9]+\) NotImplemented\(\))",
        "",
        new_src,
        flags=re.MULTILINE,
    )
    new_src = re.sub(
        r"(^[ \t]*)(func \([A-Za-z0-9]+\) NotImplemented\(\))",
        rf"\1// not implemented: {note}\n\1\2",
        new_src,
        flags=re.MULTILINE,
    )

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
