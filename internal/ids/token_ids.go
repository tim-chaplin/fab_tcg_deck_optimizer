package ids

// Token ability IDs. Anchored after the weapon ability IDs in the shared CardID space —
// each token's activated ability (Gold spend, Silver spend, Copper spend) flows through
// the attack-turn runner like a weapon ability, so cardMetaCache keys off the same CardID space.
const (
	GoldTokenAbilityID CardID = TalisharAbilityID + iota + 1
	SilverTokenAbilityID
	CopperTokenAbilityID
)

// Token aura / item IDs. The token itself (not its activated ability) carries this ID so
// cache keys distinguish "Runechant aura" from "Ponder aura" without a separate TokenKind
// field. Anchored after the token ability IDs.
const (
	RunechantTokenID CardID = CopperTokenAbilityID + iota + 1
	PonderTokenID
	GoldTokenID
	SilverTokenID
	CopperTokenID
)
