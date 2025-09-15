package models

import "fmt"

// CardSuit 花色
type CardSuit int

const (
	SuitSpades   CardSuit = 1 // 黑桃
	SuitHearts   CardSuit = 2 // 红桃
	SuitDiamonds CardSuit = 3 // 方块
	SuitClubs    CardSuit = 4 // 梅花
	SuitJoker    CardSuit = 5 // 王
)

var SuitNames = map[CardSuit]string{
	SuitSpades:   "♠",
	SuitHearts:   "♥",
	SuitDiamonds: "♦",
	SuitClubs:    "♣",
	SuitJoker:    "王",
}

// CardValue 牌值
type CardValue int

const (
	Value3          CardValue = 3  // 3
	Value4          CardValue = 4  // 4
	Value5          CardValue = 5  // 5
	Value6          CardValue = 6  // 6
	Value7          CardValue = 7  // 7
	Value8          CardValue = 8  // 8
	Value9          CardValue = 9  // 9
	Value10         CardValue = 10 // 10
	ValueJack       CardValue = 11 // J
	ValueQueen      CardValue = 12 // Q
	ValueKing       CardValue = 13 // K
	ValueAce        CardValue = 14 // A
	Value2          CardValue = 15 // 2
	ValueSmallJoker CardValue = 16 // 小王
	ValueBigJoker   CardValue = 17 // 大王
)

var ValueNames = map[CardValue]string{
	Value3:          "3",
	Value4:          "4",
	Value5:          "5",
	Value6:          "6",
	Value7:          "7",
	Value8:          "8",
	Value9:          "9",
	Value10:         "10",
	ValueJack:       "J",
	ValueQueen:      "Q",
	ValueKing:       "K",
	ValueAce:        "A",
	Value2:          "2",
	ValueSmallJoker: "小王",
	ValueBigJoker:   "大王",
}

// Card 扑克牌
type Card struct {
	Suit  CardSuit  `json:"suit"`  // 花色
	Value CardValue `json:"value"` // 牌值
}

// String 返回牌的字符串表示
func (c Card) String() string {
	if c.Suit == SuitJoker {
		return ValueNames[c.Value]
	}
	return fmt.Sprintf("%s%s", SuitNames[c.Suit], ValueNames[c.Value])
}

// NewCard 创建新牌
func NewCard(suit CardSuit, value CardValue) Card {
	return Card{Suit: suit, Value: value}
}

// NewDeck 创建一副54张的扑克牌
func NewDeck() []Card {
	var deck []Card

	// 添加普通牌（3-A, 2）
	suits := []CardSuit{SuitSpades, SuitHearts, SuitDiamonds, SuitClubs}
	values := []CardValue{
		Value3, Value4, Value5, Value6, Value7, Value8, Value9, Value10,
		ValueJack, ValueQueen, ValueKing, ValueAce, Value2,
	}

	for _, suit := range suits {
		for _, value := range values {
			deck = append(deck, NewCard(suit, value))
		}
	}

	// 添加大小王
	deck = append(deck, NewCard(SuitJoker, ValueSmallJoker))
	deck = append(deck, NewCard(SuitJoker, ValueBigJoker))

	return deck
}

// IsJoker 判断是否为王牌
func (c Card) IsJoker() bool {
	return c.Suit == SuitJoker
}

// GetWeight 获取牌的权重（用于排序和比较）
func (c Card) GetWeight() int {
	return int(c.Value)
}

// Compare 比较两张牌的大小，返回值：-1表示小于，0表示等于，1表示大于
func (c Card) Compare(other Card) int {
	if c.Value < other.Value {
		return -1
	} else if c.Value > other.Value {
		return 1
	}
	return 0
}