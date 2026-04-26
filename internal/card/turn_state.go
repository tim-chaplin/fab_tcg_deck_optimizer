package card

import "fmt"

// Per-turn shared context threaded through Card.Play. Cards mutate state directly — moving
// cards between Hand / Deck / Graveyard / Banish, registering triggers, creating runechants
// — and the sim copies the winning permutation's final state into next-turn state. There's
// no diff-signal indirection: a card that wants to draw appends to s.Hand and pops from
// s.Deck, full stop.
//
// Persistent fields (Hand, Deck, Arsenal, Graveyard, Banish, Runechants, AuraTriggers)
// carry across turns when the sim adopts the winner's snapshot. Transient fields
// (CardsPlayed, Pitched, IncomingDamage, etc.) are seeded by the sim per chain-step and
// reset at the turn boundary.

// LogEntryKind classifies a LogEntry. Triggers come in two flavours because they fire on
// opposite sides of their parent in the FaB stack — the format layer needs to know which
// side a given entry sits on so two cards with the same display name in the same chain
// don't steal each other's triggers during grouping.
type LogEntryKind int8

const (
	// LogEntryChainStep is the sim's "<Card>: <VERB>" line. Stands alone in the printout
	// and acts as the parent that triggers attach to.
	LogEntryChainStep LogEntryKind = iota
	// LogEntryPreTrigger is a trigger that fires before its parent chain entry resolves
	// (hero / aura attack-action triggers). The format layer attaches it to the next
	// chain entry whose name matches Source.
	LogEntryPreTrigger
	// LogEntryPostTrigger is a trigger that fires after its parent chain entry resolves
	// (ephemeral attack triggers). The format layer attaches it to the previous chain
	// entry whose name matches Source.
	LogEntryPostTrigger
)

// LogEntry is one chain-event entry in TurnState.Log. Text is the freeform display string
// the producer authored ("Viserai created a runechant", "Consuming Volition [R]: ATTACK")
// — the format layer renders it verbatim, no further opinions on phrasing. Kind tags the
// entry as a chain step or as a pre/post trigger so the grouping algorithm can attribute
// triggers correctly even when sibling chain entries share a name. Source names the card
// whose play caused the trigger and is matched against chain-entry names to pick the
// parent; only meaningful for triggers. N is the damage-equivalent credited to s.Value
// when the entry was added.
type LogEntry struct {
	Text   string
	Source string
	Kind   LogEntryKind
	N      int
}

