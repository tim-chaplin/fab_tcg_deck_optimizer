package trigger

// Card is the minimal source-card surface a card-backed trigger needs (display name). The
// richer v2/card.Card satisfies it structurally so consumers can pass card values through
// without conversion.
type Card interface {
	DisplayName() string
}
