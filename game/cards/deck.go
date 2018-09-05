package cards

import "math/rand"

type Deck struct {
	definition     DeckDefinition
	allCards       []Card
	remainingCards []Card
	PreviousCards  []Card
	ActiveCards    []Card
	Fresh          bool
}

func NewStandardDeck() *Deck {
	return NewDeck(standardDeck)
}

func NewDeck(def DeckDefinition) *Deck {
	stack := def.NewStack()
	remaining := make([]Card, len(stack))
	copy(remaining, stack)

	selectionSize := def.SelectionSize

	deck := &Deck{
		definition:     def,
		allCards:       stack,
		remainingCards: remaining,
		PreviousCards:  make([]Card, selectionSize),
		ActiveCards:    make([]Card, selectionSize),
	}

	deck.shuffle()
	deck.DrawNext()
	deck.Fresh = true

	return deck
}

func (d *Deck) DrawNext() {
	d.Fresh = false

	drawSize := len(d.ActiveCards)
	if len(d.remainingCards) < drawSize {
		d.remake()
		d.shuffle()
	}

	d.PreviousCards, d.ActiveCards, d.remainingCards = d.ActiveCards, d.remainingCards[:drawSize], d.remainingCards[drawSize:]
}

func (d *Deck) shuffle() {
	rand.Shuffle(len(d.remainingCards), func(i, j int) {
		d.remainingCards[i], d.remainingCards[j] = d.remainingCards[j], d.remainingCards[i]
	})
	d.Fresh = true
}

func (d *Deck) remake() {
	drawSize := len(d.ActiveCards)
	stack := d.allCards
	remaining := make([]Card, len(stack)-drawSize)

	i, j := 0, 0
	for _, card := range stack {
		active := false
		if j < drawSize {
			for _, activeCard := range d.ActiveCards {
				if activeCard.Equals(card) {
					active = true
					j++
					break
				}
			}
		}

		if !active {
			remaining[i] = card
			i++
		}
	}

	if j != drawSize {
		panic("invalid active cards, cannot remake deck")
	}

	d.remainingCards = remaining
}
