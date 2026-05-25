package gameengine

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// Pool is a fixed-capacity recycler of *GameState values. All cap slots are built
// eagerly by factory at NewPool time; Get / Put cycle the same pointers forever with no
// further allocation. Single-goroutine; Get / Put are O(1) unlocked slice-stack ops.
//
// State is NOT wiped on Put — the next Get's caller re-initialises (CopyFrom /
// CopyPersistentStateFrom / ResetEphemeralState / Reset). Slice backings survive
// Put → Get, so a factory that pre-sizes hand / cardsPlayed eliminates per-perm growth.
//
// Two invariants guard against silent fallback to per-call allocation:
//   - Get past cap panics — workload outran the budget; raise NewPool cap.
//   - Get on an empty free list panics — a caller leaked a slot past its scope.
type Pool struct {
	free          []*GameState
	cap           int
	inFlight      int // Get-Put balance since the last FreeAll.
	peakSinceFree int // max inFlight observed since the last FreeAll.
	highWater     int // max peakSinceFree observed over all cycles.
}

// NewPool returns a Pool with cap simultaneous-checkout capacity, eagerly allocating
// cap states up front via factory. cap must be > 0; factory must not return nil.
func NewPool(cap int, factory func() *GameState) *Pool {
	if cap <= 0 {
		panic("gameengine.NewPool: cap must be > 0")
	}
	if factory == nil {
		panic("gameengine.NewPool: factory must not be nil")
	}
	free := make([]*GameState, cap)
	for i := range free {
		s := factory()
		if s == nil {
			panic("gameengine.NewPool: factory returned nil")
		}
		free[i] = s
	}
	return &Pool{free: free, cap: cap}
}

// Get pops a *GameState off the free list. Contents are unspecified — the caller
// re-initialises (CopyFrom / CopyPersistentStateFrom / ResetEphemeralState / Reset).
// Panics on cap overrun or on empty-free-list-within-budget (a leaked Got slot).
func (p *Pool) Get() *GameState {
	p.inFlight++
	if p.inFlight > p.cap {
		panic("gameengine.Pool: in-flight checkouts exceeded cap; raise NewPool cap to match the workload's peak")
	}
	if p.inFlight > p.peakSinceFree {
		p.peakSinceFree = p.inFlight
	}
	n := len(p.free)
	if n == 0 {
		panic("gameengine.Pool: free list empty within budget; a caller leaked a Got slot past its scope (audit Get/Put pairing)")
	}
	s := p.free[n-1]
	p.free[n-1] = nil
	p.free = p.free[:n-1]
	return s
}

// Put returns s to the free list for later reuse. Callers must ensure no other live
// reference to s remains — once Put, the Pool may hand s out to any subsequent Get and
// the new owner is free to overwrite its fields. Putting nil is a no-op.
func (p *Pool) Put(s *GameState) {
	if s == nil {
		return
	}
	p.inFlight--
	p.free = append(p.free, s)
}

// FreeAll resets the in-flight accounting to zero so the next cycle's peak measurement
// starts clean. It does NOT reclaim slots from the free list — callers must Put every
// Got slot within scope. Typically called at end-of-shuffle.
func (p *Pool) FreeAll() {
	if p.peakSinceFree > p.highWater {
		p.highWater = p.peakSinceFree
	}
	p.peakSinceFree = 0
	p.inFlight = 0
}

// HighWaterMark returns the peak inFlight checkout count seen between any two FreeAll
// calls — the bound for sizing a future Pool's cap (peak + headroom).
func (p *Pool) HighWaterMark() int {
	peak := p.highWater
	if p.peakSinceFree > peak {
		peak = p.peakSinceFree
	}
	return peak
}

// Worst-case sizes for pooled GameState slice backings. Bounded by FaB hero intellect
// (handSize ≤ 7), equipment slots (weapons ≤ 4), and chain-runner mid-turn-drawn
// headroom (32 extra attacker slots). hand carries (held + chain + pitch); cardsPlayed
// grows for rider plays / aura fires; graveyard absorbs prior-turn grav + this turn's
// played cards.
const (
	maxHandSize           = 7
	maxWeapons            = 4
	maxDrawnExtra         = 32
	maxAttackers          = maxHandSize + maxWeapons + 1 + maxDrawnExtra
	defaultHandCap        = 2*maxHandSize + maxAttackers
	defaultCardsPlayedCap = 2 * (maxHandSize + maxAttackers)
	// maxDeckSize bounds the worst-case deck the sim runs against — Classic Constructed
	// is 60 minimum with no formal max; 80 leaves headroom for oversized lists. Graveyard
	// and banished can only ever hold cards that originated in the deck (plus a handful
	// of starting equipment), so both share the same bound.
	maxDeckSize         = 80
	defaultGraveyardCap = maxDeckSize + 10
	defaultBanishedCap  = maxDeckSize + 10
	// defaultPoolCap is the prewarmed Pool size. Measured peak in-flight on a 200-shuffle
	// Viserai run is 6; this cap leaves headroom for higher-fanout decks.
	defaultPoolCap = 24
)

// NewPrewarmedState returns a zero-value *GameState with hand / cardsPlayed / graveyard
// / banished slice backings pre-allocated to worst-case sizes, so the chain runner's
// per-permutation fills never need to grow past the pool slot's existing cap.
func NewPrewarmedState() *GameState {
	s := new(GameState)
	s.SetHandStates(make([]card.CardState, 0, defaultHandCap))
	s.SetCardsPlayed(make([]card.Card, 0, defaultCardsPlayedCap))
	s.SetGraveyard(make([]card.Card, 0, defaultGraveyardCap))
	s.SetBanished(make([]card.Card, 0, defaultBanishedCap))
	return s
}

// NewPrewarmedPool returns a Pool of prewarmed GameStates sized for the chain runner's
// worst case. The standard entry point for sim Evaluators.
func NewPrewarmedPool() *Pool {
	return NewPool(defaultPoolCap, NewPrewarmedState)
}
