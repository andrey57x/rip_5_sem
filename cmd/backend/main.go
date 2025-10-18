package main

import (
	"fmt"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"Backend/internal/app/config"
	"Backend/internal/app/dsn"
	"Backend/internal/app/handler"
	"Backend/internal/app/repository"
	"Backend/internal/pkg"
)

func main() {
	router := gin.Default()

	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = []string{"https://andrey57x.github.io"}
	configCors.AllowCredentials = true
	
	router.Use(cors.New(configCors))

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.NewRepository(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