// TurnState is the shared turn-level context passed to Card.Play alongside the per-card
// CardState wrapper.
type TurnState struct {
	// Hand is the cards currently in hand. Starts as the dealt hand minus pitched / attacker
	// / defender cards (those have been routed by the partition). Cards that draw or tutor
	// append to Hand; alt-cost effects pop from Hand. Whatever's in Hand at end of chain
	// becomes next turn's Held cards.
	Hand []Card
	// Deck is the deck top-to-bottom. Cards mutate freely: DrawOne pops Deck[0]; tutor
	// removes a specific card; alt cost prepends to Deck. Whatever's in Deck at end of
	// chain becomes next turn's deck.
	Deck []Card
	// Arsenal is the arsenal slot's contents at this point in the chain — the arsenal-in
	// card at start of turn, nil after it plays / defends, refilled post-chain by the
	// arsenal-promotion step. Cards that read "from arsenal" use CardState.FromArsenal,
	// not this field.
	Arsenal Card
	// Graveyard is cards that have entered the graveyard this turn — every card played or
	// blocked lands here after resolving. Pitched cards do not (they recycle to deck
	// bottom). Cards that destroy themselves mid-chain route through AddToGraveyard.
	Graveyard []Card
	// Banish holds cards moved into the banished zone this turn (e.g. an aura-banish-for-
	// arcane rider).
	Banish []Card
	// Runechants is the live count of Runechant aura tokens in play. Carries across turns.
	// CreateRunechants increments it; the attack pipeline consumes the running total on each
	// attack / weapon swing.
	Runechants int
	// ArcaneDamageDealt sticks true once any source of arcane damage fires this turn:
	// a Runechant token consuming itself on an attack / weapon swing, or a card whose Play
	// deals arcane directly. Effects that read "if you've dealt arcane damage this turn"
	// consult this flag rather than Runechants. Reset at turn boundary.
	ArcaneDamageDealt bool
	// AuraTriggers is the list of triggers from auras currently in play. Cards add entries
	// during Play via AddAuraTrigger; the sim fires matching entries on each trigger-Type
	// condition, decrements Count in place, and drops entries whose Count hits zero after
	// sending Self to the graveyard. Carries across turns.
	AuraTriggers []AuraTrigger

	// --- Transient: reset by the sim per turn / chain step ---

	// Value is the running damage-equivalent total for this chain — damage dealt + damage
	// prevented + every aura-token / hero-trigger credit. The dispatcher records the chain
	// step's Play+BonusAttack contribution via AddLogEntry; trigger handlers (hero, aura,
	// ephemeral) credit themselves the same way. The solver compares permutations on this
	// field. Reset by the sim per permutation.
	Value int
	// Log is the per-event chain trace — one entry per chain step / hero / aura /
	// ephemeral / weapon swing. Chain-step producers (the sim) call AddLogEntry; pre-
	// trigger handlers (hero / aura attack-action) call AddPreTriggerLogEntry; post-
	// trigger handlers (ephemeral attack) call AddPostTriggerLogEntry. The format layer
	// uses the entry's Kind plus Source to cluster triggers under the right parent.
	// Reset per permutation.
	Log []LogEntry
	// CardsPlayed is the sequence of cards played (as attacks) this turn, in order.
	// Populated by the sim after each Play returns so later cards this turn see what was
	// played before them.
	CardsPlayed []Card
	// AuraCreated is set when a card or ability creates an aura this turn (e.g. Runechant
	// tokens). Effects that check "if you've played or created an aura this turn" should
	// OR this with CardsPlayed containing an Aura-typed card.
	AuraCreated bool
	// CardsRemaining is the cards that will be played after the current one in chain order.
	// Populated by the sim before each Play so an effect can peek forward ("next X attack")
	// or grant keywords to a later card by flipping flags on its CardState entry.
	CardsRemaining []*CardState
	// Pitched is the cards pitched this turn for resources. Populated by the sim before any
	// Play. Effects that check "if an attack card was pitched" scan this list.
	Pitched []Card
	// Overpower is set when an attack with the Overpower keyword is being played. Not yet
	// consumed by the sim — blocked damage should eventually be forwarded to the hero when
	// Overpower is true.
	Overpower bool
	// NonAttackActionPlayed is set true once any non-attack action card has been appended to
	// CardsPlayed this turn. Maintained by the chain runner so hero triggers that ask "was a
	// non-attack action played earlier?" can answer in O(1).
	NonAttackActionPlayed bool
	// IncomingDamage is the opponent damage this turn — seeded by the sim from the value
	// passed to hand.Best, and decremented by ApplyAndLogEffectiveDefense as defenders block.
	// Cards reading "did we block all incoming?" against the static partition aggregate use
	// BlockTotal instead.
	IncomingDamage int
	// BlockTotal is the sum of Defense() across every Defend-role card in the current
	// partition. Uncapped: if the partition over-blocks, BlockTotal is the full sum, not
	// clamped to IncomingDamage.
	BlockTotal int
	// EphemeralAttackTriggers are same-turn, single-fire "next attack" triggers registered
	// by a card's Play (e.g. Mauvrion Skies's "if this hits, create Runechants" rider).
	// Don't carry across turns; reset per chain.
	EphemeralAttackTriggers []EphemeralAttackTrigger
	// Revealed is the side channel start-of-turn AuraTrigger handlers use to move a card
	// from the top of the post-draw deck into the hand (Sigil of the Arknight's reveal).
	Revealed []Card
	// TriggeringCard is the card whose play caused the active aura attack-action trigger
	// to fire. The sim sets it before each AuraTrigger handler runs and clears it after;
	// the handler reads it to attribute its log line back to the triggering card. Hero
	// and ephemeral handlers receive the triggering card as a direct arg already and don't
	// need this field. Nil during direct chain-step resolution and start-of-turn fires.
	TriggeringCard Card
}

// AddLogEntry appends a freeform chain-step log line and credits n damage-equivalent to
// s.Value. text is the rendered display string. Returns the clamped n so callers can fold
// the call into a Play return. Trigger handlers call AddPreTriggerLogEntry or
// AddPostTriggerLogEntry instead so the format layer can group them under their parent.
func (s *TurnState) AddLogEntry(text string, n int) int {
	return s.appendLog(LogEntry{Text: text, Kind: LogEntryChainStep, N: n})
}

