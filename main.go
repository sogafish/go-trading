package main

import (
	"fmt"
	"go-trading/config"
)

func main() {
	fmt.Print(config.Config.ApiKey)
	fmt.Print(config.Config.ApiSecret)
}
