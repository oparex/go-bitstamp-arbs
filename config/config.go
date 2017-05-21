package config

import (
	"arbot/models"
)

const (
	BITSTAMP_HTTP_URL = "https://www.bitstamp.net/api/v2/order_book/"
	WEBSOCKET_APP_KEY = "de504dc5763aeef9ff52" // bitstamp
	ARB_TRADE_THRESHOLD = 1.01 // arbitrage minimal relative profit for entering trades
)

// pair/channel map of currency pairs to track
var Channels = map[string]string{
	"btcusd": "order_book",
	"btceur": "order_book_btceur",
	"eurusd": "order_book_eurusd",
	"xrpusd": "order_book_xrpusd",
	"xrpeur": "order_book_xrpeur",
	"xrpbtc": "order_book_xrpbtc",
}

// map of paths to check for arbitrage - name: []PathNode
var Paths = map[string][]models.PathNode{
	"ubxu": {{"btcusd", "ask"}, {"xrpbtc", "ask"}, {"xrpusd", "bid"}},
	//"ubxeu": {[]PathNode{{"btcusd", "ask"}, {"xrpbtc", "ask"}, {"xrpeur", "bid"}, {"eurusd", "bid"}}},
	"ubeu": {{"btcusd", "ask"}, {"btceur", "bid"}, {"eurusd", "bid"}},
	//"ubexu": {[]PathNode{{"btcusd", "ask"}, {"btceur", "bid"}, {"xrpeur", "ask"}, {"xrpusd", "bid"}}},
	"uxbu": {{"xrpusd", "ask"}, {"xrpbtc", "bid"}, {"btcusd", "bid"}},
	//"uxbeu": {[]PathNode{{"xrpusd", "ask"}, {"xrpbtc", "bid"}, {"btceur", "bid"}, {"eurusd", "bid"}}},
	"uxeu": {{"xrpusd", "ask"}, {"xrpeur", "bid"}, {"eurusd", "bid"}},
	//"uxebu": {[]PathNode{{"xrpusd", "ask"}, {"xrpeur", "bid"}, {"btceur", "ask"}, {"btcusd", "bid"}}},
	"uebu": {{"eurusd", "ask"}, {"btceur", "ask"}, {"btcusd", "bid"}},
	//"uebxu": {[]PathNode{{"eurusd", "ask"}, {"btceur", "ask"}, {"xrpbtc", "ask"}, {"xrpusd", "bid"}}},
	"uexu": {{"eurusd", "ask"}, {"xrpeur", "ask"}, {"xrpusd", "bid"}},
	//"uexbu": {[]PathNode{{"eurusd", "ask"}, {"xrpeur", "ask"}, {"xrpbtc", "bid"}, {"btcusd", "bid"}}},
}