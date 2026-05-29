package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// StateBuilder fluently constructs a *GameState. Callers chain setters and finish
// with Build() to retrieve the *GameState, then wrap in a *GameEngine via
// &gameengine.GameEngine{GameState: ...} when they need to drive Card.Play hooks
// through the rules engine.
//
// The builder starts from the per-turn-locals defaults a just-constructed *GameState
// carries (cacheable=true, currentHookIdx=-1, empty logger) so callers only set the
// fields they care about.
type StateBuilder struct {
	gs *GameState
}

// New returns a fresh *GameEngine wrapping an empty default *GameState. Shorthand for
// &GameEngine{GameState: GameStateBuilder().Build()}; reach for the builder directly
// when you need to populate fields.
func New() *GameEngine {
	return &GameEngine{GameState: GameStateBuilder().Build()}
}

// GameStateBuilder returns a fresh *StateBuilder pre-populated with the defaults a
// just-constructed *GameState carries: cacheable=true, currentHookIdx=-1, a no-op logger
// (NoopLogger), a non-nil empty *deck.Deck (so cards that recycle into the deck don't
// panic on a zero-card state), and a no-ability fallback hero (20 health, 4 Intelligence).
// SetDeck / SetCards overrides the empty deck; SetHero overrides the default hero.
func GameStateBuilder() *StateBuilder {
	gs := &GameState{
		ephemeral: ephemeral{
			cacheable:              true,
			currentHookIdx:         -1,
			currentFiringTokenAura: -1,
			currentFiringTokenItem: -1,
			logger:                 NoopLogger{},
		},
		deck: deck.New(nil, nil, nil),
		hero: defaultHero{},
	}
	gs.initTokenSlots()
	return &StateBuilder{gs: gs}
}

// Build returns the configured *GameState.
func (b *StateBuilder) Build() *GameState { return b.gs }

// SetHero installs h as the active hero. Routes through GameState.SetHero so the cached
// heroTriggerType stays in sync.
func (b *StateBuilder) SetHero(h Hero) *StateBuilder { b.gs.SetHero(h); return b }

// SetArsenal installs c into the arsenal slot.
func (b *StateBuilder) SetArsenal(c card.Card) *StateBuilder { b.gs.arsenal = c; return b }

// AddAura appends auras to the carryover aura list. Unlike GameState.CreateAura it does
// not flip the auraCreated flag — builder auras are pre-turn carryover, not auras created
// during this turn's attack phase; reach for SetAuraCreated to set that flag explicitly.
// Token-aura entries (Runechant, Ponder) route into the pre-allocated tokenAuras slot
// for their kind (count merged via SetCount) instead of appending to the card-aura list.
func (b *StateBuilder) AddAura(auras ...Aura) *StateBuilder {
	for _, a := range auras {
		if kind, ok := tokenAuraKindForID(a.CardID()); ok {
			slot := b.gs.tokenAuras[kind]
			slot.SetCount(slot.Count() + a.Count())
			if slot.Count() > 0 {
				b.gs.tokenAurasLiveBits |= 1 << kind
			}
			continue
		}
		b.gs.auras = append(b.gs.auras, a)
	}
	return b
}

// CreateAuraFromCard plays c against an ephemeral *GameEngine wrapping the state being
// built so the card's Play body lands its aura on b.gs without test fixtures needing to
// duplicate the card's aura-construction logic. Best for pure-aura cards (Sigil of Fyendal,
// Runeblood Incantation, Blessing of Occult, etc.); cards whose Play has additional side
// effects (graveyard banishes, on-hit registration, immediate value credit) run those too.
func (b *StateBuilder) CreateAuraFromCard(c card.Card) *StateBuilder {
	b.playForCarryover(c)
	return b
}

// CreateItemFromCard is the item counterpart of CreateAuraFromCard: it plays c so the
// card's Play lands its in-play item on b.gs. Used to seed a card-sourced item (e.g.
// Talisman of Recompense) as start-of-turn carryover.
func (b *StateBuilder) CreateItemFromCard(c card.Card) *StateBuilder {
	b.playForCarryover(c)
	return b
}

// playForCarryover runs c.Play against an ephemeral *GameEngine wrapping b.gs so the
// card's own Play body lands its carryover state (aura, item) on the state being built.
func (b *StateBuilder) playForCarryover(c card.Card) {
	ge := &GameEngine{GameState: b.gs}
	c.Play(ge, NoopLogger{}, &card.CardState{Card: c})
}

// AddItem appends items to the carryover item list. Token-item entries (Gold, Silver,
// Copper) route into the pre-allocated tokenItems slot for their kind (count merged via
// SetCount) instead of appending to the card-item list.
func (b *StateBuilder) AddItem(items ...Item) *StateBuilder {
	for _, it := range items {
		if kind, ok := tokenItemKindForID(it.CardID()); ok {
			slot := b.gs.tokenItems[kind]
			slot.SetCount(slot.Count() + it.Count())
			if slot.Count() > 0 {
				b.gs.tokenItemsLiveBits |= 1 << kind
			}
			continue
		}
		b.gs.items = append(b.gs.items, it)
	}
	return b
}