// AddPreTriggerLogEntry appends a pre-trigger log line — a hero or aura-attack-action
// trigger that fires before its parent chain entry. text is the rendered display string
// ("Viserai created a runechant"); source is the DisplayName of the card whose play
// caused the trigger to fire. The format layer attaches this entry to the next chain
// entry whose name matches source. Returns the clamped n so handlers can fold the call
// into a single return:
//
//	return s.AddPreTriggerLogEntry("Viserai created a runechant",
//	    card.DisplayName(played), s.CreateRunechant())
func (s *TurnState) AddPreTriggerLogEntry(text, source string, n int) int {
	return s.appendLog(LogEntry{Text: text, Source: source, Kind: LogEntryPreTrigger, N: n})
}

// AddPostTriggerLogEntry appends a post-trigger log line — an ephemeral attack trigger
// that fires after its parent chain entry resolves. The format layer attaches this entry
// to the previous chain entry whose name matches source. Same return contract as
// AddPreTriggerLogEntry.
func (s *TurnState) AddPostTriggerLogEntry(text, source string, n int) int {
	return s.appendLog(LogEntry{Text: text, Source: source, Kind: LogEntryPostTrigger, N: n})
}

// appendLog credits the entry's N to s.Value (clamped at 0) and appends it to s.Log.
func (s *TurnState) appendLog(e LogEntry) int {
	if e.N < 0 {
		e.N = 0
	}
	s.Value += e.N
	s.Log = append(s.Log, e)
	return e.N
}

// DrawOne models a mid-turn draw: pop the top of Deck and append it to Hand. No-op on an
// empty deck. Every draw-rider card routes through this helper.
func (s *TurnState) DrawOne() {
	if len(s.Deck) == 0 {
		return
	}
	c := s.Deck[0]
	s.Deck = s.Deck[1:]
	s.Hand = append(s.Hand, c)
}

// HasPlayedType reports whether any card played this turn has the given type in its Types() set.
func (s *TurnState) HasPlayedType(t CardType) bool {
	for _, c := range s.CardsPlayed {
		if c.Types().Has(t) {
			return true
		}
	}
	return false
}

// HasPlayedOrCreatedAura reports whether an aura was played or created this turn — the
// condition behind "if you've played or created an aura this turn" riders. The aura need
// not still be on the battlefield; the flag is sticky for the rest of the turn.
func (s *TurnState) HasPlayedOrCreatedAura() bool {
	return s.AuraCreated || s.HasPlayedType(TypeAura)
}

// ClashValue returns the net damage-equivalent of a clash (see comprehensive rules 8.5.45):
// we and the opponent reveal the top card of our decks and the higher {p} wins. We model
// from our side only — our deck's top card is read from s.Deck; the opponent's top is
// approximated as 5-power. So our {p} of 6-7 wins (credit +bonus), 5 ties (credit 0), and
// anything below 5 loses (credit -bonus). Returns 0 when s.Deck is empty.
func ClashValue(s *TurnState, bonus int) int {
	if len(s.Deck) == 0 {
		return 0
	}
	switch top := s.Deck[0].Attack(); {
	case top >= 6:
		return bonus
	case top == 5:
		return 0
	default:
		return -bonus
	}
}

// RecordValue bumps s.Value by n, clamping at 0 (FaB damage / prevention can't drive the
// running total negative). Negative n is a no-op. Cards rarely call this directly — the
// AddLogEntry / AddPreTriggerLogEntry / AddPostTriggerLogEntry helpers credit Value while
// also appending a log entry; ApplyAndLogEffectiveAttack does the same for the chain step.
func (s *TurnState) RecordValue(n int) {
	if n <= 0 {
		return
	}
	s.Value += n
}

// ApplyAndLogEffectiveAttack is the canonical chain-step finisher every Card.Play invokes:
// appends the chain-step log entry "<DisplayName>: <VERB>[ from arsenal]" (where VERB is
// ATTACK for attack actions, WEAPON ATTACK for weapons, PLAY for everything else) and
// credits Card.Attack() + self.BonusAttack to s.Value, clamped at 0. Cards with separable
// rider effects (conditional arcane bonuses, runechant creation, on-hit credits) emit
// each rider as its own post-trigger child line via LogRiderOnPlay /
// CreateAndLogRunechantsOnPlay / DealAndLogArcaneDamage so the rider's contribution is
// visible in the printout instead of bundled into the chain step's (+N).
func (s *TurnState) ApplyAndLogEffectiveAttack(self *CardState) {
	n := self.EffectiveAttack()
	if n < 0 {
		n = 0
	}
	s.AddLogEntry(chainStepText(self), n)
}

