package main

import (
	"fmt"
	"go-trading/config"
	"go-trading/utils"
)

func main() {
	utils.LogginSettings(config.Config.LogFile)
	fmt.Print(config.Config.ApiKey)
	fmt.Print(config.Config.ApiSecret)
}
