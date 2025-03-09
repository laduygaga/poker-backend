package utils

import (
	"math/rand"
	"time"
)

type Deck struct {
	Cards []string
}

func NewDeck() *Deck {
	cards := make([]string, 52)
	suits := []string{"♠", "♥", "♦", "♣"}
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	idx := 0
	for _, suit := range suits {
		for _, rank := range ranks {
			cards[idx] = rank + suit
			idx++
		}
	}
	return &Deck{Cards: cards}
}

func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func (d *Deck) Deal(n int) []string {
	cards := d.Cards[:n]
	d.Cards = d.Cards[n:]
	return cards
}
