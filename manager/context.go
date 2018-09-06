package manager

import (
	"gopkg.in/olahol/melody.v1"
)

const PlayerIdContextKey = "playerId"

type PlayerFilter func(playerId string) bool

type MessageContext struct {
	Session  *melody.Session
	PlayerId string
}

func SessionContext(playerId string) map[string]interface{} {
	return map[string]interface{}{
		PlayerIdContextKey: playerId,
	}
}
