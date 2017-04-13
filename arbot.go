package main

import (
	"github.com/toorop/go-pusher"
	"log"
	"encoding/json"
	"strconv"
	"errors"
	"net/http"
	"io/ioutil"
)

const (
	APP_KEY = "de504dc5763aeef9ff52" // bitstamp
)

type PricePoint struct {
	bid float64
	ask float64
}

type PathNode struct {
	pair string
	side string
}

type Path struct {
	path []PathNode
}

type Paths struct {
	paths map[string]*Path
}

// check path for arbitrage
func (pth *Path) checkPath(bestPrices map[string]*PricePoint) float64 {
	outcome := 1.0
	bpKey := ""
	for _, p := range pth.path {
		if p.pair == "btcusd" {
			bpKey = "order_book"
		} else {
			bpKey = "order_book_" + p.pair
		}
		if p.side == "ask" {
			outcome /= bestPrices[bpKey].ask
		}
		if p.side == "bid" {
			outcome *= bestPrices[bpKey].bid
		}
	}
	return outcome
}

// check all paths for arbitrage
func (pths *Paths) checkPaths(bestPrices map[string]*PricePoint) map[string]float64 {
	outcome := make(map[string]float64)
	for name, path := range pths.paths {
		outcome[name] = path.checkPath(bestPrices)
	}
	return outcome
}

// get PricePoint from a socket message
func parseSocketPricePoint(raw_message string, msg_buf map[string][][]string) (*PricePoint, error) {
	err := json.Unmarshal([]byte(raw_message), &msg_buf)
	if err != nil {
		return nil, errors.New("Could not unmarshal raw message into json: " + raw_message)
	}
	bid, err := strconv.ParseFloat(msg_buf["bids"][0][0], 64)
	if err != nil {
		return nil, errors.New("Could not convert best bid to float: " + msg_buf["bids"][0][0])
	}
	ask, err := strconv.ParseFloat(msg_buf["asks"][0][0], 64)
	if err != nil {
		return nil, errors.New("Could not convert best ask to float: " + msg_buf["asks"][0][0])
	}
	return &PricePoint{bid, ask}, nil
}

// get PricePoint from a http message
func parseHttpPricePoint(raw_message []byte, msg_buf map[string]interface{}) (*PricePoint, error) {
	err := json.Unmarshal(raw_message, &msg_buf)
	if err != nil {
		return nil, errors.New("Could not unmarshal raw message into json from http " + string(raw_message))
	}
	bid, err := strconv.ParseFloat(msg_buf["bids"].([]interface{})[0].([]interface{})[0].(string), 64)
	if err != nil {
		return nil, errors.New("Could not convert best bid to float: " + msg_buf["bids"].([]interface{})[0].([]interface{})[0].(string))
	}
	ask, err := strconv.ParseFloat(msg_buf["asks"].([]interface{})[0].([]interface{})[0].(string), 64)
	if err != nil {
		return nil, errors.New("Could not convert best ask to float: " + msg_buf["asks"].([]interface{})[0].([]interface{})[0].(string))
	}
	return &PricePoint{bid, ask}, nil
}

// get all initial prices for pairs/channels via http
func initBestPrices(channels map[string]string, msg_buf map[string]interface{}) (map[string]*PricePoint, error) {
	bestPrices := make(map[string]*PricePoint)

	for pair, channel := range channels {
		resp, err := http.Get("https://www.bitstamp.net/api/v2/order_book/" + pair)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		pp, err := parseHttpPricePoint(body, msg_buf)
		if err != nil {
			return nil, err
		}
		bestPrices[channel] = pp
	}
	return bestPrices, nil
}

// pair/channel map of currency pairs to track
var channels = map[string]string{
	"btcusd": "order_book",
	"btceur": "order_book_btceur",
	"eurusd": "order_book_eurusd",
	"xrpusd": "order_book_xrpusd",
	"xrpeur": "order_book_xrpeur",
	"xrpbtc": "order_book_xrpbtc",
}

// map of paths to check for arbitrage - name: []PathNode
var paths = Paths{map[string]*Path{
	"ubxu": {[]PathNode{{"btcusd", "ask"}, {"xrpbtc", "ask"}, {"xrpusd", "bid"}}},
	"ubxeu": {[]PathNode{{"btcusd", "ask"}, {"xrpbtc", "ask"}, {"xrpeur", "bid"}, {"eurusd", "bid"}}},
	"ubeu": {[]PathNode{{"btcusd", "ask"}, {"btceur", "bid"}, {"eurusd", "bid"}}},
	"ubexu": {[]PathNode{{"btcusd", "ask"}, {"btceur", "bid"}, {"xrpeur", "ask"}, {"xrpusd", "bid"}}},
	"uxbu": {[]PathNode{{"xrpusd", "ask"}, {"xrpbtc", "bid"}, {"btcusd", "bid"}}},
	"uxbeu": {[]PathNode{{"xrpusd", "ask"}, {"xrpbtc", "bid"}, {"btceur", "bid"}, {"eurusd", "bid"}}},
	"uxeu": {[]PathNode{{"xrpusd", "ask"}, {"xrpeur", "bid"}, {"eurusd", "bid"}}},
	"uxebu": {[]PathNode{{"xrpusd", "ask"}, {"xrpeur", "bid"}, {"btceur", "ask"}, {"btcusd", "bid"}}},
	"uebu": {[]PathNode{{"eurusd", "ask"}, {"btceur", "ask"}, {"btcusd", "bid"}}},
	"uebxu": {[]PathNode{{"eurusd", "ask"}, {"btceur", "ask"}, {"xrpbtc", "ask"}, {"xrpusd", "bid"}}},
	"uexu": {[]PathNode{{"eurusd", "ask"}, {"xrpeur", "ask"}, {"xrpusd", "bid"}}},
	"uexbu": {[]PathNode{{"eurusd", "ask"}, {"xrpeur", "ask"}, {"xrpbtc", "bid"}, {"btcusd", "bid"}}},
}}

func main() {

	// init pusher client
	pusherClient, err := pusher.NewClient(APP_KEY)
	if err != nil {
		log.Fatalln(err)
	}
	// subscribe to all channels
	for _, c := range channels {
		err = pusherClient.Subscribe(c)
		if err != nil {
			log.Fatalln("Subscription error : ", err, c)
		}
	}
	// bind events
	dataChannelOrderBook, err := pusherClient.Bind("data")
	if err != nil {
		log.Fatalln("Bind error: ", err)
	}

	// init buffers
	http_buf := make(map[string]interface{})
	socket_buf := make(map[string][][]string)

	// get initial best prices via http for all pairs
	bestPrices, err := initBestPrices(channels, http_buf)
	if err != nil {
		log.Fatalln(err)
	}

	// start listening
	for {
		dataEvt := <-dataChannelOrderBook
		pp, err := parseSocketPricePoint(dataEvt.Data, socket_buf)
		if err != nil {
			log.Println(err)
		}
		bestPrices[dataEvt.Channel] = pp
		// print out which channel emitted this message and what is the arb profit for each path
		log.Println(dataEvt.Channel, paths.checkPaths(bestPrices))
	}
}
