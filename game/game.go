package game

import (
	"fmt"

	"github.com/gwwfps/assembly-lines/game/cards"
)

type Game struct {
	Deck        *cards.Deck
	InProgress  bool
	Players     map[string]*Player
	playerNames map[string]string
}

func NewStandardGame() *Game {
	return &Game{
		Players:     map[string]*Player{},
		Deck:        cards.NewStandardDeck(),
		playerNames: map[string]string{},
	}
}

func (g *Game) AddPlayer(id string, name string, sheetName string) error {
	if g.InProgress {
		return fmt.Errorf("game already in progress")
	}

	_, inGame := g.playerNames[id]
	if inGame {
		return fmt.Errorf("player already in game")
	}

	_, nameExists := g.Players[name]
	if nameExists {
		return fmt.Errorf("name already taken")
	}

	g.playerNames[id] = name
	g.Players[name] = NewPlayer(name, sheetName)
	return nil
}

func (g *Game) WithdrawPlayer(id string) error {
	if g.InProgress {
		return fmt.Errorf("game already in progress")
	}

	name, exists := g.playerNames[id]
	if !exists {
		return fmt.Errorf("player not in game")
	}

	delete(g.playerNames, id)
	delete(g.Players, name)
	return nil
}

func (g *Game) Start() error {
	if g.InProgress {
		return fmt.Errorf("game already in progress")
	}

	g.InProgress = true
	return nil
}
