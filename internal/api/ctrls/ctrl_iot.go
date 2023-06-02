package ctrls

import (
	"log"
	"strconv"

	"github.com/Dcarbon/go-shared/libs/broker"
	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/Dcarbon/iott-cloud/internal/repo"
	"github.com/Dcarbon/iott-cloud/internal/rss"
	"github.com/gin-gonic/gin"
)

type IotCtrl struct {
	separator *esign.TypedDataDomain // Domain seperator
	iot       domain.IIot
	sensor    domain.ISensor
	pusher    *broker.IOTEvent
}

func NewIotCtrl(typedDomain *esign.TypedDataDomain,
) (*IotCtrl, error) {
	var dMinter, err = esign.NewERC712(
		typedDomain,
		esign.MustNewTypedDataField(
			"Mint",
			esign.TypedDataStruct,
			esign.MustNewTypedDataField("iot", esign.TypedDataAddress),
			esign.MustNewTypedDataField("amount", esign.TypedDataUint256),
			esign.MustNewTypedDataField("nonce", esign.TypedDataUint256),
		),
	)
	if nil != err {
		return nil, err
	}
	utils.Dump("Type domain config ", typedDomain)
	irepo, err := repo.NewIOTRepo(dMinter)
	if nil != err {
		return nil, err
	}

	var ctrl = &IotCtrl{
		iot:       irepo,
		separator: typedDomain,
		pusher:    broker.NewIOTEvent(rss.GetRabbitPusher()),
	}
	return ctrl, nil
}

func (ctrl *IotCtrl) SetSensor(sensor domain.ISensor) {
	ctrl.sensor = sensor
}

// Create godoc
// @Summary      Create
// @Description  create iot
// @Tags         Iots
// @Accept       json
// @Produce      json
// @Param        iot   				body		RIotCreate			true	"IOT information"
// @Param        Authorization		header		string				true	"Authorization token (`Bearer $token`)"
// @Success      200				{object}	IOTDevice
// @Failure      400				{object}	Error
// @Failure      404				{object}	Error
// @Failure      500				{object}	Error
// @Router       /iots/				[post]
func (ctrl *IotCtrl) Create(r *gin.Context) {
	var payload = &domain.RIotCreate{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, models.ErrBadRequest(err.Error()))
		return
	}

	iot, err := ctrl.iot.Create(payload)
	if nil != err {
		r.JSON(500, err)
		return
	}

	r.JSON(200, iot)
	log.Println("Publish iot created")
	ctrl.pusher.PushIOTCreate(&broker.EventIOTCreate{
		ID:      iot.ID,
		Status:  broker.DeviceStatus(iot.Status),
		Address: string(iot.Address),
		Location: &broker.GPS{
			Lng: iot.Position.Lng,
			Lat: iot.Position.Lat,
		},
	})
}

// Create godoc
// @Summary      GetIot
// @Description  Get iot by id
// @Tags         Iots
// @Accept       json
// @Produce      json
// @Param        iotId				path  		int 				true	"IOT id"
// @Success      200				{object}	IOTDevice
// @Failure      400				{object}	Error
// @Failure      404				{object}	Error
// @Failure      500				{object}	Error
// @Router       /iots/{iotId}		[get]
func (ctrl *IotCtrl) GetIot(r *gin.Context) {
	iotId, err := strconv.Atoi(r.Param("iotId"))
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Invalid iot id (Must be integer)"))
		return
	}

	iot, err := ctrl.iot.GetIOT(int64(iotId))
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(200, iot)
	}
}

// Create godoc
// @Summary      ChangeStatus
// @Description  Change iot device status
// @Tags         Iots
// @Accept       json
// @Produce      json
// @Param        payload			body		RIotChangeStatus	true	"Payload"
// @Param        iotId				path  		int 				true	"IOT id"
// @Param        Authorization		header		string				true	"Authorization token (`Bearer $token`)"
// @Success      200				{object}	IOTDevice
// @Failure      400				{object}	Error
// @Router       /iots/{iotId}/change-status [put]
func (ctrl *IotCtrl) ChangeStatus(r *gin.Context) {
	var payload = &domain.RIotChangeStatus{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, models.ErrBadRequest(err.Error()))
		return
	}

	iotId, err := strconv.Atoi(r.Param("iotId"))
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Invalid iot id (Must be integer)"))
		return
	}

	iot, err := ctrl.iot.ChangeStatus(&domain.RIotChangeStatus{
		IotId:  int64(iotId),
		Status: payload.Status,
	})
	if nil != err {
		r.JSON(500, err)
		return
	}

	if ctrl.sensor != nil {
		sensors, err := ctrl.sensor.GetSensors(&domain.RGetSensors{})
		if nil != err {
			log.Println("Get list sensor error: ", err)
		} else {
			for _, ss := range sensors {
				// if ss.Status != models.DeviceStatusRegister {
				// 	continue
				// }

				ctrl.sensor.ChangeSensorStatus(&domain.RChangeSensorStatus{
					Status: *payload.Status,
					ID:     ss.ID,
				})
			}
		}
	}

	r.JSON(200, iot)
	ctrl.pusher.PushIOTChangeStatus(&broker.EventIOTChangeStatus{
		ID:     iot.ID,
		Status: broker.DeviceStatus(iot.Status),
	})

}