// LogPlay is the chain-step finisher for non-attack cards (auras, non-attack actions, items)
// — emits "<DisplayName>: PLAY[ from arsenal]" with no value contribution. The "(+0)" suffix
// is dropped because these cards never deal printed damage; any value they contribute lands
// via separate AddPostTriggerLogEntry / aura trigger paths.
func (s *TurnState) LogPlay(self *CardState) {
	s.AddLogEntry(chainStepText(self), 0)
}

// ApplyAndLogEffectiveDefense is the Defense Reaction counterpart to ApplyAndLogEffectiveAttack:
// emits the "<DisplayName>: DEFENSE REACTION[ from arsenal]" chain step and credits the
// effective Defense (printed Defense + BonusDefense + ArsenalDefenseBonus when from arsenal)
// to s.Value, clamped at the remaining IncomingDamage so an over-blocked DR doesn't credit
// past what was actually prevented. The credited amount is decremented from s.IncomingDamage
// so a later defender sees the reduced pool. Cards with separable rider effects (arcane
// pings, runechant creation, on-hit credits) emit each rider as its own post-trigger child
// line via DealAndLogArcaneDamage / CreateAndLogRunechantsOnPlay / LogRiderOnPlay after the
// chain step; conditional "+N{d}" bonuses fold into BonusDefense before the chain step so
// they roll into the same (+N), mirroring how BonusAttack feeds ApplyAndLogEffectiveAttack.
func (s *TurnState) ApplyAndLogEffectiveDefense(self *CardState) {
	n := self.EffectiveDefense()
	if n > s.IncomingDamage {
		n = s.IncomingDamage
	}
	if n < 0 {
		n = 0
	}
	s.IncomingDamage -= n
	s.AddLogEntry(chainStepText(self), n)
}

// chainStepText renders the "<DisplayName>: <VERB>[ from arsenal]" prefix the chain-step
// helper writes. VERB picks WEAPON ATTACK for weapon-typed cards, ATTACK for attack-action
// cards, DEFENSE REACTION for Defense Reactions, and PLAY for everything else. The "from
// arsenal" suffix tags entries played out of the arsenal slot.
func chainStepText(self *CardState) string {
	types := self.Card.Types()
	var verb string
	switch {
	case types.Has(TypeWeapon):
		verb = "WEAPON ATTACK"
	case types.IsAttackAction():
		verb = "ATTACK"
	case types.IsDefenseReaction():
		verb = "DEFENSE REACTION"
	default:
		verb = "PLAY"
	}
	if self.FromArsenal {
		verb += " from arsenal"
	}
	return DisplayName(self.Card) + ": " + verb
}

// CreateRunechants adds n Runechant token auras to the count, sets AuraCreated so effects
// that key on "aura created this turn" see it, and returns n — each token is credited as
// +1 damage at creation time. Tokens that never fire (end-of-sim leftovers) are slightly
// over-credited — accepted.
func (s *TurnState) CreateRunechants(n int) int {
	if n > 0 {
		s.AuraCreated = true
		s.Runechants += n
	}
	return n
}

// CreateRunechant is shorthand for CreateRunechants(1).
func (s *TurnState) CreateRunechant() int {
	return s.CreateRunechants(1)
}

// CreateAndLogRunechants creates n Runechant tokens, writes the canonical pre-trigger log
// line ("<selfName> created a runechant" for n==1, "<selfName> created N runechants" for
// n>1) sourced under sourceName, and returns the damage-equivalent credited. Trigger
// handlers that fire before their parent (Viserai's hero ability, Malefic Incantation's
// aura) call this in a single return statement.
func (s *TurnState) CreateAndLogRunechants(selfName, sourceName string, n int) int {
	return s.AddPreTriggerLogEntry(selfName+" "+runechantsCreatedPhrase(n), sourceName, s.CreateRunechants(n))
}

// CreateAndLogRunechantsOnHit is the post-trigger variant of CreateAndLogRunechants —
// the trigger log line reads "<selfName> created N runechants on hit" so the conditional
// gate on the ephemeral attack trigger (Mauvrion Skies, Runic Reaping) is visible in the
// printout. Same return contract as CreateAndLogRunechants.
func (s *TurnState) CreateAndLogRunechantsOnHit(selfName, sourceName string, n int) int {
	return s.AddPostTriggerLogEntry(selfName+" "+runechantsCreatedPhrase(n)+" on hit", sourceName, s.CreateRunechants(n))
}

