package main

import "fmt"

func (g *Gateio) GetTicker(symbol string) error {
	url := fmt.Sprintf("%s/%s/%s/%s", g.APIUrlSecondary, gateioAPIVersion, gateioTicker, symbol)

	var res TickerResponse
	err := g.SendHTTPRequest(url, &res)
	if err != nil {
		return res, err
	}
	return nil
}
