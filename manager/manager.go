package manager

import (
	"fmt"
	"sync"

	"github.com/dustinkirkland/golang-petname"
	"github.com/gwwfps/assembly-lines/game"
	"github.com/json-iterator/go"
	"github.com/labstack/echo"
	"gopkg.in/olahol/melody.v1"
)

type GameManager struct {
	m                     *melody.Melody
	logger                echo.Logger
	activeGamesById       map[string]*game.Game
	activeGamesByPlayerId map[string]*game.Game
	mutex                 *sync.Mutex
}

func NewGameManager(m *melody.Melody, logger echo.Logger) *GameManager {
	return &GameManager{
		m:                     m,
		logger:                logger,
		activeGamesById:       map[string]*game.Game{},
		activeGamesByPlayerId: map[string]*game.Game{},
		mutex:                 &sync.Mutex{},
	}
}

func (gm *GameManager) findGame(c MessageContext) *game.Game {
	return gm.activeGamesByPlayerId[c.PlayerId]
}

func (gm *GameManager) broadcastData(val interface{}, filter PlayerFilter) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("non-error panic: %v", r)
			}
			gm.logger.Error("recovered panic in broadcast", err)
		}
	}()

	data, err := jsoniter.ConfigFastest.Marshal(val)
	if err != nil {
		gm.logger.Error("cannot serialize data", err)
	}

	err = gm.m.BroadcastFilter(data, func(s *melody.Session) bool {
		id := s.MustGet(PlayerIdContextKey).(string)
		return filter(id)
	})
	if err != nil {
		gm.logger.Error("cannot write to WebSocket for broadcast", err)
	}
}

func (gm *GameManager) broadcastState(g *game.Game) {
	gm.broadcastData(g, g.IsPlayerJoined)
}

func (gm *GameManager) getLobbies() map[string][]string {
	lobbies := map[string][]string{}
	for id, g := range gm.activeGamesById {
		if g.Phase == game.GamePhaseLobby {
			lobbies[id] = g.GetPlayerNames()
		}
	}
	return lobbies
}

func (gm *GameManager) broadcastLobbies() {
	lobbies := gm.getLobbies()
	gm.broadcastData(lobbies, func(id string) bool {
		_, exists := gm.activeGamesByPlayerId[id]
		return !exists
	})
}

func (gm *GameManager) generateLobbyId() (string, error) {
	i := 0
	for {
		id := petname.Generate(3, "-")
		if _, exists := gm.activeGamesById[id]; !exists {
			return id, nil
		}
		i++
		if i > 10 {
			break
		}
	}
	return "", fmt.Errorf("cannot generate game id")
}

func (gm *GameManager) StartLobby(c MessageContext, args struct {
	Name      string
	SheetName string
}) error {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if gm.findGame(c) != nil {
		return fmt.Errorf("cannot start game when already part of lobby or game")
	}

	id, err := gm.generateLobbyId()
	if err != nil {
		return err
	}
	g := game.NewStandardGame()
	err = g.AddPlayer(c.PlayerId, args.Name, args.SheetName)
	if err != nil {
		return err
	}
	gm.activeGamesById[id] = g
	gm.activeGamesByPlayerId[c.PlayerId] = g

	gm.broadcastLobbies()
	gm.broadcastState(g)

	return nil
}
