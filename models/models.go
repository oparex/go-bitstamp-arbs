package models

type PricePoint struct {
	Price float64
	Amount float64
}

type MarketPoint struct {
	Bid *PricePoint
	Ask *PricePoint
}

type PathNode struct {
	Pair string
	Side string
}