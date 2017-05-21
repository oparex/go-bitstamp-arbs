package main

import (
	"log"
	"github.com/toorop/go-pusher"
	"github.com/ajph/bitstamp-go"
	"arbot/config"
	"arbot/parser"
)

func main() {

	// init pusher client
	pusherClient, err := pusher.NewClient(config.WEBSOCKET_APP_KEY)
	if err != nil {
		log.Fatalln(err)
	}
	// subscribe to all channels
	for _, c := range config.Channels {
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
	httpBuf := make(map[string]interface{})
	socketBuf := make(map[string][][]string)

	// get initial best prices via http for all pairs
	bestPrices, err := parser.InitBestPrices(config.Channels, httpBuf)
	if err != nil {
		log.Fatalln(err)
	}

	bitstamp.SetAuth(config.CLIENTID, config.KEY, config.SECRET)

	// start listening
	for {
		dataEvt := <-dataChannelOrderBook
		pp, err := parser.ParseSocketPricePoint(dataEvt.Data, socketBuf)
		if err != nil {
			log.Println(err)
			continue
		}
		bestPrices[dataEvt.Channel] = pp

		bestPrices.CheckPaths(config.Paths)
	}
}