// Create godoc
// @Summary      Get by bounding box
// @Description  Get iot by bounding box
// @Tags         Iots
// @Accept       json
// @Produce      json
// @Param        minLng   	query      	number  true  "Min longitude"
// @Param        minLat   	query      	number  true  "Min latitude"
// @Param        maxLng   	query      	number  true  "Max longitude"
// @Param        maxLat   	query      	number  true  "Max latitude"
// @Success      200		{array}		IOTDevice
// @Failure      400		{object}	Error
// @Router       /iots/by-bb [get]
func (ctrl *IotCtrl) GetByBB(r *gin.Context) {
	minLng, err := strconv.ParseFloat(r.Query("minLng"), 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Min longitude must be double"))
		return
	}
	minLat, err := strconv.ParseFloat(r.Query("minLat"), 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Min latitude must be double"))
		return
	}

	maxLng, err := strconv.ParseFloat(r.Query("maxLng"), 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Min longitude must be double"))
		return
	}
	maxLat, err := strconv.ParseFloat(r.Query("maxLat"), 64)
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
	data, err := ctrl.iot.GetByBB(min, max)
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(200, data)
	}
}

// GetRawMetric		godoc
// @Summary			IOT save mint signature
// @Description		IOT save mint signature
// @Tags			Iots
// @Accept			json
// @Produce			json
// @Param			iotAddr			path		string				true	"IOT address"
// @Param			iot				body		models.MintSign		true	"Signature"
// @Success			200				{object}	models.MintSign
// @Failure			400				{object}	Error
// @Failure			404				{object}	Error
// @Failure			500				{object}	Error
// @Router			/iots/{iotAddr}/mint-sign	[post]
func (ctrl *IotCtrl) CreateMint(r *gin.Context) {
	var mint = &models.MintSign{}
	var err = r.BindJSON(mint)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Payload must be json: "+err.Error()))
		return
	}

	err = ctrl.iot.CreateMint(mint)
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(200, mint)
	}
}

// GetRawMetric		godoc
// @Summary			Get mint signature of iot
// @Description		Get mint signature of iot
// @Tags			Iots
// @Accept			json
// @Produce			json
// @Param			iotId			path		number  			true  "IOT id"
// @Param			from			query		number				true  "Duration start"
// @Param			to				query		number				true  "Duration end"
// @Success			200				{array}		models.MintSign
// @Failure			400				{object}	Error
// @Failure			404				{object}	Error
// @Failure			500				{object}	Error
// @Router			/iots/{iotId}/mint-sign/ 	[get]
func (ctrl *IotCtrl) GetMintSigns(r *gin.Context) {
	iotId, err := strconv.ParseInt(r.Param("iotId"), 10, 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Invalid iot id (Must be integer)"))
		return
	}
	from, err := strconv.ParseInt(r.Query("from"), 10, 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Missing duration start"))
		return
	}

	to, err := strconv.ParseInt(r.Query("to"), 10, 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Missing duration start"))
		return
	}

	signeds, err := ctrl.iot.GetMintSigns(&domain.RIotGetMintSignList{
		From:  from,
		To:    to,
		IotId: iotId,
	})
	if nil != err {
		r.JSON(500, err)
		return
	} else {
		r.JSON(200, signeds)
	}
}

// GetDomainSeperator		godoc
// @Summary			GetDomainSeperator
// @Description		Get domain separator
// @Tags			Iots
// @Accept			json
// @Produce			json
// @Success			200				{integer}	esign.TypedDataDomain
// @Failure			400				{object}	Error
// @Failure			404				{object}	Error
// @Failure			500				{object}	Error
// @Router			/iots/seperator [get]
func (ctrl *IotCtrl) GetDomainSeperator(r *gin.Context) {
	r.JSON(200, ctrl.separator)
}

func (ctrl *IotCtrl) GetIOTRepo() domain.IIot {
	return ctrl.iot
}

// type RIOTChangeStatus struct {
// 	// Address string              `json:"address"`
// 	Status models.DeviceStatus `json:"status"`
// }

