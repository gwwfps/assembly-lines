package manager

import (
	"fmt"
	"sync"
	"time"

	"github.com/gwwfps/assembly-lines/game"
	"github.com/json-iterator/go"
	"github.com/labstack/echo"
	"gopkg.in/olahol/melody.v1"
)

var unexpectedError = fmt.Errorf("unexpected error")

type GameManager struct {
	m                     *melody.Melody
	logger                echo.Logger
	activeGamesById       map[string]*game.Game
	activeGamesByPlayerId map[string]*game.Game
	mutex                 *sync.Mutex
	broadcastChan         chan string
}

func NewGameManager(m *melody.Melody, logger echo.Logger) *GameManager {
	return &GameManager{
		m:                     m,
		logger:                logger,
		activeGamesById:       map[string]*game.Game{},
		activeGamesByPlayerId: map[string]*game.Game{},
		mutex:                 &sync.Mutex{},
		broadcastChan:         make(chan string, 2048),
	}
}

func Watch(gm *GameManager) {
	flagged := make(map[string]bool)
	broadcastTicker := time.NewTicker(50 * time.Millisecond)
	gcTicker := time.NewTicker(100 * time.Second)
	for {
		select {
		case name := <-gm.broadcastChan:
			flagged[name] = true
		case <-broadcastTicker.C:
			for id, broadcast := range flagged {
				if broadcast {
					flagged[id] = false
					if id == broadcastLobbies {
						go gm.broadcastLobbies()
					} else {
						go gm.broadcastState(gm.activeGamesById[id])
					}
				}
			}
		case <-gcTicker.C:
			go gm.garbageCollectLobbies()
		}
	}
}

func (gm *GameManager) findGame(c MessageContext) *game.Game {
	return gm.activeGamesByPlayerId[c.PlayerId]
}

func (gm *GameManager) reply(c MessageContext, payload interface{}) error {
	data, err := jsoniter.ConfigFastest.Marshal(payload)
	if err != nil {
		gm.logger.Error("cannot serialize data", err)
		return unexpectedError
	}

	err = c.Session.Write(data)
	if err != nil {
		gm.logger.Error("cannot write to WebSocket", err)
		return unexpectedError
	}

	return nil
}
