package card

// Universal is the marker for cards whose printed type-line is "extended" by the active
// hero's class — Wage Gold's "Universal" keyword reads as "this counts as the hero's
// class for class-gated triggers" (Viserai's Runeblade trigger fires on a Universal
// Wage Gold). Universal cards ask the engine for the active hero's class inside their
// own Types(g) body; v2/card holds no hero state.
type Universal interface {
	Universal()
}
