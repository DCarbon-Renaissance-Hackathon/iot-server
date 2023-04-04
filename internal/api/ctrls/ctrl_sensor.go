package ctrls

import (
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

// Create godoc
// @Summary      Create
// @Description  create sensor
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        sensor				body		domain.RCreateSensor	true	"Sensor information"
// @Param        Authorization		header		string					true	"Authorization token (`Bearer $token`)"
// @Success      200				{object}	models.Sensor
// @Failure      400				{object}	models.Error
// @Failure      404				{object}	models.Error
// @Failure      500				{object}	models.Error
// @Router       /sensors/ 			[post]
func (ctrl *SensorCtrl) Create(r *gin.Context) {
	var payload = &domain.RCreateSensor{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, models.ErrBadRequest(err.Error()))
	} else {
		sensor, err := ctrl.sensorRepo.CreateSensor(payload)
		if nil != err {
			r.JSON(500, err)
		} else {
			r.JSON(http.StatusCreated, sensor)
		}
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
// @Success      200						{object}	models.Sensor
// @Failure      400						{object}	models.Error
// @Failure      404						{object}	models.Error
// @Failure      500						{object}	models.Error
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
			r.JSON(http.StatusCreated, sensor)
		}
	}
}

// Create godoc
// @Summary      GetSensor
// @Description  Get sensor by id
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        id   				path		int				true	"Sensor id"
// @Success      200				{object}	models.Sensor
// @Failure      400				{object}	models.Error
// @Failure      404				{object}	models.Error
// @Failure      500				{object}	models.Error
// @Router       /sensors/{id} 		[get]
func (ctrl *SensorCtrl) GetSensor(r *gin.Context) {
	var id, err = strconv.ParseInt(r.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		r.JSON(400, models.ErrBadRequest("Invalid sensor id "))
	} else {
		sensor, err := ctrl.sensorRepo.GetSensor(&domain.RGetSensor{ID: id})
		if nil != err {
			r.JSON(500, err)
		} else {
			r.JSON(http.StatusCreated, sensor)
		}
	}
}

// Create godoc
// @Summary      GetSensors
// @Description  Get list of sensors
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        skip				query		int			true	"Skip"
// @Param        limit				query		int			true	"Limit"
// @Success      200				{object}	models.Sensor
// @Failure      400				{object}	models.Error
// @Failure      404				{object}	models.Error
// @Failure      500				{object}	models.Error
// @Router       /sensors/			[get]
func (ctrl *SensorCtrl) GetSensors(r *gin.Context) {
	var skip, _ = strconv.ParseInt(r.Query("skip"), 10, 64)
	var limit, _ = strconv.ParseInt(r.Query("limit"), 10, 64)
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	sensor, err := ctrl.sensorRepo.GetSensors(&domain.RGetSensors{
		Skip:  int(skip),
		Limit: int(limit),
	})
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(http.StatusCreated, sensor)
	}
}

// Create godoc
// @Summary      Create sm
// @Description  create sensor metric (for signature-enabled sensor)
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        payload			body		domain.RCreateSM		true	"Signature of metric was signed by sensor"
// @Success      200				{object}	models.Sensor
// @Failure      400				{object}	models.Error
// @Failure      404				{object}	models.Error
// @Failure      500				{object}	models.Error
// @Router       /sensors/sm/create	[post]
func (ctrl *SensorCtrl) CreateSm(r *gin.Context) {
	var payload = &domain.RCreateSM{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, models.ErrBadRequest(err.Error()))
	} else {
		sensor, err := ctrl.sensorRepo.CreateSM(payload)
		if nil != err {
			r.JSON(500, err)
		} else {
			r.JSON(http.StatusCreated, sensor)
		}
	}
}

// Create godoc
// @Summary      Create sm by iot
// @Description  create sensor metric (for signature-disabled sensor)
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        iot   				body		domain.RCreateSMFromIOT		true	"Signature of metric was signed by iot"
// @Success      200				{object}	models.Sensor
// @Failure      400				{object}	models.Error
// @Failure      404				{object}	models.Error
// @Failure      500				{object}	models.Error
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
		r.JSON(http.StatusCreated, sensor)
	}
}

// Create godoc
// @Summary      Create
// @Description  create sensor
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        payload				body		domain.RGetSM	true	"Payload"
// @Success      200					{object}	models.Sensor
// @Failure      400					{object}	models.Error
// @Failure      404					{object}	models.Error
// @Failure      500					{object}	models.Error
// @Router       /sensors/sm/			[get]
func (ctrl *SensorCtrl) GetMetrics(r *gin.Context) {
	var payload = &domain.RGetSM{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, models.ErrBadRequest(err.Error()))
	} else {
		sensor, err := ctrl.sensorRepo.GetMetrics(payload)
		if nil != err {
			r.JSON(500, err)
		} else {
			r.JSON(http.StatusCreated, sensor)
		}
	}
}
