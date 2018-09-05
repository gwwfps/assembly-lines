package server

import (
	"fmt"

	"github.com/gwwfps/assembly-lines/game"
)

type GameManager struct {
	games           []*game.Game
	gamesByPlayerId map[string]*game.Game
}

func (gm *GameManager) ExampleAction(playerId string, args struct {
	Question int
}) error {
	return fmt.Errorf("not implemented")
}
