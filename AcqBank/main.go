package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	server *http.Server
}

func NewHandler() *Handler {
	return &Handler{}
}

func (app *App) Run(port string) error {
	router := gin.Default()

	router.Use(gin.Recovery())

	RegisterHTTPEndpoints(router)

	app.server = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return app.server.ListenAndServe()
}

func main() {
	router := gin.Default()

	router.Use(gin.Recovery())

	app := App{}

	if err := app.Run("9090"); err != nil {
		log.Fatal("Failed to listen and serve: ", err)
	}
}
