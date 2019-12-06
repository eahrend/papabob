package main

import (
	api "github.com/eahrend/papabob/api"
	"github.com/eahrend/papabob/common/nflapi"
	mw "github.com/eahrend/papabob/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)


func main(){
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app := gin.Default()
	app.Use(logger.SetLogger())
	app.Use(cors.Default())
	nflClient, err := nflapi.NewNFLClient()
	if err != nil {
		panic(err)
	}
	app.Use(mw.NFLClientMW(nflClient))
	api.ApplyRoutes(app)
	panic(app.Run(":" + port))
}
