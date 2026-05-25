package gameengine

// Pool is a fixed-capacity recycler of *GameState values, sharing one set of slice
// backings across acquisitions to amortise the per-permutation allocation cost in the
// chain runner's hot loop.
//
// The Pool is single-goroutine: callers must not share a Pool across goroutines (the
// parallel evaluator already gives each worker its own scratch, so this matches the
// existing locality). Get / Put are O(1) slice-stack ops with no locking.
//
// Get / Put / FreeAll maintain an in-flight counter (Gets minus Puts since the last
// FreeAll). Get panics when in-flight would exceed the Pool's capacity — that's a hard
// signal to the operator that the workload outran the budgeted cap and we'd otherwise
// silently start allocating fresh state every call. Sizing the cap: run a workload,
// read HighWaterMark, set cap to peak + ~50% headroom.
//
// State is NOT wiped on Put; the next Get's caller is responsible for re-initialising
// (typically via CopyFrom / CopyPersistentStateFrom). Slice backings on the recycled
// GameState (hand / graveyard / banished / auras / items) survive Put → Get, so the
// follow-up CopyPersistentStateFrom's cap checks reuse them without reallocating —
// that backing reuse is the whole point of pooling here.
type Pool struct {
	free          []*GameState
	cap           int
	inFlight      int // Get-Put balance since the last FreeAll.
	peakSinceFree int // max inFlight observed since the last FreeAll.
	highWater     int // max peakSinceFree observed over all cycles.
}

// NewPool returns a Pool with the given maximum simultaneous-checkout capacity. Get
// allocates lazily up to cap; beyond cap it panics. cap must be > 0.
func NewPool(cap int) *Pool {
	if cap <= 0 {
		panic("gameengine.NewPool: cap must be > 0")
	}
	return &Pool{
		free: make([]*GameState, 0, cap),
		cap:  cap,
	}
}

// Get returns a *GameState from the free list, allocating a fresh zero-value GameState
// when the list is empty. The returned state's contents are unspecified — the caller
// must initialise it (typically via CopyFrom / CopyPersistentStateFrom) before use.
//
// Panics when in-flight (Gets minus Puts since the last FreeAll) would exceed the cap
// the Pool was constructed with. The fix is to either tighten the caller's Put
// discipline or raise the cap to match a measured HighWaterMark + headroom.
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
		return new(GameState)
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

// FreeAll declares that every *GameState the Pool has handed out is no longer in use
// and resets the in-flight counter to zero. Slice backings on the recycled states are
// retained, so subsequent Gets reuse them. Callers must ensure no live reference to any
// previously-Got state remains when FreeAll runs — typically called at a clean lifecycle
// boundary (end of a shuffle, end of a deck eval).
func (p *Pool) FreeAll() {
	if p.peakSinceFree > p.highWater {
		p.highWater = p.peakSinceFree
	}
	p.peakSinceFree = 0
	p.inFlight = 0
}

// HighWaterMark returns the peak inFlight checkout count seen between any two FreeAll
// calls — the bound a fixed-capacity Pool would need. Used to size a pre-allocated Pool
// by running a workload, reading the peak, and budgeting cap = peak + headroom.
func (p *Pool) HighWaterMark() int {
	peak := p.highWater
	if p.peakSinceFree > peak {
		peak = p.peakSinceFree
	}
	return peak
}
