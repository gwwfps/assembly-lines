package sheet

const slotUnfilled int = -1
const slotExpandLeft int = -2
const slotExpandRight int = -3
const slotPackaging int = -4

type Sheet struct {
	Name         string
	Rows         [3]*Row
	Objectives   [3]int
	QCStations   int
	CustomOrders int
	Trainings    [6]int
	AddOns       int
	Skips        int
}

type rowDefinition struct {
	Size       int
	RobotLimit int
	QCSpots    [3]int
}

type Row struct {
	Slots  []int
	Robots int
}

var row1Definition = rowDefinition{
	Size:       10,
	RobotLimit: 3,
	QCSpots:    [3]int{2, 6, 7},
}

var row2Definition = rowDefinition{
	Size:       11,
	RobotLimit: 4,
	QCSpots:    [3]int{0, 3, 7},
}

var row3Definition = rowDefinition{
	Size:       12,
	RobotLimit: 4,
	QCSpots:    [3]int{1, 6, 10},
}

func NewSheet(name string) *Sheet {
	return &Sheet{
		Name: name,
		Rows: [3]*Row{
			newRow(row1Definition),
			newRow(row2Definition),
			newRow(row3Definition),
		},
	}
}

func newRow(def rowDefinition) *Row {
	slots := make([]int, def.Size)
	for i := range slots {
		slots[i] = slotUnfilled
	}
	return &Row{
		Slots: slots,
	}
}
