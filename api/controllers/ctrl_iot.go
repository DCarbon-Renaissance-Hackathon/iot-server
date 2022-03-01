package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Dcarbon/domain"
	"github.com/Dcarbon/models"
	"github.com/Dcarbon/repo"
	"github.com/gin-gonic/gin"
)

type IOTCtrl struct {
	repo domain.IIot
}

func NewIOTCtrl(dbUrl string) (*IOTCtrl, error) {
	var irepo, err = repo.NewIOTRepo(dbUrl)
	if nil != err {
		return nil, err
	}
	var ctrl = &IOTCtrl{
		repo: irepo,
	}
	return ctrl, nil
}

// Create godoc
// @Summary      Create
// @Description  create iot
// @Tags         IOTS
// @Accept       json
// @Produce      json
// @Param        id   				body		models.IOTDevice	true	"IOT information"
// @Param        Authorization		header		string				true	"Authorization"
// @Success      200				{object}	models.IOTDevice
// @Failure      400				{object}	models.Error
// @Failure      404				{object}	models.Error
// @Failure      500				{object}	models.Error
// @Router       /iots/ [post]
func (ctrl *IOTCtrl) Create(r *gin.Context) {
	var iot = &models.IOTDevice{}
	var err = r.Bind(iot)
	if nil != err {
		r.JSON(400, models.ErrBadRequest(err.Error()))
	} else {
		err = ctrl.repo.Create(iot)
		if nil != err {
			r.JSON(500, err)
		} else {
			r.JSON(http.StatusCreated, iot)
		}
	}
}

// Create godoc
// @Summary      ChangeStatus
// @Description  Change iot device status
// @Tags         IOTS
// @Accept       json
// @Produce      json
// @Param        id   		body      	RIOTChangeStatus	true	"IOT address"
// @Success      200		{object}	models.IOTDevice
// @Failure      400		{object}	models.Error
// @Failure      404  		{object}	models.Error
// @Failure      500  		{object}	models.Error
// @Router       /iots/{iot_id}/change-status [put]
func (ctrl *IOTCtrl) ChangeStatus(r *gin.Context) {
	var payload = &RIOTChangeStatus{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, models.ErrBadRequest(err.Error()))
	} else {
		iot, err := ctrl.repo.ChangeStatus(payload.Address, payload.Status)
		if nil != err {
			r.JSON(500, err)
		} else {
			r.JSON(200, iot)
		}
	}
}

// Create godoc
// @Summary      GetByBB
// @Description  Get iot by bounding box
// @Tags         IOTS
// @Accept       json
// @Produce      json
// @Param        min_lng   	query      	number  true  "Min longitude"
// @Param        min_lat   	query      	number  true  "Min latitude"
// @Param        max_lng   	query      	number  true  "Max longitude"
// @Param        max_lat   	query      	number  true  "Max latitude"
// @Success      200		{array}		models.IOTDevice
// @Failure      400		{object}	models.Error
// @Failure      404  		{object}	models.Error
// @Failure      500  		{object}	models.Error
// @Router       /iots/by-bb [get]
func (ctrl *IOTCtrl) GetByBB(r *gin.Context) {
	minLng, err := strconv.ParseFloat(r.Query("min_lng"), 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Min longitude must be double"))
		return
	}
	minLat, err := strconv.ParseFloat(r.Query("min_lat"), 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Min latitude must be double"))
		return
	}

	maxLng, err := strconv.ParseFloat(r.Query("max_lng"), 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Min longitude must be double"))
		return
	}
	maxLat, err := strconv.ParseFloat(r.Query("max_lat"), 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Min latitude must be double"))
		return
	}

	var min = &models.Point4326{
		Lng: minLng,
		Lat: minLat,
	}
	var max = &models.Point4326{
		Lng: maxLng,
		Lat: maxLat,
	}
	data, err := ctrl.repo.GetByBB(min, max)
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(200, data)
	}

}

// GetRawMetric godoc
// @Summary      Create metrics (only for iot)
// @Description  Get metrics
// @Tags         IOTS
// @Accept       json
// @Produce      json
// @Param        payload	body			models.Metric	true	"IOT address"
// @Success      200		{object}		models.Metric
// @Failure      400		{object}		models.Error
// @Failure      404		{object}		models.Error
// @Failure      500		{object}		models.Error
// @Router       /iots/:iotAddr/metrics 	[post]
func (ctrl *IOTCtrl) CreateMetric(r *gin.Context) {
	var payload = &models.Metric{}
	var err = r.BindJSON(payload)
	if nil != err {
		r.JSON(400, "Bad request")
		return
	}

	err = payload.VerifySignature()
	if nil != err {
		r.JSON(400, err)
	} else {
		err = ctrl.repo.CreateMetric(payload)
		if nil != err {
			r.JSON(http.StatusCreated, payload)
		} else {
			r.JSON(500, err)
		}
	}

}

// GetRawMetric godoc
// @Summary      Get list metric of iot
// @Description  Get metrics
// @Tags         IOTS
// @Accept       json
// @Produce      json
// @Param        iotAddr	path		string  true  "IOT address"
// @Param        from		query		integer	true  "Timestamp"
// @Param        to			query		integer	true  "Timestamp"
// @Success      200		{array}		models.Metric
// @Failure      400		{object}	models.Error
// @Failure      404		{object}	models.Error
// @Failure      500		{object}	models.Error
// @Router       /iots/{iotAddr}/metrics [get]
func (ctrl *IOTCtrl) GetMetrics(r *gin.Context) {
	from, err := strconv.ParseInt(r.Query("from"), 10, 64)
	if nil != err {
		r.JSON(400, models.ErrQueryParam("from must be int"))
		return
	}
	to, err := strconv.ParseInt(r.Query("to"), 10, 64)
	if nil != err {
		r.JSON(400, models.ErrQueryParam("to must be int"))
		return
	}

	var addr = r.Param("iotAddr")
	log.Printf("Add:%s from:%d to:%d\n", addr, from, to)
	data, err := ctrl.repo.GetMetrics(addr, from, to)
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(200, data)
	}
}

// GetRawMetric		godoc
// @Summary			Show raw metric from iot
// @Description		Get all information of metric
// @Tags			IOTS
// @Accept			json
// @Produce			json
// @Param			iotAddr			path		string  true  "IOT address"
// @Param			metricId		path		string  true  "Metric id"
// @Success			200				{object}	models.Metric
// @Failure			400				{object}	models.Error
// @Failure			404				{object}	models.Error
// @Failure			500				{object}	models.Error
// @Router			/iots/{iotAddr}/metrics/{metricId} [get]
func (ctrl *IOTCtrl) GetRawMetric(r *gin.Context) {
	var iotAddr = r.Param("iotAddr")
	var metricId = r.Param("metricId")
	if iotAddr == "" || metricId == "" {
		r.JSON(500, models.ErrBadRequest(""))
		return
	}

	data, err := ctrl.repo.GetRawMetric(iotAddr, metricId)
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(200, data)
	}
}

type RIOTChangeStatus struct {
	Address string           `json:"address"`
	Status  models.IOTStatus `json:"status"`
}