// GetRawMetric godoc
// @Summary      Create metrics (only for iot)
// @Description  Get metrics
// @Tags         Iots
// @Accept       json
// @Produce      json
// @Param        iotAddr	path			string				true	"IOT address"
// @Param        payload	body			models.Metric		true	"Metric payload"
// @Success      200		{object}		models.Metric
// @Failure      400		{object}		models.Error
// @Router       /iots/{iotAddr}/metrics 	[post]
// func (ctrl *IotCtrl) CreateMetric(r *gin.Context) {
// 	var payload = &models.Metric{}
// 	var err = r.BindJSON(payload)
// 	if nil != err {
// 		r.JSON(400, models.ErrBadRequest("Payload must be json: "+err.Error()))
// 		return
// 	}

// 	err = payload.VerifySignature()
// 	if nil != err {
// 		r.JSON(400, err)
// 	} else {

// 		err = ctrl.iotRepo.CreateMetric(payload)
// 		if nil != err {
// 			r.JSON(500, err)
// 		} else {
// 			r.JSON(http.StatusOK, payload)
// 		}
// 	}
// }

// GetRawMetric godoc
// @Summary      Get list metric of iot
// @Description  Get metrics
// @Tags         Iots
// @Accept       json
// @Produce      json
// @Param        iotAddr	path		string 			true  "IOT address"
// @Param        from		query		integer			true  "Timestamp"
// @Param        to			query		integer			true  "Timestamp"
// @Success      200		{array}		models.Metric
// @Failure      400		{object}	Error
// @Failure      404		{object}	Error
// @Failure      500		{object}	Error
// @Router       /iots/{iotAddr}/metrics [get]
// func (ctrl *IotCtrl) GetMetrics(r *gin.Context) {
// 	from, err := strconv.ParseInt(r.Query("from"), 10, 64)
// 	if nil != err {
// 		r.JSON(400, models.ErrQueryParam("from must be int"))
// 		return
// 	}

// 	to, err := strconv.ParseInt(r.Query("to"), 10, 64)
// 	if nil != err {
// 		r.JSON(400, models.ErrQueryParam("to must be int"))
// 		return
// 	}

// 	var addr = r.Param("iotAddr")
// 	log.Printf("Add:%s from:%d to:%d\n", addr, from, to)
// 	data, err := ctrl.iotRepo.GetMetrics(addr, from, to)
// 	if nil != err {
// 		r.JSON(500, err)
// 	} else {
// 		r.JSON(200, data)
// 	}
// }

// GetRawMetric		godoc
// @Summary			Show raw metric from iot
// @Description		Get all information of metric
// @Tags			Iots
// @Accept			json
// @Produce			json
// @Param			iotAddr			path				string  true  "IOT address"
// @Param			metricId		path				string  true  "Metric id"
// @Success			200				{object}			models.Metric
// @Failure			400				{object}			models.Error
// @Failure			404				{object}			models.Error
// @Failure			500				{object}			models.Error
// @Router			/iots/{iotAddr}/metrics/{metricId} [get]
// func (ctrl *IotCtrl) GetRawMetric(r *gin.Context) {
// 	var iotAddr = r.Param("iotAddr")
// 	var metricId = r.Param("metricId")
// 	if iotAddr == "" || metricId == "" {
// 		r.JSON(500, models.ErrBadRequest("Missing metric id "))
// 		return
// 	}

// 	data, err := ctrl.iotRepo.GetRawMetric(metricId)
// 	if nil != err {
// 		r.JSON(500, err)
// 	} else {
// 		r.JSON(200, data)
// 	}
// }

// GetRawMetric		godoc
// @Summary			Get mint signature of iot
// @Description		Get mint signature of iot
// @Tags			Iots
// @Accept			json
// @Produce			json
// @Param			iotAddr			path		string  			true  "IOT address"
// @Param			fromNonce		query		number				true  "LatestNonce"
// @Success			200				{integer}	integer
// @Failure			400				{object}	Error
// @Failure			404				{object}	Error
// @Failure			500				{object}	Error
// @Router			/iots/{iotAddr}/get-tt [get]
// func (ctrl *IotCtrl) GetTT(r *gin.Context) {
// 	var iotAddress = r.Param("iotAddr")
// 	if iotAddress == "" {
// 		r.JSON(400, models.ErrBadRequest("Missing iot address"))
// 		return
// 	}

// 	var fromNonce, err = strconv.ParseInt(r.Query("fromNonce"), 10, 64)
// 	if nil != err {
// 		r.JSON(400, models.ErrBadRequest("Missing iot nonce"))
// 		return
// 	}

// 	signeds, err := ctrl.iotRepo.GetMintSigns(iotAddress, int(fromNonce))
// 	if nil != err {
// 		r.JSON(500, err)
// 		return
// 	} else {
// 		r.JSON(200, signeds)
// 	}
// }
