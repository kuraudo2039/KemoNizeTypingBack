package cardObj

import (
	"math/rand"
)

type Deck struct {
	Cards     []Card `json:"-"`
	IsShuffle bool   `json:"-"`
	Counts    int    `json:"counts"`
}

type Card struct {
	ID  int `json:"id"`
	Num int `json:"num"`
}

var deckTemp = []Card{
	Card{0, 0},
	Card{0, 0},
	Card{0, 0},
	Card{1, 1},
	Card{1, 1},
	Card{1, 1},
	Card{1, 1},
	Card{2, 2},
	Card{2, 2},
	Card{2, 2},
	Card{2, 2},
	Card{3, 3},
	Card{3, 3},
	Card{3, 3},
	Card{3, 3},
	Card{4, 4},
	Card{4, 4},
	Card{4, 4},
	Card{4, 4},
	Card{5, 5},
	Card{5, 5},
	Card{5, 5},
	Card{5, 5},
	Card{6, 10},
	Card{6, 10},
	Card{6, 10},
	Card{7, 15},
	Card{7, 15},
	Card{8, 20},
	Card{9, -5},
	Card{9, -5},
	Card{10, -10},
	// ゲームフロー影響カード
	Card{100, 0}, // 夜カード：次ターン山札シャッフル
	Card{101, 0}, // 洞穴カード：計算時に山札からカードを引く
	// 計算時影響カード
	Card{200, 0}, // 酋長カード：すべての基本カードを2倍
	Card{201, 0}, // 狐カード：最大のコヨーテカードを0に
}

func CreateDeck() Deck {
	deck := Deck{initCards(), false, len(deckTemp)}
	return deck
}

func (deck *Deck) DrawCard() Card {
	// 山札がなった場合シャッフル
	if len(deck.Cards) == 0 {
		deck.Shuffle()
	}

	index := rand.Intn(len(deck.Cards))
	drawenCard := deck.Cards[index]
	deck.Cards = append(deck.Cards[:index], deck.Cards[index+1:]...)

	// 夜カードを引いた場合シャッフルフラグtrue
	if drawenCard.ID == 100 {
		deck.IsShuffle = true
	}

	deck.Counts--
	return drawenCard
}

func (deck *Deck) Shuffle() {
	deck.Cards = initCards()
	deck.Counts = len(deckTemp)
	deck.IsShuffle = false
}

/* local func*/

func initCards() []Card {
	cards := make([]Card, 0)
	copy(cards, deckTemp)
	return cards
}
