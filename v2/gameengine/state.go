package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// GameState owns the raw per-turn data — every slice, every scalar, every flag the
// engine reads or writes during a chain run. Internal machinery (sim's per-permutation
// scratch, snapshots that carry into the next turn, the find-best winning pointer) holds
// *GameState because it cares about the data, not the rules. GameEngine wraps a
// *GameState and adds the rules-engine API cards see; the unexported fields stay
// package-private so external callers can only touch them through methods.
type GameState struct {
	hero Hero

	hand      []card.Card
	deck      *deck.Deck // owned scratch; refilled on Reset via CopyFrom
	arsenal   card.Card
	graveyard []card.Card
	banished  []card.Card
	auras     []Aura
	triggers  []Trigger
	items     []Item

	cardsPlayed    []card.Card
	cardsRemaining []*card.CardState
	pitched        []card.Card
	defenders      []card.Card

	logger               *turnlogger.TurnLogger
	triggeringCard       card.Card
	attackReactionTarget *card.CardState

	actionPoints         int
	value                int
	cardsDrawn           int
	incomingDamage       int
	arcaneIncomingDamage int
	blockTotal           int
	currentAuraIdx       int

	cardBanished          bool
	arcaneDamageDealt     bool
	opponentMarked        bool
	auraCreated           bool
	nonAttackActionPlayed bool
	currentAuraDestroyed  bool
	currentStepRerouted   bool
	cacheable             bool
}

// Engine wraps s in a *GameEngine so the chain runner can drive Card.Play hooks against
// the rules-engine API while the underlying state remains the same pointer the caller
// holds. Cheap (single struct allocation); the engine doesn't copy state.
func (s *GameState) Engine() *GameEngine { return &GameEngine{state: s} }

// Copy returns a deep copy of s. Slice and *deck.Deck fields get fresh backing storage;
// Aura / Item entries are deep-copied via their Copy() methods so per-permutation
// Count / FiredThisTurn mutations stay isolated. Triggers are effectively immutable after
// construction, so only the slice header is duplicated. Logger is reset to nil — the
// caller installs a fresh per-clone logger when recording, leaving find-best copies
// allocation-free.
func (s *GameState) Copy() *GameState {
	out := *s
	out.hand = appendCopy(nil, s.hand)
	if s.deck != nil {
		out.deck = s.deck.Copy()
	}
	out.graveyard = appendCopy(nil, s.graveyard)
	out.banished = appendCopy(nil, s.banished)
	out.pitched = appendCopy(nil, s.pitched)
	out.defenders = appendCopy(nil, s.defenders)
	out.cardsPlayed = appendCopy(nil, s.cardsPlayed)
	if len(s.cardsRemaining) > 0 {
		out.cardsRemaining = append([]*card.CardState(nil), s.cardsRemaining...)
	} else {
		out.cardsRemaining = nil
	}
	if len(s.auras) > 0 {
		out.auras = make([]Aura, len(s.auras))
		for i, a := range s.auras {
			out.auras[i] = a.Copy()
		}
	} else {
		out.auras = nil
	}
	if len(s.triggers) > 0 {
		out.triggers = append([]Trigger(nil), s.triggers...)
	} else {
		out.triggers = nil
	}
	if len(s.items) > 0 {
		out.items = make([]Item, len(s.items))
		for i, it := range s.items {
			out.items[i] = it.Copy()
		}
	} else {
		out.items = nil
	}
	out.logger = nil
	return &out
}

// BeginPermutation resets per-chain locals so the state is ready to play a fresh chain
// from this permutation's hand order. Auras, items, banished, graveyard, deck, arsenal,
// pitched, hero, and OpponentMarked carry over untouched — they represent the leaf's
// pre-chain state. logger is installed verbatim (pass nil for the find-best path, a
// fresh logger for the recording path).
func (s *GameState) BeginPermutation(hand []card.Card, incomingDamage int, logger *turnlogger.TurnLogger) {
	s.hand = hand
	s.cardsPlayed = nil
	s.cardsRemaining = nil
	s.triggers = nil
	s.triggeringCard = nil
	s.attackReactionTarget = nil
	s.actionPoints = 1
	s.value = 0
	s.cardsDrawn = 0
	s.incomingDamage = incomingDamage
	s.cardBanished = false
	s.arcaneDamageDealt = false
	s.nonAttackActionPlayed = false
	s.currentAuraDestroyed = false
	s.currentStepRerouted = false
	s.currentAuraIdx = -1
	s.cacheable = true
	s.logger = logger
	s.auraCreated = len(s.auras) > 0
}

