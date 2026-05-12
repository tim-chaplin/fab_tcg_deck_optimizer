package gameengine

import "github.com/tim-chaplin/fab-deck-optimizer/v2/card"

// DealArcaneDamage credits n arcane damage to Value, writes a "Dealt n arcane damage"
// rider line under source, and flips ArcaneDamageDealt when LikelyDamageHits(n, false)
// approves so same-turn triggers reading "if you've dealt arcane damage this turn" fire.
// Routes through dealtArcaneText[n] so the hot path avoids per-call fmt.Sprintf and
// variadic-int boxing.
func (g *GameEngine) DealArcaneDamage(l card.Logger, source string, n int) {
	g.AddValue(n)
	if g.LikelyDamageHits(n, false) {
		g.arcaneDamageDealt = true
	}
	if n >= 0 && n < len(dealtArcaneText) {
		l.AppendPostTrigger(source, dealtArcaneText[n], n)
		return
	}
	l.AppendPostTriggerf(source, n, "Dealt %d arcane damage", n)
}

// dealtArcaneText is the pre-built rider-line cache indexed by arcane-damage count, keeping
// DealArcaneDamage alloc-free on the hot path. Extend if a new card prints higher arcane.
var dealtArcaneText = [...]string{
	0: "Dealt 0 arcane damage",
	1: "Dealt 1 arcane damage",
	2: "Dealt 2 arcane damage",
	3: "Dealt 3 arcane damage",
	4: "Dealt 4 arcane damage",
}
