package routers

import (
	"github.com/Dcarbon/api/controllers"
	"github.com/Dcarbon/api/mids"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port   int
	DBUrl  string
	JwtKey string
}

type Router struct {
	*gin.Engine
	config       Config
	auth         *mids.A2M
	iotCtrl      *controllers.IOTCtrl
	projectCtrl  *controllers.ProjectCtrl
	userCtrl     *controllers.UserCtrl
	proposalCtrl *controllers.ProposalCtrl
}

func NewRouter(config Config) (*Router, error) {
	var iotCtrl, err = controllers.NewIOTCtrl(config.DBUrl)
	if nil != err {
		return nil, err
	}

	projectCtrl, err := controllers.NewProjectCtrl(config.DBUrl)
	if nil != err {
		return nil, err
	}

	proposalCtrl, err := controllers.NewProposalCtrl(config.DBUrl)
	if nil != err {
		return nil, err
	}

	userCtrl, err := controllers.NewUserCtrl(config.DBUrl, config.JwtKey)
	if nil != err {
		return nil, err
	}

	var r = &Router{
		Engine:       gin.Default(),
		config:       config,
		iotCtrl:      iotCtrl,
		projectCtrl:  projectCtrl,
		proposalCtrl: proposalCtrl,
		userCtrl:     userCtrl,
		auth:         &mids.A2M{},
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
		iotRoute.GET("/", iotCtrl.GetByBB)

		iotRoute.POST("/:iotAddr/metrics", iotCtrl.CreateMetric)
		iotRoute.GET("/:iotAddr/metrics", iotCtrl.GetMetrics)
		iotRoute.GET("/:iotAddr/metrics/:metricId", iotCtrl.GetRawMetric)

		iotRoute.POST("/:iotAddr/mint-sign/", iotCtrl.CreateMint)
		iotRoute.GET("/:iotAddr/mint-sign/", iotCtrl.GetMintSigns)
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
