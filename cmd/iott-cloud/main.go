package main

import (
	"fmt"
	"log"

	"github.com/Dcarbon/iott-cloud/api/routers"
	"github.com/Dcarbon/iott-cloud/libs/utils"

	"github.com/Dcarbon/iott-cloud/cmd/iott-cloud/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var config = routers.Config{
	Port:          utils.IntEnv("DCENTER_PORT", 8081),
	DBUrl:         utils.StringEnv("DB_URL", ""),
	JwtKey:        utils.StringEnv("JWT_KEY", ""),
	TokenDuration: utils.Int64Env("TOKEN_DURATION", 1*365*86400),
	ChainID:       utils.Int64Env("CHAIN_ID", 1),
	CarbonVersion: utils.StringEnv("CARBON_VERSION", "1"),
	CarbonAddress: utils.StringEnv("CARBON_ADDRESS", ""),
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
	docs.SwaggerInfo.Schemes = []string{
		utils.StringEnv("SERVER_SCHEME", "http"),
	}
	docs.SwaggerInfo.Host = utils.StringEnv("SERVER_HOST", "localhost:8081")

	var rt, err = routers.NewRouter(config)
	utils.PanicError("Create router", err)

	rt.GET("/swg/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Run server at ", config.Port)
	err = rt.Run(fmt.Sprintf(":%d", config.Port))
	if nil != err {
		log.Println("Listen and serve error: ", err)
	}
}
