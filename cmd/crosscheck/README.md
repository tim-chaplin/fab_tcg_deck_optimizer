# cmd/crosscheck

A diagnostic that compares our deck-construction pool against the **official** Flesh and Blood
card database (cardvault.fabtcg.com), so we can spot cards the official site says are legal
that we're missing — or cards in our pool the official site no longer lists.

It pulls every Silver-Age-legal Runeblade and Generic print from the official advanced-search
API, keeps the no-talent main-deck cards (what Viserai can legally run), and diffs that set
against `registry.LegalCardsFor(SilverAge, Viserai)`. Each difference is annotated:

- **[A] Official-legal but not in our pool** — split into cards we never implemented (no yaml)
  vs. cards we implemented but quarantined in the `unplayable/` / `notimplemented/`
  subpackages (which the registry doesn't compile in, so `registry.AllCards` can't see them).
- **[B] In our pool but not official-legal** — cards we'd let Viserai run that the official
  query doesn't return (a recent banlist change, a non-deck/talented typebox, or a name
  mismatch). Empty is the healthy state.

## Run

```
go run ./cmd/crosscheck            # needs outbound HTTPS to the cardvault API
go run ./cmd/crosscheck -cards-dir internal/card/cards   # override the yaml root
```

## How the typebox filter works

A printed typebox reads `<supertypes> <CardType> - <subtypes>` (double-faced cards join faces
with `||`). `inScope` reduces the grammar to a word blacklist — a card is out of scope if
any word across either face is a talent (Lightning, Shadow, Earth, …), an off-class (Warrior,
Wizard, …), or a non-deck card-type (Weapon, Equipment, Hero, Token, …). The blacklist is
order-independent, so multi-word types like "Defense Reaction" don't trip a positional parse.
Keep the talent / class lists in sync with the official type line as new sets land.
