package ctrls

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/Dcarbon/iott-cloud/internal/repo"
	"github.com/gin-gonic/gin"
)

type SensorCtrl struct {
	iotRepo    domain.IIot
	sensorRepo domain.ISensor
}

func NewSensorCtrl(iotRepo domain.IIot) (*SensorCtrl, error) {
	sensor, err := repo.NewSensorRepo()
	if nil != err {
		return nil, err
	}
	var ctrl = &SensorCtrl{
		iotRepo:    iotRepo,
		sensorRepo: sensor,
	}
	return ctrl, nil
}

func (ctrl *SensorCtrl) GetSensorRepo() domain.ISensor {
	return ctrl.sensorRepo
}

// Create godoc
// @Summary      Create
// @Description  create sensor
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        sensor				body		RCreateSensor	true	"Sensor information"
// @Param        Authorization		header		string					true	"Authorization token (`Bearer $token`)"
// @Success      200				{object}	Sensor
// @Failure      400				{object}	Error
// @Failure      404				{object}	Error
// @Failure      500				{object}	Error
// @Router       /sensors/ 			[post]
func (ctrl *SensorCtrl) Create(r *gin.Context) {
	var payload = &domain.RCreateSensor{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, models.ErrBadRequest(err.Error()))
		return
	}

	_, err = ctrl.iotRepo.GetIOT(payload.IotID)
	if nil != err {
		r.JSON(500, models.ErrBadRequest(err.Error()))
		return
	}

	sensor, err := ctrl.sensorRepo.CreateSensor(payload)
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(http.StatusOK, sensor)
	}
}

// Create godoc
// @Summary      ChangeStatus
// @Description  Change status of sensor
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        payload					body		domain.RChangeSensorStatus	true	"Request payload"
// @Param        Authorization				header		string						true	"Authorization token (`Bearer $token`)"
// @Success      200						{object}	Sensor
// @Failure      400						{object}	Error
// @Failure      404						{object}	Error
// @Failure      500						{object}	Error
// @Router       /sensors/change-status		[put]
func (ctrl *SensorCtrl) ChangeStatus(r *gin.Context) {
	var payload = &domain.RChangeSensorStatus{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, models.ErrBadRequest(err.Error()))
	} else {
		sensor, err := ctrl.sensorRepo.ChangeSensorStatus(payload)
		if nil != err {
			r.JSON(500, err)
		} else {
			r.JSON(http.StatusOK, sensor)
		}
	}
}

// Create godoc
// @Summary      GetSensor
// @Description  Get sensor by id
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        id					path		int					true	"Sensor id"
// @Success      200				{object}	Sensor
// @Failure      400				{object}	Error
// @Failure      404				{object}	Error
// @Failure      500				{object}	Error
// @Router       /sensors/{id} 		[get]
func (ctrl *SensorCtrl) GetSensor(r *gin.Context) {
	var id, err = strconv.ParseInt(r.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		r.JSON(400, models.ErrBadRequest("Invalid sensor id "))
	} else {
		sensor, err := ctrl.sensorRepo.GetSensor(&domain.SensorID{ID: id})
		if nil != err {
			r.JSON(500, err)
		} else {
			r.JSON(http.StatusOK, sensor)
		}
	}
}

