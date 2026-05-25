package gameengine

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
