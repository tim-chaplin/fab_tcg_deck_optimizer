package card

// Role is what a card did on a given turn cycle.
type Role uint8

const (
	Pitch Role = iota
	Attack
	Defend
	Held
	// Arsenal marks the card placed into the arsenal slot at end of turn.
	Arsenal
)

// String returns a human-readable role name.
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

// CardAssignment is a single card + the role it took this turn. Hand cards produce one
// per card; an arsenal-in card contributes one with FromArsenal set so a turn fits in
// one slice.
type CardAssignment struct {
	Card        Card
	Role        Role
	FromArsenal bool
}
