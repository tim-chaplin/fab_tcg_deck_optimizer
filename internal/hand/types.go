package hand

// Turn-summary data shapes returned by Best: Role plus the per-card and aggregate records that
// encode what the solver chose to do with a hand.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Role is what a card does on a given turn cycle.
type Role uint8

const (
	Pitch Role = iota
	Attack
	Defend
	Held
	// Arsenal marks the card placed into the arsenal at end of turn. Contributes nothing to this
	// turn's Value (it scores on the future turn it's played) and carries across the turn boundary
	// as Play.ArsenalCard. At most one card in BestLine takes this role; when the partition
	// enumerator leaves arsenal empty, Best post-hoc promotes one Held card.
	Arsenal
)

// CardAssignment is a single card + the role it took this turn. Hand cards produce one per card;
// an arsenal-in card contributes one with FromArsenal set so a turn fits in one slice.
// Contribution is the per-card credit toward TurnSummary.Value (damage dealt, block share, or
// pitch resource depending on Role), filled by fillContributions once the winner is picked.
type CardAssignment struct {
	Card         card.Card
	Role         Role
	FromArsenal  bool
	Contribution float64
}

// TurnSummary is the result of running Best on a hand: the winning card-role assignments plus
// aggregate metadata about the turn.
type TurnSummary struct {
	// BestLine is the winning partition. Hand cards come first in canonical (post-sort) order;
	// the previous-turn arsenal card, if any, is the last entry with FromArsenal=true. Never
	// mutate — memoized results alias this slice.
	BestLine []CardAssignment
	// AttackChain is the winning chain in play order: attack-role cards from BestLine interleaved
	// with any swung weapons at the positions the solver picked. Each entry carries its Play-time
	// damage (plus hero-trigger damage for cards) so callers can attribute contribution for
	// weapons, which have no BestLine entry. Swung weapons are recoverable by type-asserting
	// AttackChainEntry.Card to weapon.Weapon. Empty when no attacks were played.
	AttackChain []AttackChainEntry
	// Value is the turn's total score (damage dealt + damage prevented). Equals the sum of
	// BestLine[].Contribution plus any weapon-swing damage (weapons aren't in BestLine).
	Value int
	// LeftoverRunechants is the Runechant token count at end of the chosen chain; the caller
	// feeds it back as the next turn's carryover. For partitions with no attacks, equals the
	// carryover the caller passed in.
	LeftoverRunechants int
	// ArsenalCard is the card occupying the arsenal slot at end of turn — either a hand card
	// just arsenaled (role=Arsenal) or a previous-turn arsenal card that stayed. Nil when empty.
	// The caller feeds this back as the next turn's arsenalCardIn.
	ArsenalCard card.Card
	// Drawn is the cards the winning chain drew mid-turn, in draw order, each paired with the
	// disposition the solver picked for it. Populated from state.Drawn during fillContributions's
	// tracked replay. Role is one of Pitch (consumed to fund the chain, Contribution = Pitch()),
	// Attack (played as a free-cost chain extension, Contribution = damage dealt), Arsenal
	// (promoted into an empty slot post-enumeration, Contribution 0), or Held (carries into
	// the next hand, Contribution 0). Nil when no draw rider fired.
	Drawn []CardAssignment
	// AuraTriggers is the surviving AuraTrigger set at end of this turn — triggers added by
	// this turn's winning Play chain. The deck loop feeds this into next turn's start-of-turn
	// trigger pass, closing the cross-turn loop. Nil when the turn played no trigger-creating
	// aura.
	AuraTriggers []card.AuraTrigger
	// TriggersFromLastTurn records the AuraTriggers whose start-of-turn handlers fired at the
	// top of this turn, each with the damage-equivalent it credited. Populated by the deck
	// loop before FormatBestTurn is called; Value already includes the sum.
	TriggersFromLastTurn []TriggerContribution
	// StartOfTurnAuras lists the aura cards that were in play at the top of this turn — one
	// entry per AuraTrigger carried in from the previous turn (so an aura that registered
	// multiple triggers appears multiple times, and two copies of the same aura show twice).
	// Populated by the deck loop before the start-of-turn fires run; surfaced in FormatBestTurn
	// so the reader can see which carryover auras fed mid-chain "(+M aura trigger)" damage.
	StartOfTurnAuras []card.Card
	// ReturnedToTopOfDeck surfaces card.TurnState.ReturnedToTopOfDeck from the winning permutation; see
	// that field for the alt-cost contract (deck loop skips Held → nextHeld carries for
	// these cards, plus inserts each at the next-turn deck top).
	ReturnedToTopOfDeck []card.Card
	// DeckRemoved surfaces card.TurnState.DeckRemoved from the winning permutation; see
	// that field for the buf-patch contract.
	DeckRemoved []card.Card
}

// TriggerContribution is one start-of-turn AuraTrigger fire: the aura that fired plus the
// Damage it credited (folded into Value) and the card (if any) the handler revealed onto the
// hand. Surfaced in TurnSummary.TriggersFromLastTurn so FormatBestTurn can print a "(from
// previous turn)" line naming the outcome.
//
// Revealed is the deck-top card the handler moved into the hand (e.g. Sigil of the Arknight
// revealing an attack action). Nil when the handler didn't reveal anything.
type TriggerContribution struct {
	Card     card.Card
	Damage   int
	Revealed card.Card
}

// AttackChainEntry is a single played attack — a card with role=Attack or a swung weapon —
// carrying the damage it contributed when it resolved in the winning chain. Damage is the Play()
// return; TriggerDamage is the hero's OnCardPlayed contribution (e.g. Viserai creating a
// Runechant); AuraTriggerDamage is the mid-chain AuraTrigger contribution (e.g. a Malefic
// Incantation played on a prior turn whose TriggerAttackAction fires when this card resolves).
// The two trigger buckets are split so the display can attribute hero OnCardPlayed damage and
// mid-chain aura damage on their own lines. For BestLine Attack entries Damage + TriggerDamage
// + AuraTriggerDamage equals CardAssignment.Contribution; weapons live only here.
type AttackChainEntry struct {
	Card              card.Card
	Damage            float64
	TriggerDamage     float64
	AuraTriggerDamage float64
}

// ArsenalIn returns the assignment for the card that started the turn in the arsenal, if any.
// Lets callers treat the arsenal-in card differently from hand cards without scanning BestLine.
func (t TurnSummary) ArsenalIn() (CardAssignment, bool) {
	for _, a := range t.BestLine {
		if a.FromArsenal {
			return a, true
		}
	}
	return CardAssignment{}, false
}

// String returns a human-readable role name ("PITCH", "ATTACK", "DEFEND", "HELD", "ARSENAL").
func (r Role) String() string {
	switch r {
	case Pitch:
		return "PITCH"
	case Attack:
		return "ATTACK"
	case Defend:
		return "DEFEND"
	case Held:
		return "HELD"
	case Arsenal:
		return "ARSENAL"
	}
	return "UNKNOWN"
}
