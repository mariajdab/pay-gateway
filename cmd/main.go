package main

import (
	"github.com/mariajdab/pay-gateway/internal/server"
	"log"
)

const port = "8080"

func main() {
	app := server.NewApp()

	if err := app.Run(port); err != nil {
		log.Fatal("Failed to listen and serve: ", err)
	}
}