// === Pure state accessors. No cacheable flips; sim uses these to drive the chain runner. ===

func (s *GameState) Hero() Hero          { return s.hero }
func (s *GameState) SetHero(h Hero)      { s.hero = h }
func (s *GameState) IsCacheable() bool   { return s.cacheable }
func (s *GameState) SetCacheable(v bool) { s.cacheable = v }

func (s *GameState) Hand() []card.Card     { return s.hand }
func (s *GameState) SetHand(h []card.Card) { s.hand = h }
func (s *GameState) AppendHandRaw(c card.Card) {
	s.hand = append(s.hand, c)
}

// RemoveFromHand removes the first matching card from the hand without flipping
// IsCacheable. Returns true if a card was removed.
func (s *GameState) RemoveFromHand(c card.Card) bool {
	for i := range s.hand {
		if s.hand[i] == c {
			s.hand = append(s.hand[:i], s.hand[i+1:]...)
			return true
		}
	}
	return false
}

func (s *GameState) Deck() *deck.Deck     { return s.deck }
func (s *GameState) SetDeck(d *deck.Deck) { s.deck = d }

func (s *GameState) Graveyard() []card.Card      { return s.graveyard }
func (s *GameState) SetGraveyard(gv []card.Card) { s.graveyard = gv }
func (s *GameState) AppendGraveyard(c card.Card) { s.graveyard = append(s.graveyard, c) }

func (s *GameState) Arsenal() card.Card     { return s.arsenal }
func (s *GameState) SetArsenal(c card.Card) { s.arsenal = c }

func (s *GameState) Banished() []card.Card     { return s.banished }
func (s *GameState) SetBanished(b []card.Card) { s.banished = b }

func (s *GameState) Auras() []Aura       { return s.auras }
func (s *GameState) ClearAuras()         { s.auras = nil }
func (s *GameState) Triggers() []Trigger { return s.triggers }
func (s *GameState) ClearTriggers()      { s.triggers = nil }
func (s *GameState) Items() []Item       { return s.items }
func (s *GameState) ClearItems()         { s.items = nil }

// CreateAura appends a to the aura list. Flips AuraCreated so same-turn "if you've
// played or created an aura" riders see the entry.
func (s *GameState) CreateAura(a Aura) {
	s.auras = append(s.auras, a)
	s.auraCreated = true
}

// CreateTrigger appends t to the one-shot trigger queue.
func (s *GameState) CreateTrigger(t Trigger) { s.triggers = append(s.triggers, t) }

// CreateItem appends i to the item list.
func (s *GameState) CreateItem(i Item) { s.items = append(s.items, i) }

// AuraCount returns the count of live auras. Used by gates like Yinti Yanti's "while you
// control an aura" rider.
func (s *GameState) AuraCount() int { return len(s.auras) }

// RunechantCount / PonderCount / GoldCount / SilverCount / CopperCount return the live
// token-aura / token-item count by display name. Both GameState and GameEngine expose
// these so end-of-turn callers (TurnSummary readers) can read counts off the state
// pointer directly without needing an engine wrapper.
func (s *GameState) RunechantCount() int { return auraCountByName(s.auras, tokenNameRunechant) }
func (s *GameState) PonderCount() int    { return auraCountByName(s.auras, tokenNamePonder) }
func (s *GameState) GoldCount() int      { return itemCountByName(s.items, tokenNameGold) }
func (s *GameState) SilverCount() int    { return itemCountByName(s.items, tokenNameSilver) }
func (s *GameState) CopperCount() int    { return itemCountByName(s.items, tokenNameCopper) }

func (s *GameState) Pitched() []card.Card     { return s.pitched }
func (s *GameState) SetPitched(p []card.Card) { s.pitched = p }

func (s *GameState) Defenders() []card.Card     { return s.defenders }
func (s *GameState) SetDefenders(d []card.Card) { s.defenders = d }

func (s *GameState) CardsPlayed() []card.Card               { return s.cardsPlayed }
func (s *GameState) SetCardsPlayed(cs []card.Card)          { s.cardsPlayed = cs }
func (s *GameState) CardsRemaining() []*card.CardState      { return s.cardsRemaining }
func (s *GameState) SetCardsRemaining(cs []*card.CardState) { s.cardsRemaining = cs }

