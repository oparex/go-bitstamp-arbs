# Check for arbitrage on Bitstamp

This program depends on github.com/toorop/go-pusher. Install it first with
```
go get github.com/toorop/go-pusher
```
Then run the arbot program with
```
go run arbot.go
```
or build/install the program.

It will connect to Bitstamp's socket and listen to order_book(_currencypair) channels (best 100 bids/asks). 
It has some predefined arbitrage paths (hops from one currency to another that end in the starting one) that it checks. 
When a new message is received, all paths are checked for arbitrage and updated result map is printed to stdout.