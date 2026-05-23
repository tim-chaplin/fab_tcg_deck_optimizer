package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// StateBuilder fluently constructs a *GameState. Callers chain setters and finish
// with Build() to retrieve the *GameState, then wrap in a *GameEngine via
// &gameengine.GameEngine{GameState: ...} when they need to drive Card.Play hooks
// through the rules engine.
//
// The builder starts from the chain-locals defaults a just-constructed *GameState
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
	return &StateBuilder{
		gs: &GameState{
			cacheable:      true,
			currentHookIdx: -1,
			logger:         NoopLogger{},
			deck:           deck.New(nil, nil, nil),
			hero:           defaultHero{},
		},
	}
}

// Build returns the configured *GameState.
func (b *StateBuilder) Build() *GameState { return b.gs }

// SetHero installs h as the active hero.
func (b *StateBuilder) SetHero(h Hero) *StateBuilder { b.gs.hero = h; return b }

// SetWeapons installs the equipped weapons. Persistent across turns; mid-game equip /
// unequip card abilities mutate via GameState.SetWeapons.
func (b *StateBuilder) SetWeapons(w []weapon.Weapon) *StateBuilder { b.gs.weapons = w; return b }

// SetArsenal installs c into the arsenal slot.
func (b *StateBuilder) SetArsenal(c card.Card) *StateBuilder { b.gs.arsenal = c; return b }

// AddAura appends auras to the carryover aura list. Unlike GameState.CreateAura it does
// not flip the auraCreated flag — builder auras are pre-turn carryover, not auras created
// during this turn's chain; reach for SetAuraCreated to set that flag explicitly.
func (b *StateBuilder) AddAura(auras ...Aura) *StateBuilder {
	b.gs.auras = append(b.gs.auras, auras...)
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

// AddItem appends items to the carryover item list.
func (b *StateBuilder) AddItem(items ...Item) *StateBuilder {
	b.gs.items = append(b.gs.items, items...)
	return b
}

// SetBanished replaces the banished-zone slice.
func (b *StateBuilder) SetBanished(cs []card.Card) *StateBuilder { b.gs.banished = cs; return b }

// SetGraveyard replaces the graveyard slice.
func (b *StateBuilder) SetGraveyard(cs []card.Card) *StateBuilder { b.gs.graveyard = cs; return b }

// SetPitched replaces the pitched-this-turn slice.
func (b *StateBuilder) SetPitched(cs []card.Card) *StateBuilder { b.gs.pitched = cs; return b }

// SetDefenders replaces the defenders slice.
func (b *StateBuilder) SetDefenders(cs []card.Card) *StateBuilder { b.gs.defenders = cs; return b }

// SetCardsPlayed replaces the cards-played-this-turn slice.
func (b *StateBuilder) SetCardsPlayed(cs []card.Card) *StateBuilder { b.gs.cardsPlayed = cs; return b }

// SetCardsRemaining replaces the cards-scheduled-after-current-chain-step slice.
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

// SetIncomingDamage replaces the matchup's incoming-damage tally.
func (b *StateBuilder) SetIncomingDamage(n int) *StateBuilder { b.gs.incomingDamage = n; return b }

// SetArcaneIncomingDamage replaces the matchup's arcane-incoming-damage tally.
func (b *StateBuilder) SetArcaneIncomingDamage(n int) *StateBuilder {
	b.gs.arcaneIncomingDamage = n
	return b
}

// SetBlockTotal replaces the partition's uncapped defense sum.
func (b *StateBuilder) SetBlockTotal(n int) *StateBuilder { b.gs.blockTotal = n; return b }

// SetAttackReactionTarget installs the buff target for the AR resolving next.
func (b *StateBuilder) SetAttackReactionTarget(cs *card.CardState) *StateBuilder {
	b.gs.attackReactionTarget = cs
	return b
}

// SetDeck replaces the chain-runner deck.
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