func (s *GameState) Logger() *turnlogger.TurnLogger     { return s.logger }
func (s *GameState) SetLogger(l *turnlogger.TurnLogger) { s.logger = l }
func (s *GameState) LogEntries() []turnlogger.LogEntry  { return s.logger.Entries() }

func (s *GameState) TriggeringCard() card.Card                  { return s.triggeringCard }
func (s *GameState) SetTriggeringCard(c card.Card)              { s.triggeringCard = c }
func (s *GameState) AttackReactionTarget() *card.CardState      { return s.attackReactionTarget }
func (s *GameState) SetAttackReactionTarget(cs *card.CardState) { s.attackReactionTarget = cs }

func (s *GameState) ActionPoints() int     { return s.actionPoints }
func (s *GameState) SetActionPoints(n int) { s.actionPoints = n }
func (s *GameState) AddActionPoints(n int) { s.actionPoints += n }

func (s *GameState) Value() int     { return s.value }
func (s *GameState) SetValue(v int) { s.value = v }
func (s *GameState) AddValue(n int) { s.value += n }

func (s *GameState) CardsDrawn() int     { return s.cardsDrawn }
func (s *GameState) SetCardsDrawn(n int) { s.cardsDrawn = n }

func (s *GameState) IncomingDamage() int     { return s.incomingDamage }
func (s *GameState) SetIncomingDamage(n int) { s.incomingDamage = n }

func (s *GameState) ArcaneIncomingDamage() int     { return s.arcaneIncomingDamage }
func (s *GameState) SetArcaneIncomingDamage(n int) { s.arcaneIncomingDamage = n }

func (s *GameState) BlockTotal() int     { return s.blockTotal }
func (s *GameState) SetBlockTotal(n int) { s.blockTotal = n }

func (s *GameState) ArcaneDamageDealt() bool     { return s.arcaneDamageDealt }
func (s *GameState) SetArcaneDamageDealt(v bool) { s.arcaneDamageDealt = v }

func (s *GameState) OpponentMarked() bool     { return s.opponentMarked }
func (s *GameState) SetOpponentMarked(v bool) { s.opponentMarked = v }
func (s *GameState) MarkOpponent()            { s.opponentMarked = true }
func (s *GameState) ClearOpponentMarked()     { s.opponentMarked = false }

func (s *GameState) AuraCreated() bool     { return s.auraCreated }
func (s *GameState) SetAuraCreated(v bool) { s.auraCreated = v }

func (s *GameState) CardBanished() bool     { return s.cardBanished }
func (s *GameState) SetCardBanished(v bool) { s.cardBanished = v }

func (s *GameState) NonAttackActionPlayed() bool     { return s.nonAttackActionPlayed }
func (s *GameState) SetNonAttackActionPlayed(v bool) { s.nonAttackActionPlayed = v }

func (s *GameState) CurrentStepRerouted() bool     { return s.currentStepRerouted }
func (s *GameState) SetCurrentStepRerouted(v bool) { s.currentStepRerouted = v }

// AmendLastChainStepN adds n to the most recent ChainStep entry's N field. No-op when
// the logger is nil or when no chain-step entry exists yet.
func (s *GameState) AmendLastChainStepN(n int) { s.logger.AmendLastChainStepN(n) }

// HeroWantsLowerHealth reports whether the current hero opts into the LowerHealthWanter
// marker. Returns false when no hero is set.
func (s *GameState) HeroWantsLowerHealth() bool {
	if s.hero == nil {
		return false
	}
	_, ok := s.hero.(LowerHealthWanter)
	return ok
}

// CurrentHeroClass returns the active hero's primary class. Zero when no hero is set.
func (s *GameState) CurrentHeroClass() card.CardType {
	if s.hero == nil {
		return 0
	}
	return s.hero.Class()
}

// HasPlayedType reports whether any card played this turn has the given type. Universal
// cards' Types() folds the active hero's class through the engine, but at the GameState
// level we just check the recorded types directly.
func (s *GameState) HasPlayedType(t card.CardType) bool {
	for _, c := range s.cardsPlayed {
		if c.Types(nil).Has(t) {
			return true
		}
	}
	return false
}

// appendCopy is the small slice-clone helper Copy uses. Inlined into its own helper so
// each persistent slice's branch is a one-liner.
func appendCopy(dst []card.Card, src []card.Card) []card.Card {
	if len(src) == 0 {
		return dst
	}
	return append(dst, src...)
}
