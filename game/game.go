package game

import (
	"fmt"
	"sync"
	"time"

	"github.com/gwwfps/assembly-lines/game/cards"
)

type GamePhase int

const (
	GamePhaseLobby GamePhase = iota
	GamePhaseInProgress
	GamePhaseComplete
)

type Game struct {
	Id               string
	CreationTime     time.Time
	Deck             *cards.Deck
	Phase            GamePhase
	Players          map[string]*Player
	playerNames      map[string]string
	globalStateMutex *sync.Mutex
}

func NewStandardGame(id string) *Game {
	return &Game{
		Id:               id,
		CreationTime:     time.Now(),
		Phase:            GamePhaseLobby,
		Players:          map[string]*Player{},
		Deck:             cards.NewStandardDeck(),
		playerNames:      map[string]string{},
		globalStateMutex: &sync.Mutex{},
	}
}

func (g *Game) GetPlayerById(id string) *Player {
	name, exists := g.playerNames[id]
	if exists {
		return g.Players[name]
	}
	return nil
}

func (g *Game) IsPlayerJoined(id string) bool {
	_, exists := g.playerNames[id]
	return exists
}

func (g *Game) GetPlayerNames() []string {
	var names []string
	for _, name := range g.playerNames {
		names = append(names, name)
	}
	return names
}

func (g *Game) AddPlayer(id string, name string, sheetName string) error {
	g.globalStateMutex.Lock()
	defer g.globalStateMutex.Unlock()

	if g.Phase != GamePhaseLobby {
		return fmt.Errorf("game already started")
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
	g.globalStateMutex.Lock()
	defer g.globalStateMutex.Unlock()

	if g.Phase != GamePhaseLobby {
		return fmt.Errorf("game already started")
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
	g.globalStateMutex.Lock()
	defer g.globalStateMutex.Unlock()

	if g.Phase != GamePhaseLobby {
		return fmt.Errorf("game already started")
	}

	g.Phase = GamePhaseInProgress
	return nil
}
