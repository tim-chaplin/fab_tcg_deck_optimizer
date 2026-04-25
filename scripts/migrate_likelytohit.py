"""Bulk rewrite LikelyToHit call sites for the new CardState-typed signature.

Patterns handled:
  card.LikelyToHit(attack, self.EffectiveDominate())             → card.LikelyToHit(self)
  card.LikelyToHit(c.Attack(), self.EffectiveDominate())          → card.LikelyToHit(self)
  card.LikelyToHit(target.Card.Attack(), target.EffectiveDominate()) → card.LikelyToHit(target)
  card.LikelyToHit(pc.Card.Attack(), pc.EffectiveDominate())      → card.LikelyToHit(pc)

Raw-integer probes (Fragile Aura's runechant check):
  card.LikelyToHit(<expr>, false)                                 → card.LikelyDamageHits(<expr>, false)

Nebula Blade's internal self-buff (`dmg += 3` for NonAttackActionPlayed) and the resulting
`LikelyToHit(dmg, …)` call need a manual edit — that case migrates to BonusAttack on self —
so we leave its file untouched here and edit it by hand.
"""

import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
TARGETS = [
    "internal/card/generic/blanch.go",
    "internal/card/generic/blow_for_a_blow.go",
    "internal/card/generic/crash_down_the_gates.go",
    "internal/card/generic/cut_down_to_size.go",
    "internal/card/generic/destructive_deliberation.go",
    "internal/card/generic/down_but_not_out.go",
    "internal/card/generic/fact_finding_mission.go",
    "internal/card/generic/hand_behind_the_pen.go",
    "internal/card/generic/humble.go",
    "internal/card/generic/jack_be_nimble.go",
    "internal/card/generic/jack_be_quick.go",
    "internal/card/generic/life_for_a_life.go",
    "internal/card/generic/money_or_your_life.go",
    "internal/card/generic/performance_bonus.go",
    "internal/card/generic/pursue_to_the_edge_of_oblivion.go",
    "internal/card/generic/pursue_to_the_pits_of_despair.go",
    "internal/card/generic/smash_up.go",
    "internal/card/generic/snatch.go",
    "internal/card/generic/strike_gold.go",
    "internal/card/generic/tongue_tied.go",
    "internal/card/generic/walk_the_plank.go",
    "internal/card/generic/wreck_havoc.go",
    "internal/card/runeblade/consuming_volition.go",
    "internal/card/runeblade/mauvrion_skies.go",
    "internal/card/runeblade/meat_and_greet.go",
    "internal/card/runeblade/reek_of_corruption.go",
    "internal/card/runeblade/runic_reaping.go",
]

# (regex, replacement) pairs applied in order. Each captures any leading whitespace via the
# match itself — we don't need to preserve indentation explicitly because the replacement
# leaves the column alone.
REWRITES = [
    # Most common pattern: helper-passed attack int + self.EffectiveDominate.
    (r"card\.LikelyToHit\(attack, self\.EffectiveDominate\(\)\)",
     "card.LikelyToHit(self)"),
    # Inlined c.Attack() variant.
    (r"card\.LikelyToHit\(c\.Attack\(\), self\.EffectiveDominate\(\)\)",
     "card.LikelyToHit(self)"),
    # Mauvrion / Runic Reaping handler shape: target is a *CardState.
    (r"card\.LikelyToHit\(target\.Card\.Attack\(\), target\.EffectiveDominate\(\)\)",
     "card.LikelyToHit(target)"),
    # Fragile Aura's pc-typed peek.
    (r"card\.LikelyToHit\(pc\.Card\.Attack\(\), pc\.EffectiveDominate\(\)\)",
     "card.LikelyToHit(pc)"),
    # Fragile Aura's raw runechant probe — no CardState; route to LikelyDamageHits.
    (r"card\.LikelyToHit\(runechants, false\)",
     "card.LikelyDamageHits(runechants, false)"),
]


def main():
    changed = 0
    for rel in TARGETS:
        path = os.path.join(ROOT, rel)
        with open(path, encoding="utf-8") as f:
            src = f.read()
        new = src
        for pat, repl in REWRITES:
            new = re.sub(pat, repl, new)
        if new != src:
            with open(path, "w", encoding="utf-8", newline="\n") as f:
                f.write(new)
            print("rewrote", path)
            changed += 1

    # Fragile Aura is in TARGETS implicitly via its file path — add it.
    fa = os.path.join(ROOT, "internal/card/runeblade/fragile_aura.go")
    with open(fa, encoding="utf-8") as f:
        src = f.read()
    new = src
    for pat, repl in REWRITES:
        new = re.sub(pat, repl, new)
    if new != src:
        with open(fa, "w", encoding="utf-8", newline="\n") as f:
            f.write(new)
        print("rewrote", fa)
        changed += 1

    print(f"changed {changed} files")


if __name__ == "__main__":
    main()
