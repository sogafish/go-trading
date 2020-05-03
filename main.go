package main

import (
	"fmt"
	"go-trading/bitflyer"
	"go-trading/config"
	"go-trading/utils"
)

func main() {
	utils.LogginSettings(config.Config.LogFile)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)

	ticker, _ := apiClient.GetTicker("BTC_JPY")
	fmt.Println(ticker)
}
