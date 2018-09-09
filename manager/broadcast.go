package manager

import (
	"fmt"

	"github.com/gwwfps/assembly-lines/game"
	"github.com/json-iterator/go"
	"gopkg.in/olahol/melody.v1"
)

const broadcastLobbies = "lobbies"

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

func (gm *GameManager) scheduleBroadcastState(g *game.Game) {
	go func() {
		gm.broadcastChan <- g.Id
	}()
}

func (gm *GameManager) broadcastState(g *game.Game) {
	gm.broadcastData(stateForGame(g), g.IsPlayerJoined)
}

func (gm *GameManager) scheduleBroadcastLobbies() {
	go func() {
		gm.broadcastChan <- broadcastLobbies
	}()
}

func (gm *GameManager) broadcastLobbies() {
	lobbies := gm.getLobbies()
	gm.broadcastData(stateForLobby(lobbies), func(id string) bool {
		_, exists := gm.activeGamesByPlayerId[id]
		return !exists
	})
}
