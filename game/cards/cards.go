package cards

type Action int
type Machine int

const (
	ActionQC Action = iota
	ActionCustomOrder
	ActionAddOn
	ActionRobot
	ActionTraining
	ActionBarrier
)

type Card struct {
	Action  Action
	Machine Machine
}

func (c Card) Equals(other Card) bool {
	return c.Action == other.Action && c.Machine == other.Machine
}
