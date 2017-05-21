package arber

import (
	"math"
	"github.com/ajph/bitstamp-go"
	"arbot/models"
	"arbot/config"
)

type BestPrices map[string]*models.MarketPoint

// check all paths for arbitrage
func (bp *BestPrices) CheckPaths(paths map[string][]models.PathNode) {
	for _, path := range paths {
		if bp.CheckPath(path) > config.ARB_TRADE_THRESHOLD {
			go bp.ActOnArb(path)
		}
	}
}

// check path for arbitrage
func (bp BestPrices) CheckPath(path []models.PathNode) float64 {
	outcome := 1.0
	bpKey := ""
	for _, p := range path {
		if p.Pair == "btcusd" {
			bpKey = "order_book"
		} else {
			bpKey = "order_book_" + p.Pair
		}
		if p.Side == "ask" {
			outcome /= bp[bpKey].Ask.Price
		}
		if p.Side == "bid" {
			outcome *= bp[bpKey].Bid.Price
		}
	}
	return outcome
}

func (bp BestPrices) ActOnArb(path []models.PathNode) {
	amounts := [3]float64{}
	if path[1].Side == "ask" {
		if math.Min(bp[path[1].Pair].Ask.Amount, bp[path[2].Pair].Bid.Amount) *
			bp[path[1].Pair].Ask.Price > bp[path[0].Pair].Ask.Amount {
			amounts[0] = bp[path[0].Pair].Ask.Amount
			amounts[1] = bp[path[0].Pair].Ask.Amount / bp[path[1].Pair].Ask.Price
			amounts[2] = amounts[1]
		} else {
			amounts[1] = math.Min(bp[path[1].Pair].Ask.Amount, bp[path[2].Pair].Bid.Amount)
			amounts[2] = amounts[1]
			amounts[0] = amounts[1] * bp[path[1].Pair].Ask.Price
		}
		bitstamp.BuyLimitOrder(path[0].Pair, amounts[0], bp[path[0].Pair].Ask.Price)
		bitstamp.BuyLimitOrder(path[1].Pair, amounts[1], bp[path[1].Pair].Ask.Price)
		bitstamp.SellLimitOrder(path[2].Pair, amounts[2], bp[path[2].Pair].Bid.Price)
	} else {
		if math.Min(bp[path[0].Pair].Ask.Amount, bp[path[1].Pair].Bid.Amount) *
			bp[path[1].Pair].Bid.Price > bp[path[2].Pair].Bid.Amount {
			amounts[2] = bp[path[2].Pair].Bid.Amount
			amounts[1] = bp[path[2].Pair].Bid.Amount * bp[path[1].Pair].Bid.Price
			amounts[0] = amounts[1]
		} else {
			amounts[0] = math.Min(bp[path[0].Pair].Ask.Amount, bp[path[1].Pair].Bid.Amount)
			amounts[1] = amounts[0]
			amounts[2] = amounts[0] * bp[path[1].Pair].Bid.Price
		}
		bitstamp.BuyLimitOrder(path[0].Pair, amounts[0], bp[path[0].Pair].Ask.Price)
		bitstamp.SellLimitOrder(path[1].Pair, amounts[1], bp[path[1].Pair].Bid.Price)
		bitstamp.SellLimitOrder(path[2].Pair, amounts[2], bp[path[2].Pair].Bid.Price)
	}
}