// CreateAndLogRunechantsOnPlay is the on-play self-rider variant: the chain step's own
// "Created N runechants" sub-line, sourced under self so the format layer attaches it as
// a child of self's chain entry. The line uses indentation to convey source (no card-name
// prefix, sentence-cap leading verb) since the format layer renders it indented under
// self's chain entry. n>0 only — n=0 returns 0 without writing a line.
func (s *TurnState) CreateAndLogRunechantsOnPlay(self *CardState, n int) int {
	if n <= 0 {
		return 0
	}
	var text string
	if n == 1 {
		text = "Created a runechant"
	} else {
		text = fmt.Sprintf("Created %d runechants", n)
	}
	return s.AddPostTriggerLogEntry(text, DisplayName(self.Card), s.CreateRunechants(n))
}

// runechantsCreatedPhrase returns "created a runechant" / "created N runechants" — the
// canonical verb phrase for runechant-creation log lines.
func runechantsCreatedPhrase(n int) string {
	if n == 1 {
		return "created a runechant"
	}
	return fmt.Sprintf("created %d runechants", n)
}

// DealArcaneDamage credits n arcane damage and, when LikelyDamageHits(n, false) approves,
// flips ArcaneDamageDealt so same-turn triggers reading "if you've dealt arcane damage this
// turn" fire. The value is credited unconditionally — even if the opponent expends a card or
// resource to negate it, that's still net tempo gained — so only the trigger flag is gated by
// the hit heuristic. Returns n so callers can fold the arcane damage into their Play return.
func (s *TurnState) DealArcaneDamage(n int) int {
	if LikelyDamageHits(n, false) {
		s.ArcaneDamageDealt = true
	}
	return n
}

// DealAndLogArcaneDamage is the rider-line variant: credits n arcane damage (routed through
// DealArcaneDamage so the same hit-gated ArcaneDamageDealt flip applies) and writes a
// "Dealt N arcane damage" sub-line sourced under self so the format layer attaches it as a
// child of self's chain entry. n>0 only — n=0 returns 0 without writing a line.
func (s *TurnState) DealAndLogArcaneDamage(self *CardState, n int) int {
	if n <= 0 {
		return 0
	}
	var text string
	if n == 1 {
		text = "Dealt 1 arcane damage"
	} else {
		text = fmt.Sprintf("Dealt %d arcane damage", n)
	}
	return s.AddPostTriggerLogEntry(text, DisplayName(self.Card), s.DealArcaneDamage(n))
}

// LogRiderOnPlay writes a freeform rider sub-line under self's chain entry. text is a
// terse description of what the rider did (e.g. "On-hit discarded a card", "Gained 3
// health"); the format layer renders it indented under self's chain step, so the line
// should read as a complete utterance without a card-name prefix. n is the damage-
// equivalent credit. Returns the credited n (clamped at 0 by the underlying log helper).
func (s *TurnState) LogRiderOnPlay(self *CardState, text string, n int) int {
	return s.AddPostTriggerLogEntry(text, DisplayName(self.Card), n)
}

// AddToGraveyard appends c to s.Graveyard so later-resolving cards see it. Persistent-type
// cards (Auras, Items) don't enter the graveyard on play, so effects that destroy or banish
// themselves mid-chain route through here to make the move visible to downstream readers.
func (s *TurnState) AddToGraveyard(c Card) {
	s.Graveyard = append(s.Graveyard, c)
}

// AddAuraTrigger is the Play-side combo every Action - Aura card reaches for: flip
// AuraCreated so same-turn "if you've played or created an aura" riders see the entry, and
// append t to s.AuraTriggers so the sim fires it on its matching Type condition.
func (s *TurnState) AddAuraTrigger(t AuraTrigger) {
	s.AuraCreated = true
	s.AuraTriggers = append(s.AuraTriggers, t)
}

// AddEphemeralAttackTrigger registers a same-turn, fire-once "next attack" trigger. The sim
// stamps t.SourceIndex after the registering card's Play returns. Fires on the next
// matching attack action's resolution; fizzles silently at end of turn if no match.
func (s *TurnState) AddEphemeralAttackTrigger(t EphemeralAttackTrigger) {
	s.EphemeralAttackTriggers = append(s.EphemeralAttackTriggers, t)
}
