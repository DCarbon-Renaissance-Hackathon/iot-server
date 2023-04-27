package routers

import (
	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/Dcarbon/iott-cloud/internal/api/ctrls"
	"github.com/Dcarbon/iott-cloud/internal/api/mids"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/Dcarbon/iott-cloud/internal/rss"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port          int
	DBUrl         string
	RedisUrl      string
	JwtKey        string
	TokenDuration int64
	ChainID       int64
	CarbonVersion string
	CarbonAddress string

	StorageHost string
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
	operatorCtrl *ctrls.OperatorCtrl
}

func NewRouter(config Config,
) (*Router, error) {
	rss.SetUrl(config.DBUrl, config.RedisUrl)

	isvToken, err := GetInternalToken(config.JwtKey)
	if nil != err {
		return nil, err
	}

	projectCtrl, err := ctrls.NewProjectCtrl(config.DBUrl, config.StorageHost, isvToken)
	if nil != err {
		return nil, err
	}

	proposalCtrl, err := ctrls.NewProposalCtrl(config.DBUrl)
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

	sensorCtrl, err := ctrls.NewSensorCtrl(iotCtrl.GetIOTRepo())
	if nil != err {
		return nil, err
	}
	iotCtrl.SetSensor(sensorCtrl.GetSensorRepo())

	userCtrl, err := ctrls.NewUserCtrl(config.JwtKey, config.TokenDuration)
	if nil != err {
		return nil, err
	}

	opCtrl, err := ctrls.NewOperatorCtrl(iotCtrl.GetIOTRepo(), sensorCtrl.GetSensorRepo())
	if nil != err {
		return nil, err
	}

	// signVerifier := mids.NewSignedAuth()

	var r = &Router{
		Engine:       gin.Default(),
		auth:         &mids.A2M{},
		config:       config,
		iotCtrl:      iotCtrl,
		projectCtrl:  projectCtrl,
		proposalCtrl: proposalCtrl,
		userCtrl:     userCtrl,
		sensorCtrl:   sensorCtrl,
		operatorCtrl: opCtrl,
	}

	r.Engine.MaxMultipartMemory = 25 << 20
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
			"/:iotId/change-status",
			mids.NewA2(config.JwtKey, "iot-change-status").HandlerFunc,
			iotCtrl.ChangeStatus,
		)

		iotRoute.GET("/:iotId", iotCtrl.GetIot)
		iotRoute.GET("/:iotId/mint-sign/", iotCtrl.GetMintSigns)
		iotRoute.GET("/seperator", iotCtrl.GetDomainSeperator)

		iotRoute.POST("/:iotAddr/mint-sign/", iotCtrl.CreateMint)

		// iotRoute.GET("/by-bb", iotCtrl.GetByBB)
		// iotRoute.POST("/:iotAddr/metrics", iotCtrl.CreateMetric)
		// iotRoute.GET("/:iotAddr/metrics", iotCtrl.GetMetrics)
		// iotRoute.GET("/:iotAddr/metrics/:metricId", iotCtrl.GetRawMetric)
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

		sensorRoute.POST("/sm/create", sensorCtrl.CreateSm)
		sensorRoute.POST("/sm/create-by-iot", sensorCtrl.CreateSMByIOT)
	}

	var opRoute = v1.Group("/op")
	{
		opRoute.GET("/status/:iotId", opCtrl.GetStatus)
		opRoute.GET("/metrics/:iotId", opCtrl.GetMetrics)
	}

	var projectRoute = v1.Group("/projects")
	{
		projectRoute.POST(
			"/",
			mids.NewA2(config.JwtKey, "project-create").HandlerFunc,
			projectCtrl.Create,
		)
		projectRoute.POST(
			"/add-image",
			mids.NewA2(config.JwtKey, "project-create").HandlerFunc,
			projectCtrl.AddImage,
		)
		projectRoute.GET("/", projectCtrl.GetList)
		projectRoute.GET("/:projectId", projectCtrl.GetByID)

		// projectRoute.GET("/by-bb", projectCtrl.GetByBB)
		// projectRoute.PUT("/:projectId/change-status", projectCtrl.ChangeStatus)
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

func GetInternalToken(jwtKey string) (string, error) {
	return mids.EncodeJWT(jwtKey, &models.User{
		ID:   1,
		Role: "super-admin",
	}, 10*365*86400)
}
