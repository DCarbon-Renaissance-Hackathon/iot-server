package routers

import (
	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/Dcarbon/iott-cloud/internal/api/ctrls"
	"github.com/Dcarbon/iott-cloud/internal/api/mids"
	"github.com/Dcarbon/iott-cloud/internal/repo"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port          int
	DBUrl         string
	JwtKey        string
	TokenDuration int64
	ChainID       int64
	CarbonVersion string
	CarbonAddress string
}

type Router struct {
	*gin.Engine
	config       Config
	auth         *mids.A2M
	iotCtrl      *ctrls.IotCtrl
	projectCtrl  *ctrls.ProjectCtrl
	userCtrl     *ctrls.UserCtrl
	proposalCtrl *ctrls.ProposalCtrl
	sensorCtrl   *ctrls.SensorCtrl
}

func NewRouter(config Config) (*Router, error) {
	err := repo.InitRepo(config.DBUrl)
	if nil != err {
		return nil, err
	}

	iotCtrl, err := ctrls.NewIotCtrl(
		&esign.TypedDataDomain{
			ChainId:           config.ChainID,
			Version:           config.CarbonVersion,
			VerifyingContract: config.CarbonAddress,
			Name:              "CARBON",
		},
	)
	if nil != err {
		return nil, err
	}

	projectCtrl, err := ctrls.NewProjectCtrl(config.DBUrl)
	if nil != err {
		return nil, err
	}

	proposalCtrl, err := ctrls.NewProposalCtrl(config.DBUrl)
	if nil != err {
		return nil, err
	}

	sensorCtrl, err := ctrls.NewSensorCtrl(iotCtrl.GetIOTRepo())
	if nil != err {
		return nil, err
	}

	userCtrl, err := ctrls.NewUserCtrl(config.DBUrl, config.JwtKey, config.TokenDuration)
	if nil != err {
		return nil, err
	}

	var r = &Router{
		Engine:       gin.Default(),
		auth:         &mids.A2M{},
		config:       config,
		iotCtrl:      iotCtrl,
		projectCtrl:  projectCtrl,
		proposalCtrl: proposalCtrl,
		userCtrl:     userCtrl,
		sensorCtrl:   sensorCtrl,
	}

	r.Use(mids.GetCORS())

	var v1 = r.Group("/api/v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	var iotRoute = v1.Group("/iots")
	{
		iotRoute.POST(
			"/",
			mids.NewA2(config.JwtKey, "iot-create").HandlerFunc,
			iotCtrl.Create,
		)
		iotRoute.PUT(
			"/:id/change-status",
			mids.NewA2(config.JwtKey, "iot-change-status").HandlerFunc,
			iotCtrl.ChangeStatus,
		)
		iotRoute.GET("/by-bb", iotCtrl.GetByBB)

		iotRoute.POST("/:iotAddr/metrics", iotCtrl.CreateMetric)
		iotRoute.GET("/:iotAddr/metrics", iotCtrl.GetMetrics)
		iotRoute.GET("/:iotAddr/metrics/:metricId", iotCtrl.GetRawMetric)

		iotRoute.POST("/:iotAddr/mint-sign/", iotCtrl.CreateMint)
		iotRoute.GET("/:iotAddr/mint-sign/", iotCtrl.GetMintSigns)

		iotRoute.GET("/seperator", iotCtrl.GetDomainSeperator)
	}

	var sensorRoute = v1.Group("/sensors")
	{
		sensorRoute.POST("/",
			mids.NewA2(config.JwtKey, "sensor-create").HandlerFunc,
			sensorCtrl.Create,
		)

		sensorRoute.PUT("/change-status",
			mids.NewA2(config.JwtKey, "sensor-change-status").HandlerFunc,
			sensorCtrl.ChangeStatus,
		)

		sensorRoute.GET("/:id", sensorCtrl.GetSensor)
		sensorRoute.GET("/", sensorCtrl.GetSensors)

		sensorRoute.POST("/sm/create", sensorCtrl.Create)
		sensorRoute.POST("/sm/create-by-iot", sensorCtrl.CreateSMByIOT)
	}

	var projectRoute = v1.Group("/projects")
	{
		projectRoute.POST(
			"/",
			mids.NewA2(config.JwtKey, "project-create").HandlerFunc,
			projectCtrl.Create,
		)
		projectRoute.GET("/", projectCtrl.GetList)
		projectRoute.GET("/by-bb", projectCtrl.GetByBB)
		projectRoute.GET("/:project_id", projectCtrl.GetByID)
		// projectRoute.PUT("/:project_id/change-status", projectCtrl.ChangeStatus)
	}

	var proposalRoute = v1.Group("/proposals")
	{
		proposalRoute.POST(
			"/",
			mids.NewA2(config.JwtKey, "").HandlerFunc,
			proposalCtrl.Create,
		)
		proposalRoute.GET("/", proposalCtrl.GetList)
		projectRoute.PUT(
			"/change-status",
			mids.NewA2(config.JwtKey, "proposals-change-status").HandlerFunc,
			proposalCtrl.ChangeStatus,
		)
	}

	var userRoute = v1.Group("/users")
	{
		userRoute.POST("/login", userCtrl.Login)
		userRoute.PUT(
			"/:id",
			mids.NewA2(config.JwtKey, "").HandlerFunc,
			userCtrl.Update,
		)
	}

	return r, nil
}
