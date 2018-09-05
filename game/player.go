package game

import (
	"github.com/gwwfps/assembly-lines/game/sheet"
)

type PlayerPhase int

const (
	PhaseLobby PlayerPhase = iota
	PhaseSelecting
	PhaseCompletingObjectives
	PhaseWaitingForOthers
	PhaseGameEnd
)

type Player struct {
	Name  string
	Phase PlayerPhase
	Sheet *sheet.Sheet
}

func NewPlayer(name string, sheetName string) *Player {
	return &Player{
		Name:  name,
		Phase: PhaseLobby,
		Sheet: sheet.NewSheet(sheetName),
	}
}
