package manager

import (
	"github.com/gwwfps/assembly-lines/game"
	"github.com/json-iterator/go"
)

type PlayerStatus int

const (
	statusLobby PlayerStatus = iota
	statusInGame
)

type PlayerState struct {
	Status PlayerStatus
	State  interface{}
}

func stateForLobby(lobbies Lobbies) *PlayerState {
	return &PlayerState{
		Status: statusLobby,
		State:  lobbies,
	}
}

func stateForGame(g *game.Game) *PlayerState {
	return &PlayerState{
		Status: statusInGame,
		State:  g,
	}
}

func (gm *GameManager) FetchState(c MessageContext) error {
	g, inGame := gm.activeGamesByPlayerId[c.PlayerId]
	var state *PlayerState
	if inGame {
		state = stateForGame(g)
	} else {
		state = stateForLobby(gm.getLobbies())
	}

	data, err := jsoniter.ConfigFastest.Marshal(state)
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
