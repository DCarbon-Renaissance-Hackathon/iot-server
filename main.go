package main

import (
	"fmt"
	"log"

	"github.com/Dcarbon/api/routers"
	"github.com/Dcarbon/libs/utils"

	"github.com/Dcarbon/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var config = routers.Config{
	Port:   utils.IntEnv("DCENTER_PORT", 8081),
	DBUrl:  utils.StringEnv("DB_URL", ""),
	JwtKey: utils.StringEnv("JWT_KEY", ""),
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8081
// @BasePath  /api/v1

func main() {
	docs.SwaggerInfo.Title = "Internet of trusted thing cloud"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Description = "Internet of trusted thing cloud"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	var rt, err = routers.NewRouter(config)
	utils.PanicError("Create router", err)

	rt.GET("/swg/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Run server at ", config.Port)
	err = rt.Run(fmt.Sprintf(":%d", config.Port))
	if nil != err {
		log.Println("Listen and serve error: ", err)
	}
}