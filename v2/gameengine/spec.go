package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// Spec is the test-friendly engine constructor input. Mirrors the scalar / card-typed
// fields of GameState in exported form so callers outside the package don't have to
// chain per-field setters after construction. Build via NewFromSpec.
//
// Spec carries no Auras / Triggers / Items list: callers seed those individually after
// NewFromSpec via CreateAura / CreateTrigger / CreateItem on the underlying GameState.
type Spec struct {
	Hero                  Hero
	Arsenal               card.Card
	Banished              []card.Card
	Graveyard             []card.Card
	Pitched               []card.Card
	Defenders             []card.Card
	CardsPlayed           []card.Card
	CardsRemaining        []*card.CardState
	OpponentMarked        bool
	ArcaneDamageDealt     bool
	AuraCreated           bool
	CardBanished          bool
	NonAttackActionPlayed bool
	ActionPoints          int
	IncomingDamage        int
	ArcaneIncomingDamage  int
	BlockTotal            int
	Value                 int
	TriggeringCard        card.Card
	AttackReactionTarget  *card.CardState
}

// New returns a bare *GameEngine wrapping a fresh GameState. Sim's chain runner
// builds higher-level constructs (master engines, per-permutation states) on top of
// this; tests typically use NewFromSpec or NewFromCards instead.
func New() *GameEngine {
	return &GameEngine{state: NewState()}
}

// NewState returns a bare *GameState — cacheable=true, currentAuraIdx=-1, fresh logger.
// Use directly when only the state container is needed (e.g. by sim's chain runner that
// holds *GameState and wraps in *GameEngine when calling card hooks).
func NewState() *GameState {
	return &GameState{
		cacheable:      true,
		currentAuraIdx: -1,
		logger:         turnlogger.New(),
	}
}

// NewFromSpec builds a *GameEngine from a Spec, sealing the unexported fields inside
// the package. Use when an external caller — turntest, hero test, EvalOneTurnForTesting
// — needs to construct a prior-turn state. Auras / Triggers / Items aren't on Spec;
// callers chain CreateAura / CreateTrigger / CreateItem on the state after NewFromSpec.
func NewFromSpec(spec Spec) *GameEngine {
	s := &GameState{
		hero:                  spec.Hero,
		arsenal:               spec.Arsenal,
		banished:              spec.Banished,
		graveyard:             spec.Graveyard,
		pitched:               spec.Pitched,
		defenders:             spec.Defenders,
		cardsPlayed:           spec.CardsPlayed,
		cardsRemaining:        spec.CardsRemaining,
		opponentMarked:        spec.OpponentMarked,
		arcaneDamageDealt:     spec.ArcaneDamageDealt,
		auraCreated:           spec.AuraCreated,
		cardBanished:          spec.CardBanished,
		nonAttackActionPlayed: spec.NonAttackActionPlayed,
		actionPoints:          spec.ActionPoints,
		incomingDamage:        spec.IncomingDamage,
		arcaneIncomingDamage:  spec.ArcaneIncomingDamage,
		blockTotal:            spec.BlockTotal,
		value:                 spec.Value,
		triggeringCard:        spec.TriggeringCard,
		attackReactionTarget:  spec.AttackReactionTarget,
		cacheable:             true,
		currentAuraIdx:        -1,
		logger:                turnlogger.New(),
	}
	return &GameEngine{state: s}
}

// NewFromCards is a test-only constructor that wraps a Card slice in a fresh *deck.Deck
// and seeds the graveyard. Lets card tests build a GameEngine with a hand-rolled deck
// order without each test plumbing the deck construction inline.
func NewFromCards(deckCards, graveyard []card.Card) *GameEngine {
	dc := make([]deck.Card, len(deckCards))
	for i, c := range deckCards {
		dc[i] = c
	}
	g := New()
	g.state.deck = deck.New(nil, nil, dc)
	g.state.graveyard = graveyard
	return g
}
