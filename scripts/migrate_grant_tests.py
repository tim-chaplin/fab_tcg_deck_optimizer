"""Migrate per-card unit tests for "next attack +N{p}" granters from the prior contract
(granter's Play returned the bonus) to the new BonusAttack contract (granter writes
BonusAttack on the target's CardState and returns 0).

Each affected test follows the shape:
  s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stub...}}}
  for _, tc := range cases {
      if got := tc.c.Play(&s, &card.CardState{}); got != tc.want { ... }
  }

The rewrite:
  - Creates a fresh `s` per iteration (so BonusAttack doesn't accumulate across cases).
  - Asserts the granter's Play returns 0 (its own contribution).
  - Asserts the first CardState in CardsRemaining picked up the expected BonusAttack.
"""

import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# Test files where the "ReturnsBonus" / "GrantsBonus" pattern lives.
TARGETS = [
    "internal/card/generic/clearwater_elixir_test.go",
    "internal/card/generic/come_to_fight_test.go",
    "internal/card/generic/force_sight_test.go",
    "internal/card/generic/minnowism_test.go",
    "internal/card/generic/money_where_ya_mouth_is_test.go",
    "internal/card/generic/nimblism_test.go",
    "internal/card/generic/plunder_run_test.go",
    "internal/card/generic/prime_the_crowd_test.go",
    "internal/card/generic/public_bounty_test.go",
    "internal/card/generic/regain_composure_test.go",
    "internal/card/generic/restvine_elixir_test.go",
    "internal/card/generic/sapwood_elixir_test.go",
    "internal/card/generic/sloggism_test.go",
    "internal/card/generic/smashing_good_time_test.go",
    "internal/card/generic/warmongers_recital_test.go",
]

# The shared "ReturnsBonus" test pattern. Captures the stub call so each rewritten test
# uses the same one. Multi-line.
PATTERN = re.compile(
    r"^(?P<indent>[ \t]*)s := card\.TurnState\{CardsRemaining: \[\]\*card\.CardState\{\{Card: (?P<stub>stub[A-Za-z0-9_]+\([^)]*\))\}\}\}\n"
    r"(?P<indent2>[ \t]*)cases := \[\]struct \{\n"
    r"(?P<fields>(?:[ \t]*[^\n]+\n)+?)"
    r"(?P<indent3>[ \t]*)\}\{\n"
    r"(?P<rows>(?:[ \t]*\{[^\n]+\}[,\n]*\n)+?)"
    r"(?P<indent4>[ \t]*)\}\n"
    r"(?P<loopindent>[ \t]*)for _, tc := range cases \{\n"
    r"(?P<bodyindent>[ \t]*)if got := tc\.c\.Play\(&s, &card\.CardState\{\}\); got != tc\.want \{\n"
    r"(?P<errindent>[ \t]*)t\.Errorf\((?P<errmsg>\"[^\"]*\"), tc\.c\.Name\(\), got, tc\.want\)\n"
    r"(?P<endif>[ \t]*)\}\n"
    r"(?P<endloop>[ \t]*)\}\n",
    re.MULTILINE,
)


def rewrite(src):
    def repl(m):
        i = m.group("indent")
        i2 = m.group("indent2")
        i3 = m.group("indent3")
        i4 = m.group("indent4")
        li = m.group("loopindent")
        bi = m.group("bodyindent")
        ei = m.group("errindent")
        endif = m.group("endif")
        endloop = m.group("endloop")
        stub = m.group("stub")
        fields = m.group("fields")
        rows = m.group("rows")
        new = (
            f"{i2}cases := []struct {{\n"
            f"{fields}"
            f"{i3}}}{{\n"
            f"{rows}"
            f"{i4}}}\n"
            f"{li}for _, tc := range cases {{\n"
            f"{bi}target := &card.CardState{{Card: {stub}}}\n"
            f"{bi}s := card.TurnState{{CardsRemaining: []*card.CardState{{target}}}}\n"
            f"{bi}if got := tc.c.Play(&s, &card.CardState{{}}); got != 0 {{\n"
            f"{ei}t.Errorf(\"%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)\", tc.c.Name(), got)\n"
            f"{endif}}}\n"
            f"{bi}if target.BonusAttack != tc.want {{\n"
            f"{ei}t.Errorf(\"%s: target BonusAttack = %d, want %d\", tc.c.Name(), target.BonusAttack, tc.want)\n"
            f"{endif}}}\n"
            f"{endloop}}}\n"
        )
        return new

    return PATTERN.sub(repl, src)


def main():
    changed = 0
    for rel in TARGETS:
        path = os.path.join(ROOT, rel)
        with open(path, encoding="utf-8") as f:
            src = f.read()
        new = rewrite(src)
        if new != src:
            with open(path, "w", encoding="utf-8", newline="\n") as f:
                f.write(new)
            print("rewrote", rel)
            changed += 1
        else:
            print("(no change)", rel, file=sys.stderr)
    print(f"changed {changed} files")


if __name__ == "__main__":
    main()
