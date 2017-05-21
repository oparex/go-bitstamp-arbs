package parser

import (
	"encoding/json"
	"strconv"
	"io/ioutil"
	"errors"
	"arbot/models"
	"net/http"
	"arbot/config"
	"arbot/arber"
)

// get all initial prices for pairs/channels via http
func InitBestPrices(channels map[string]string, msgBuf map[string]interface{}) (arber.BestPrices, error) {
	bestPrices := make(map[string]*models.MarketPoint)

	for pair, channel := range channels {
		pp, err := GetBestPrice(pair, msgBuf)
		if err != nil {
			return nil, err
		}
		bestPrices[channel] = pp
	}
	return bestPrices, nil
}

func GetBestPrice(pair string, msgBuf map[string]interface{}) (*models.MarketPoint, error) {
	resp, err := http.Get(config.BITSTAMP_HTTP_URL + pair)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	pp, err := ParseHttpPricePoint(body, msgBuf)
	if err != nil {
		return nil, err
	}
	return pp, nil
}

// get PricePoint from a http message
func ParseHttpPricePoint(raw_message []byte, msgBuf map[string]interface{}) (*models.MarketPoint, error) {
	err := json.Unmarshal(raw_message, &msgBuf)
	if err != nil {
		return nil, errors.New("Could not unmarshal raw message into json from http " + string(raw_message))
	}
	bidPrice, err := strconv.ParseFloat(msgBuf["bids"].([]interface{})[0].([]interface{})[0].(string), 64)
	if err != nil {
		return nil, errors.New("Could not convert best bid price to float: " + msgBuf["bids"].([]interface{})[0].([]interface{})[0].(string))
	}
	askPrice, err := strconv.ParseFloat(msgBuf["asks"].([]interface{})[0].([]interface{})[0].(string), 64)
	if err != nil {
		return nil, errors.New("Could not convert best ask price to float: " + msgBuf["asks"].([]interface{})[0].([]interface{})[0].(string))
	}
	bidAmount, err := strconv.ParseFloat(msgBuf["bids"].([]interface{})[0].([]interface{})[1].(string), 64)
	if err != nil {
		return nil, errors.New("Could not convert best bid amount to float: " + msgBuf["bids"].([]interface{})[0].([]interface{})[1].(string))
	}
	askAmount, err := strconv.ParseFloat(msgBuf["asks"].([]interface{})[0].([]interface{})[1].(string), 64)
	if err != nil {
		return nil, errors.New("Could not convert best ask amount to float: " + msgBuf["asks"].([]interface{})[0].([]interface{})[1].(string))
	}
	return &models.MarketPoint{
		Bid: &models.PricePoint{bidPrice, bidAmount},
		Ask: &models.PricePoint{askPrice, askAmount},
	}, nil
}

// get PricePoint from a socket message
func ParseSocketPricePoint(raw_message string, msgBuf map[string][][]string) (*models.MarketPoint, error) {
	err := json.Unmarshal([]byte(raw_message), &msgBuf)
	if err != nil {
		return nil, errors.New("Could not unmarshal raw message into json: " + raw_message)
	}
	bidPrice, err := strconv.ParseFloat(msgBuf["bids"][0][0], 64)
	if err != nil {
		return nil, errors.New("Could not convert best bid price to float: " + msgBuf["bids"][0][0])
	}
	askPrice, err := strconv.ParseFloat(msgBuf["asks"][0][0], 64)
	if err != nil {
		return nil, errors.New("Could not convert best ask price to float: " + msgBuf["asks"][0][0])
	}
	bidAmount, err := strconv.ParseFloat(msgBuf["bids"][0][1], 64)
	if err != nil {
		return nil, errors.New("Could not convert best bid amount to float: " + msgBuf["bids"][0][1])
	}
	askAmount, err := strconv.ParseFloat(msgBuf["asks"][0][1], 64)
	if err != nil {
		return nil, errors.New("Could not convert best ask amount to float: " + msgBuf["asks"][0][1])
	}
	return &models.MarketPoint{
		Bid: &models.PricePoint{bidPrice, bidAmount},
		Ask: &models.PricePoint{askPrice, askAmount},
	}, nil
}
