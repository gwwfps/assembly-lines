package cards

import "math/rand"

type DeckDefinition struct {
	Actions       map[Action]int
	Machines      map[Machine]int
	MaxMachine    Machine
	SelectionSize int
}

func (d DeckDefinition) NewStack() []Card {
	var actions []Action
	for action, n := range d.Actions {
		for i := 0; i < n; i++ {
			actions = append(actions, action)
		}
	}
	var machines []Machine
	for machine, n := range d.Machines {
		for i := 0; i < n; i++ {
			machines = append(machines, machine)
		}
	}

	if len(machines) != len(actions) {
		panic("machine and action size mismatch in deck definition")
	}

	rand.Shuffle(len(machines), func(i, j int) {
		machines[i], machines[j] = machines[j], machines[i]
	})

	cards := make([]Card, len(machines))
	for i, machine := range machines {
		action := actions[i]
		cards[i] = Card{
			Machine: machine,
			Action:  action,
		}
	}

	return cards
}

var standardDeck = DeckDefinition{
	Actions: map[Action]int{
		ActionQC:          9,
		ActionCustomOrder: 9,
		ActionAddOn:       9,
		ActionRobot:       18,
		ActionTraining:    18,
		ActionBarrier:     18,
	},
	Machines: map[Machine]int{
		1:  3,
		2:  3,
		3:  4,
		4:  5,
		5:  6,
		6:  7,
		7:  8,
		8:  9,
		9:  8,
		10: 7,
		11: 6,
		12: 5,
		13: 4,
		14: 3,
		15: 3,
	},
	SelectionSize: 3,
}
