// Package card defines the Card interface every Flesh and Blood card implements, the
// per-chain-step CardState wrapper that carries mutable flags between resolution
// phases, and the narrow GameEngine / Logger / Aura / Trigger interfaces cards
// consume from the sim.
//
// The package owns the contract; it does NOT import the sim. *sim.TurnState satisfies
// GameEngine structurally and *turnlogger.TurnLogger satisfies Logger structurally,
// so the sim hands either through cards via the explicit Play / Block / OnHit / Cost
// args without any adapter in the middle. Cards interact with the GameEngine /
// Logger / Aura / Trigger interfaces only, freeing the sim package to evolve its
// concrete representations without breaking the card contract.
package card

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
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

	// Auras: per-trigger-type registration. The engine builds the underlying aura
	// internally; cards supply only the handler and an initial count. Source is
	// derived from self.Card.
	AddStartOfTurnAura(self *CardState, handler AuraHandler, count int)
	AddAttackActionAura(self *CardState, handler AuraHandler, count int)
	AddOncePerTurnAttackActionAura(self *CardState, handler AuraHandler, count int)
	AddAttackAura(self *CardState, handler AuraHandler, count int)
	AddEndOfTurnAura(self *CardState, handler AuraHandler, count int)

	// Triggers: one-shot, per-trigger-type. AddHitTrigger's filter narrows the
	// firing event to a card-type predicate (typically TypeSet.IsAttack or
	// TypeSet.IsAttackAction); nil = no filter.
	AddHitTrigger(self *CardState, handler TriggerHandler, filter func(TypeSet) bool)
	AddEndOfTurnTrigger(self *CardState, handler TriggerHandler)

	// Token economy
	CreateRunechants(int)
	CreatePonder(int)
	CreateGold(int)
	CreateSilver(int)
	CreateCopper(int)
	Runechants() int
	Ponders() int
	Gold() int
	Silver() int
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
	HasPlayedType(CardType) bool
	CardsPlayed() []Card
	SetCardsPlayed([]Card)
	CardsRemaining() []*CardState
	TriggeringCard() Card
	// AuraCount is the count of live auras — used by gates like Yinti Yanti's "while
	// you control an aura" rider. Cards don't get a typed slice view; the engine
	// owns the live aura set.
	AuraCount() int

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
	AppendPreTriggerf(source string, n int, format string, args ...any)
	AmendLastChainStepN(n int)
}

// Aura is the minimal view cards' aura handlers see of the firing aura. The engine
// constructs and tracks the full concrete representation internally; the handler
// reads counts / self-identity through this interface and ends the aura's life via
// Destroy.
type Aura interface {
	// Count returns the aura's current count. Handlers that fire-then-destroy read
	// it as a payload (e.g. Blessing of Occult's runechant count); multi-fire auras
	// read it as fires-remaining.
	Count() int
	// DecrementCount decrements the count by 1 and returns the new value. Multi-fire
	// aura handlers (Malefic / Runeblood pattern) call this and Destroy when the
	// result reaches 0.
	DecrementCount() int
	// SelfName is the originating card or token's display name — used for log
	// attribution from inside the handler.
	SelfName() string
	// SelfCardID is the originating card's registry ID, or ids.InvalidCard when the
	// aura belongs to a token. Handlers that match against a specific source card
	// (e.g. Sigil of Silphidae's "destroy a different aura") gate on this.
	SelfCardID() ids.CardID
	// Destroy ends the aura's life. addToGraveyard sends the originating card to
	// the graveyard (token auras with no originating card skip the append).
	Destroy(addToGraveyard bool)
}

// Trigger is the minimal view cards' one-shot trigger handlers see of the firing
// trigger. Today only SourceName is needed; expand as more triggers appear.
type Trigger interface {
	SourceName() string
}

// AuraHandler is the card-facing aura handler signature. Cards write functions of
// this type and pass them to GameEngine.AddXxxAura.
type AuraHandler func(g GameEngine, l Logger, a Aura)

// TriggerHandler is the card-facing one-shot trigger handler signature. Cards write
// functions of this type and pass them to GameEngine.AddHitTrigger (and any future
// AddXxxTrigger method).
type TriggerHandler func(g GameEngine, l Logger, t Trigger)
