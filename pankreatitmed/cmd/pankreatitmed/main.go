package main

import (
	"pankreatitmed/internal/app/config"
	"pankreatitmed/internal/app/dsn"
	"pankreatitmed/internal/app/handler"
	"pankreatitmed/internal/app/middleware"
	"pankreatitmed/internal/app/repository"
	"pankreatitmed/internal/app/services"
	"pankreatitmed/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title PankreatitMed API
// @version 1.0
// @description Ranscon Counter
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in   header
// @name Authorization
// @description Bearer <token>
func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	jwtCfg := middleware.JWTConfig{
		Secret: conf.JWT.Secret,
		Issuer: conf.JWT.Issuer,
		TTL:    conf.JWT.TTL,
	}
	blacklist := middleware.NewRedisBlacklist(conf.Redis.Addr, conf.Redis.Password, conf.Redis.DB)
	router.Use(middleware.Auth(jwtCfg, blacklist))

	postgresString := dsn.FromEnv()

	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	svcs := services.NewServices(services.Reps{
		CriteriaRepo:             rep,
		PankreatitOrdersRepo:     rep,
		PankreatitOrderItemsRepo: rep,
		MedUsersRepo:             rep,
	}, services.Configs{JWTConfig: jwtCfg, JWTBlackList: blacklist})

	hand := handler.NewHandler(svcs)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
