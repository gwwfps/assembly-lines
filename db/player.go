package db

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

func getPlayerKey(id string) string {
	return fmt.Sprintf("player-%s", id)
}

func (db *DB) GetPlayerGameId(playerId string) (string, error) {
	gameId, err := db.rc.Get(getPlayerKey(playerId)).Result()
	return gameId, errors.Wrap(err, "cannot retrieve gameId from db")
}

func (db *DB) JoinGame(playerId string, gameId string) error {
	playerKey := getPlayerKey(playerId)
	err := db.rc.Watch(func(tx *redis.Tx) error {
		currGameId, err := tx.Get(playerKey).Result()
		if err != nil {
			return err
		}
		if currGameId != "" {
			return fmt.Errorf("player %s is already in %s", playerId, currGameId)
		}

		_, err = tx.Set(playerId, gameId, 0).Result()
		return err
	}, playerKey)
	return errors.Wrapf(err, "cannot join game %s -> %s", playerId, gameId)
}

func (db *DB) UpdatePlayerName() {

}