// tokenAuraKindForID maps a token aura's CardID to its tokenAuraKind index, or returns
// false if id is not a recognised aura-token kind.
func tokenAuraKindForID(id ids.CardID) (tokenAuraKind, bool) {
	switch id {
	case ids.RunechantTokenID:
		return tokenAuraRunechant, true
	case ids.PonderTokenID:
		return tokenAuraPonder, true
	}
	return 0, false
}

// tokenItemKindForID maps a token item's CardID to its tokenItemKind index, or returns
// false if id is not a recognised item-token kind.
func tokenItemKindForID(id ids.CardID) (tokenItemKind, bool) {
	switch id {
	case ids.GoldTokenID:
		return tokenItemGold, true
	case ids.SilverTokenID:
		return tokenItemSilver, true
	case ids.CopperTokenID:
		return tokenItemCopper, true
	}
	return 0, false
}

// SetBanished replaces the banished-zone slice.
func (b *StateBuilder) SetBanished(cs []card.Card) *StateBuilder { b.gs.banished = cs; return b }

// SetGraveyard replaces the graveyard slice.
func (b *StateBuilder) SetGraveyard(cs []card.Card) *StateBuilder { b.gs.graveyard = cs; return b }

// SetPitched replaces the pitched-this-turn slice. Each entry is the physical
// pitched-zone copy (carries in-attack-turn state).
func (b *StateBuilder) SetPitched(cs []*card.CardState) *StateBuilder { b.gs.pitched = cs; return b }

// SetDefenders replaces the defenders slice.
func (b *StateBuilder) SetDefenders(cs []card.Card) *StateBuilder { b.gs.defenders = cs; return b }

// SetCardsPlayed replaces the cards-played-this-turn slice.
func (b *StateBuilder) SetCardsPlayed(cs []card.Card) *StateBuilder { b.gs.cardsPlayed = cs; return b }

// SetCardsRemaining replaces the cards-scheduled-after-current-attack-step slice.
func (b *StateBuilder) SetCardsRemaining(cs []*card.CardState) *StateBuilder {
	b.gs.cardsRemaining = cs
	return b
}

// SetOpponentMarked sets the Mark token state on the opposing hero.
func (b *StateBuilder) SetOpponentMarked(v bool) *StateBuilder { b.gs.opponentMarked = v; return b }

// SetArcaneDamageDealt overrides the "arcane damage has been dealt this turn" sticky
// flag.
func (b *StateBuilder) SetArcaneDamageDealt(v bool) *StateBuilder {
	b.gs.arcaneDamageDealt = v
	return b
}

// SetAuraCreated overrides the "an aura was created this turn" sticky flag.
func (b *StateBuilder) SetAuraCreated(v bool) *StateBuilder { b.gs.auraCreated = v; return b }

// SetCardBanished overrides the "a card was banished this turn" sticky flag.
func (b *StateBuilder) SetCardBanished(v bool) *StateBuilder { b.gs.cardBanished = v; return b }

// SetNonAttackActionPlayed overrides the "a non-attack action has resolved" sticky
// flag.
func (b *StateBuilder) SetNonAttackActionPlayed(v bool) *StateBuilder {
	b.gs.nonAttackActionPlayed = v
	return b
}

// SetIncomingPhysicalDamage replaces the matchup's incoming-damage tally.
func (b *StateBuilder) SetIncomingPhysicalDamage(n int) *StateBuilder {
	b.gs.incomingPhysicalDamage = n
	return b
}

// SetIncomingArcaneDamage replaces the matchup's arcane-incoming-damage tally.
func (b *StateBuilder) SetIncomingArcaneDamage(n int) *StateBuilder {
	b.gs.incomingArcaneDamage = n
	return b
}

// SetAttackReactionTarget installs the buff target for the AR resolving next.
func (b *StateBuilder) SetAttackReactionTarget(cs *card.CardState) *StateBuilder {
	b.gs.attackReactionTarget = cs
	return b
}

// SetDeck replaces the attack-turn runner deck.
func (b *StateBuilder) SetDeck(d *deck.Deck) *StateBuilder { b.gs.deck = d; return b }

// SetCards wraps the supplied cards in a fresh *deck.Deck and assigns it as the
// state's deck. Equivalent to SetDeck(deck.New(nil, nil, ...cards)) — convenient for
// tests that want to pin a deterministic deck order.
func (b *StateBuilder) SetCards(cs []card.Card) *StateBuilder {
	dc := make([]deck.Card, len(cs))
	for i, c := range cs {
		dc[i] = c
	}
	b.gs.deck = deck.New(nil, nil, dc)
	return b
}
