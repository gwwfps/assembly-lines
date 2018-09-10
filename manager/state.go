package manager

import (
	"github.com/gwwfps/assembly-lines/game"
)

type PlayerStatus int

const (
	statusLobby PlayerStatus = iota
	statusInGame
)

type PlayerState struct {
	Status PlayerStatus
	State  interface{}
	Init   *InitState
}

type InitState struct {
	Name string
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
		state.Init = &InitState{
			Name: g.GetPlayerById(c.PlayerId).Name,
		}
	} else {
		state = stateForLobby(gm.getLobbies())
	}

	return gm.reply(c, state)
}
