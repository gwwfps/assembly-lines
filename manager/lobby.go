package manager

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustinkirkland/golang-petname"
	"github.com/gwwfps/assembly-lines/game"
)

const activeGameLimit = 20

type Lobbies map[string][]string

type JoinLobbyArgs struct {
	Id        string
	Name      string
	SheetName string
}

func (args JoinLobbyArgs) validate() (string, string, error) {
	name := strings.TrimSpace(args.Name)
	sheetName := strings.TrimSpace(args.SheetName)
	if name == "" {
		return "", "", fmt.Errorf("name is required")
	}
	if sheetName == "" {
		return "", "", fmt.Errorf("sheet name is required")
	}
	return name, sheetName, nil
}

func (gm *GameManager) getLobbies() Lobbies {
	lobbies := map[string][]string{}
	for id, g := range gm.activeGamesById {
		if g.Phase == game.GamePhaseLobby {
			lobbies[id] = g.GetPlayerNames()
		}
	}
	return lobbies
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

func (gm *GameManager) StartLobby(c MessageContext, args JoinLobbyArgs) error {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	name, sheetName, err := args.validate()
	if err != nil {
		return err
	}

	if len(gm.activeGamesById) >= activeGameLimit {
		return fmt.Errorf("maximum number of active games/lobbies reached")
	}

	if gm.findGame(c) != nil {
		return fmt.Errorf("cannot start game when already part of lobby or game")
	}

	id, err := gm.generateLobbyId()
	if err != nil {
		return err
	}
	g := game.NewStandardGame(id)
	err = g.AddPlayer(c.PlayerId, name, sheetName)
	if err != nil {
		return err
	}
	gm.activeGamesById[id] = g
	gm.activeGamesByPlayerId[c.PlayerId] = g

	gm.scheduleBroadcastLobbies()
	gm.scheduleBroadcastState(g)

	gm.logger.Infof("%s started game %s as %s", c.PlayerId, g.Id, name)

	return nil
}

func (gm *GameManager) JoinLobby(c MessageContext, args JoinLobbyArgs) error {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	name, sheetName, err := args.validate()
	if err != nil {
		return err
	}

	if gm.findGame(c) != nil {
		return fmt.Errorf("cannot join game when already part of lobby or game")
	}

	g := gm.activeGamesById[args.Id]
	if g == nil {
		return fmt.Errorf("game does not exist")
	}

	err = g.AddPlayer(c.PlayerId, name, sheetName)
	if err != nil {
		return err
	}
	gm.activeGamesByPlayerId[c.PlayerId] = g

	gm.scheduleBroadcastLobbies()
	gm.scheduleBroadcastState(g)

	gm.logger.Infof("%s joined %s as %s", c.PlayerId, g.Id, name)

	return nil
}

func (gm *GameManager) garbageCollectLobbies() {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	var gameIds, playerIds []string
	for _, g := range gm.activeGamesById {
		if (g.Phase == game.GamePhaseLobby && g.CreationTime.Add(48*time.Hour).Before(time.Now())) || g.CreationTime.Add(7*24*time.Hour).Before(time.Now()) {
			gameIds = append(gameIds, g.Id)
			for playerId, pg := range gm.activeGamesByPlayerId {
				if pg == g {
					playerIds = append(playerIds, playerId)
					break
				}
			}
		}
	}

	for _, gameId := range gameIds {
		gm.logger.Infof("garbage-collecting game %s", gameId)
		delete(gm.activeGamesById, gameId)
	}

	for _, playerId := range playerIds {
		gm.logger.Infof("garbage-collecting game for player %s", playerId)
		delete(gm.activeGamesByPlayerId, playerId)
	}
}