// Create godoc
// @Summary      GetSensors
// @Description  Get list of sensors. Only use one of iot_id or iot_address
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        skip				query		int					false	"Skip"
// @Param        limit				query		int					false	"Limit"
// @Param        iot_id				query		int					false	"IOT id, only use iot_id or iot_address"
// @Param        iot_address		query		string				false	"IOT address, only use iot_id or iot_address"
// @Success      200				{array}		Sensor
// @Failure      400				{object}	Error
// @Failure      404				{object}	Error
// @Failure      500				{object}	Error
// @Router       /sensors/			[get]
func (ctrl *SensorCtrl) GetSensors(r *gin.Context) {
	var skip, _ = strconv.ParseInt(r.Query("skip"), 10, 64)
	var limit, _ = strconv.ParseInt(r.Query("limit"), 10, 64)
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	var iotId, _ = strconv.ParseInt(r.Query("iot_id"), 10, 64)
	var iotAddr = models.EthAddress(r.Query("iot_address"))

	if iotId <= 0 && !iotAddr.IsEmpty() {
		iot, err := ctrl.iotRepo.GetIOTByAddress(iotAddr)
		if nil != err {
			r.JSON(500, err)
			return
		}
		iotId = iot.ID
	}

	sensors, err := ctrl.sensorRepo.GetSensors(&domain.RGetSensors{
		Skip:  int(skip),
		Limit: int(limit),
		IotId: iotId,
	})
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(http.StatusOK, sensors)
	}
}

// Create godoc
// @Summary      Create sm
// @Description  create sensor metric (for signature-enabled sensor)
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        payload			body		domain.RCreateSM		true	"Signature of metric was signed by sensor"
// @Success      200				{object}	Sensor
// @Failure      400				{object}	Error
// @Failure      404				{object}	Error
// @Failure      500				{object}	Error
// @Router       /sensors/sm/create	[post]
func (ctrl *SensorCtrl) CreateSm(r *gin.Context) {
	var payload = &domain.RCreateSM{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, models.ErrBadRequest(err.Error()))
		return
	}

	if payload.SensorAddress == "" || payload.Signed == "" || payload.Data == "" {
		r.JSON(400, models.ErrBadRequest("Request missing param. Please check again"))
		return
	}

	sensor, err := ctrl.sensorRepo.CreateSM(payload)
	if nil != err {
		log.Println("Create sm error: ", err)
		r.JSON(500, err)
	} else {
		r.JSON(http.StatusOK, sensor)
	}
}

// Create godoc
// @Summary      Create sm by iot
// @Description  create sensor metric (for signature-disabled sensor)
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        iot   				body		domain.RCreateSMFromIOT		true	"Signature of metric was signed by iot"
// @Success      200				{object}	Sensor
// @Failure      400				{object}	Error
// @Failure      404				{object}	Error
// @Failure      500				{object}	Error
// @Router       /sensors/sm/create-by-iot	[post]
func (ctrl *SensorCtrl) CreateSMByIOT(r *gin.Context) {
	var payload = &domain.RCreateSMFromIOT{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, models.ErrBadRequest(err.Error()))
		return
	}

	if payload.IotAddress.IsEmpty() {
		r.JSON(400, models.ErrBadRequest("Missing iot address"))
		return
	}

	iot, err := ctrl.iotRepo.GetIOTByAddress(payload.IotAddress)
	if nil != err {
		r.JSON(500, err)
		return
	}

	payload.IotID = iot.ID
	sensor, err := ctrl.sensorRepo.CreateSMFromIot(payload)
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(http.StatusOK, sensor)
	}
}

// Create godoc
// @Summary      GetSensorMetrics
// @Description  create sensor
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        from				query		int  			true	"From unix (second)"
// @Param        to					query		int  			true	"To unix (second)"
// @Param        iotId				query		int  			true	"Iot id"
// @Param        skip				query		int  			false	"Skip"
// @Param        limit				query		int  			true	"Limit (max: 50)"
// @Param        sensorId			query		int  			false	"Sensor id"
// @Success      200				{object}	SensorMetrics
// @Failure      400				{object}	Error
// @Failure      404				{object}	Error
// @Failure      500				{object}	Error
// @Router       /sensors/sm		[get]
func (ctrl *SensorCtrl) GetMetrics(r *gin.Context) {
	var payload = &domain.RGetSM{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, models.ErrBadRequest(err.Error()))
		return
	}

	metrics, err := ctrl.sensorRepo.GetMetrics(payload)
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(http.StatusOK, &SensorMetrics{Metrics: metrics})
	}
}

type SensorMetrics struct {
	Metrics []*domain.Metric `json:"metrics"`
}
