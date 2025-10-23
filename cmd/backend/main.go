package main

import (
	"fmt"
	"slices"

	"Backend/internal/app/config"
	"Backend/internal/app/dsn"
	"Backend/internal/app/handler"
	"Backend/internal/app/repository"
	"Backend/internal/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


func main() {
	router := gin.Default()
	
	// addCORS(router)

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

func addCORS(router *gin.Engine) {
	configCors := cors.DefaultConfig()

	allowed := []string{"https://andrey57x.github.io", "tauri://localhost", "http://localhost:3000"}

    configCors.AllowOriginFunc = func(origin string) bool {
        return slices.Contains(allowed, origin)
    }

    configCors.AllowCredentials = true
    
	
	router.Use(cors.New(configCors))
}
