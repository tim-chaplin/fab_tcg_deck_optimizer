// Package card defines the Card interface every Flesh and Blood card implements, the
// per-chain-step CardState wrapper that carries mutable flags between resolution
// phases, and the narrow GameEngine / Logger interfaces cards consume from the sim.
//
// The package owns the contract; it does NOT import the sim. *sim.TurnState satisfies
// GameEngine structurally and *turnlogger.TurnLogger satisfies Logger structurally,
// so the sim hands either through cards via the explicit Play / Block / OnHit / Cost
// args without any adapter in the middle. Cards interact with the GameEngine and
// Logger interfaces only, freeing the sim package to evolve its concrete TurnState
// representation without breaking the card contract.
package card

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// GameEngine is the cards-facing rules-engine handle the sim threads through every
// Card hook. Every method cards invoke during Card.Play, OnHit handlers, Trigger
// handlers, Block hooks, and Cost queries is declared here — the surface combines
// state queries (Hand, Graveyard, Runechants, …) with active-effect operations
// (DrawOne, CreateRunechants, AddValue, Opt, Clash, …). The sim's *TurnState
// satisfies it structurally.
//
// The interface is intentionally method-only: it exposes no fields. Sim-internal
// code reads and writes TurnState fields directly; card code that needs raw field
// access type-asserts the GameEngine value back to *sim.TurnState and lives with
// the sim import. Cards-side field access is the exception, not the rule — most
// cards reach only the methods below.
//
// Methods taking sim-owned shapes (Aura, Trigger) are typed as `any` so the
// interface stays sim-free. Cards still import sim to construct the concrete values
// they pass through; the type assertion happens inside sim.TurnState's
// implementation. Tightening these to local marker interfaces is a future cleanup.
type GameEngine interface {
	// Zones
	Hand() []Card
	AppendHand(Card)
	PopHandAt(int) Card
	Deck() *deck.Deck
	PeekDeck() (Card, bool)
	PeekTopN(int) []Card
	PrependToDeck(Card)
	RecycleToDeckBottom(*CardState)
	Graveyard() []Card
	AddToGraveyard(Card)
	BanishFromGraveyard(func(Card) bool) (Card, bool)
	Banished() []Card

	// Auras / Triggers (opaque — cards still import sim to construct the concrete
	// sim.Aura / sim.Trigger values these methods accept).
	AddAura(any)
	AddTrigger(any)
	DestroyAura(any, bool)

	// Token economy
	CreateRunechants(int)
	CreatePonder(int)
	CreateGold(int)
	CreateCopper(int)
	Runechants() int
	Ponders() int
	Gold() int
	Copper() int

	// Value crediting and arcane damage
	Value() int
	AddValue(int)
	SetValue(int)
	DealArcaneDamage(Logger, *CardState, int)

	// AP / Overpower (chain-step controls cards read or grant).
	ActionPoints() int
	AddActionPoints(int)
	Overpower() bool
	SetOverpower(bool)

	// Sticky flags cards read to gate their effects (and flip when their effects fire).
	ArcaneDamageDealt() bool
	SetArcaneDamageDealt(bool)
	AuraCreated() bool
	SetAuraCreated(bool)
	CardBanished() bool
	NonAttackActionPlayed() bool
	OpponentMarked() bool
	SetOpponentMarked(bool)

	// Partition / matchup state.
	IncomingDamage() int
	ArcaneIncomingDamage() int
	BlockTotal() int
	Defenders() []Card
	Pitched() []Card

	// Chain queries
	HasPlayedOrCreatedAura() bool
	HasPlayedType(card.CardType) bool
	CardsPlayed() []Card
	SetCardsPlayed([]Card)
	CardsRemaining() []*CardState
	TriggeringCard() Card
	// Auras / Triggers slice access is omitted — Aura and Trigger are sim-owned shapes
	// (mirrored as `any` in AddAura / AddTrigger / DestroyAura). The handful of cards
	// that scan the live aura set (e.g. Yinti Yanti's len(Auras) gate) type-assert back
	// to *sim.TurnState; a future cleanup could expose typed slice accessors once Aura
	// / Trigger move to a shared package.

	// Mid-chain draw / tutoring / recycling
	DrawOne()
	TutorFromDeck(func(Card) int) (Card, bool)
	RecycleFromGraveyardToTop(func(Card) bool) (Card, bool)
	RecycleFromGraveyardToBottom(func(Card) bool) (Card, bool)

	// Opt (hero-driven top-of-deck reshape)
	Opt(Logger, int)

	// Clash (top-of-deck power compare)
	Clash(win, lose func())

	// Attack reaction target accessor.
	// TODO: evaluate if we actually need this. It leaks an engine implementation
	// detail (that the engine has to remember which attack is on the stack while
	// it's processing attack reactions). A more natural shape would be passing the
	// target directly into the AttackReaction's Play.
	AttackReactionTarget() *CardState
}

// Logger is the cards-facing log sink the chain runner threads through every Card
// hook. The method set matches *turnlogger.TurnLogger structurally so the sim hands
// one of those directly to cards. A typed-nil receiver is the find-best skip sentinel
// — every method short-circuits inside the implementation, so cards never gate at the
// call site.
//
// Cards reach for AppendPostTrigger / AppendPostTriggerf for self-riders ("Created a
// runechant"), AppendPreTrigger for hero or aura attack-action triggers. The
// AppendChainStep / AppendChainStepf / AmendLastChainStepN trio is sim-internal —
// ResolveChainStep emits the chain step after Play returns, the Opt helper emits its
// own free-form chain entry, and GrantAttackReactionBuff amends the buffed attack's
// delta — but the methods sit on Logger so the same value can flow through both
// card-side and sim-side call sites.
type Logger interface {
	AppendChainStep(text string, n int)
	AppendChainStepf(n int, format string, args ...any)
	AppendPostTrigger(source, text string, n int)
	AppendPostTriggerf(source string, n int, format string, args ...any)
	AppendPreTrigger(source, text string, n int)
	AmendLastChainStepN(n int)
}
