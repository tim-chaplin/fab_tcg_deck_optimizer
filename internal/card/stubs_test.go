package card

// stubCard is a minimal Card implementation for exercising TurnState helpers. Tests only care
// about identity and the Types / Attack / GoAgain hooks some helpers probe; everything else
// returns a zero value.
type stubCard struct {
	name    string
	types   TypeSet
	attack  int
	goAgain bool
}

func (c stubCard) ID() ID                        { return Invalid }
func (c stubCard) Name() string                  { return c.name }
func (stubCard) Cost(*TurnState) int             { return 0 }
func (stubCard) Pitch() int                      { return 0 }
func (c stubCard) Attack() int                   { return c.attack }
func (stubCard) Defense() int                    { return 0 }
func (c stubCard) Types() TypeSet                { return c.types }
func (c stubCard) GoAgain() bool                 { return c.goAgain }
func (stubCard) Play(*TurnState, *CardState) int { return 0 }

// dominatingStubCard is a stubCard that implements the Dominator marker — exercises the
// printed-Dominate branch of EffectiveDominate / HasDominate.
type dominatingStubCard struct {
	stubCard
}

func (dominatingStubCard) Dominate() {}

// notImplementedStubCard is a stubCard that implements the NotImplemented marker — exercises
// the type assertion the deck legal-pool filter keys on.
type notImplementedStubCard struct {
	stubCard
}

func (notImplementedStubCard) NotImplemented() {}
