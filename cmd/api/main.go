package main

import (
	"health/configs"
	"health/server"
	"log"

	"github.com/spf13/viper"
)

func main() {
	// Загружаем конфиги, они доступны через viper.Get(...)
	configs.Init()

	app := server.InitApp()

	if err := app.Run(viper.GetString("app.port")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
