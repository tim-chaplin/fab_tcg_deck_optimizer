package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// defaultHeroTypes is the TypeSet defaultHero advertises.
var defaultHeroTypes = card.NewTypeSet(card.TypeHero, card.TypeYoung)

// defaultHero is the no-ability fallback hero used by GameStateBuilder when SetHero isn't
// called. 20 health, 4 Intelligence, no triggered ability, no Opt preference. Lives in the
// gameengine package (rather than internal/hero/heroes) so the builder can default to it
// without introducing a gameengine → heroes import cycle with heroes_test.
type defaultHero struct{}

func (defaultHero) ID() ids.HeroID                                  { return ids.DefaultHeroID }
func (defaultHero) Name() string                                    { return "Default" }
func (defaultHero) Health() int                                     { return 20 }
func (defaultHero) Intelligence() int                               { return 4 }
func (defaultHero) Types() card.TypeSet                             { return defaultHeroTypes }
func (defaultHero) Class() card.CardType                            { return card.CardType(0) }
func (defaultHero) Opt(cards []card.Card) (top, bottom []card.Card) { return nil, cards }

// No triggered ability — TriggerType == 0 makes FireTriggers skip the hero.
func (defaultHero) TriggerType() triggertype.Type     { return 0 }
func (defaultHero) OncePerTurn() bool                 { return false }
func (defaultHero) FiredThisTurn() bool               { return false }
func (defaultHero) SetFiredThisTurn(bool)             {}
func (defaultHero) Matches(card.TypeSet) bool         { return false }
func (defaultHero) Fire(card.GameEngine, card.Logger) {}
