package pkg

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"pankreatitmed/docs"
	"pankreatitmed/internal/app/config"
	"pankreatitmed/internal/app/handler"
)

type Application struct {
	Config  *config.Config
	Router  *gin.Engine
	Handler *handler.Handler
}

func NewApp(c *config.Config, r *gin.Engine, h *handler.Handler) *Application {
	return &Application{
		Config:  c,
		Router:  r,
		Handler: h,
	}
}

func (a *Application) RunApp() {
	logrus.Info("Server start up")

	docs.SwaggerInfo.Title = "PankreatitMed API"
	docs.SwaggerInfo.Description = "Ranson Criteria Counter and Order Management"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	//docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", a.Config.ServicePort)

	a.Handler.RegisterRoutes(a.Router)
	//a.Handler.RegisterStatic(a.Router)

	a.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	a.Router.GET("/openapi.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})

	serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
	fmt.Println(serverAddress)
	if err := a.Router.Run(serverAddress); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Server down")
}
