package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// StateBuilder fluently constructs a *GameState. Callers chain setters and finish
// with Build() to retrieve the *GameState, then wrap in a *GameEngine via
// &gameengine.GameEngine{GameState: ...} when they need to drive Card.Play hooks
// through the rules engine.
//
// The builder starts from the chain-locals defaults a just-constructed *GameState
// carries (cacheable=true, currentAuraIdx=-1, empty logger) so callers only set the
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
// just-constructed *GameState carries: cacheable=true, currentAuraIdx=-1, empty
// logger, and a non-nil empty *deck.Deck (so cards that recycle into the deck don't
// panic on a zero-card state). SetDeck / SetCards overrides the empty deck.
func GameStateBuilder() *StateBuilder {
	return &StateBuilder{
		gs: &GameState{
			cacheable:      true,
			currentAuraIdx: -1,
			logger:         turnlogger.New(),
			deck:           deck.New(nil, nil, nil),
		},
	}
}

// Build returns the configured *GameState.
func (b *StateBuilder) Build() *GameState { return b.gs }

// SetHero installs h as the active hero.
func (b *StateBuilder) SetHero(h Hero) *StateBuilder { b.gs.hero = h; return b }

// SetArsenal installs c into the arsenal slot.
func (b *StateBuilder) SetArsenal(c card.Card) *StateBuilder { b.gs.arsenal = c; return b }

// AddAura appends auras to the carryover aura list. Unlike GameState.CreateAura it does
// not flip the auraCreated flag — builder auras are pre-turn carryover, not auras created
// during this turn's chain; reach for SetAuraCreated to set that flag explicitly.
func (b *StateBuilder) AddAura(auras ...Aura) *StateBuilder {
	b.gs.auras = append(b.gs.auras, auras...)
	return b
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

// SetActionPoints replaces the running action-point pool.
func (b *StateBuilder) SetActionPoints(n int) *StateBuilder { b.gs.actionPoints = n; return b }

// SetIncomingDamage replaces the matchup's incoming-damage tally.
func (b *StateBuilder) SetIncomingDamage(n int) *StateBuilder { b.gs.incomingDamage = n; return b }

// SetArcaneIncomingDamage replaces the matchup's arcane-incoming-damage tally.
func (b *StateBuilder) SetArcaneIncomingDamage(n int) *StateBuilder {
	b.gs.arcaneIncomingDamage = n
	return b
}

// SetBlockTotal replaces the partition's uncapped defense sum.
func (b *StateBuilder) SetBlockTotal(n int) *StateBuilder { b.gs.blockTotal = n; return b }

// SetValue replaces the running chain value.
func (b *StateBuilder) SetValue(n int) *StateBuilder { b.gs.value = n; return b }

// SetTriggeringCard replaces the triggering-card slot.
func (b *StateBuilder) SetTriggeringCard(c card.Card) *StateBuilder {
	b.gs.triggeringCard = c
	return b
}

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
