package ctrls

import (
	"strconv"

	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/Dcarbon/iott-cloud/internal/repo"
	"github.com/gin-gonic/gin"
)

type OperatorCtrl struct {
	iot      domain.IIot
	sensor   domain.ISensor
	operator domain.IOperator
}

func NewOperatorCtrl(iot domain.IIot, sensor domain.ISensor,
) (*OperatorCtrl, error) {
	var op, err = repo.NewOperatorRepo()
	if nil != err {
		return nil, err
	}
	sensor.SetOperatorCache(op)

	var ctrl = &OperatorCtrl{
		iot:      iot,
		sensor:   sensor,
		operator: op,
	}
	return ctrl, nil
}

// Create godoc
// @Summary      SetStatus
// @Description  IOT set status itself
// @Tags         Operator
// @Accept       json
// @Produce      json
// @Param        Authorization	header		string					true	"Authorization token (`Bearer $token`), use sign token"
// @Param        payload		body		domain.ROpSetStatus 	true	"Current iot status"
// @Success      200			{array}		Empty
// @Failure      400			{object}	models.Error
// @Failure      404			{object}	models.Error
// @Failure      500			{object}	models.Error
// @Router       /op/status		[put]
// func (ctrl *OperatorCtrl) SetStatus(r *gin.Context) {
// 	token, err := mids.GetSignAuth(r.Request.Context())
// 	if nil != err {
// 		r.JSON(401, err)
// 		return
// 	}
// 	var payload = &domain.ROpSetStatus{}
// 	err = r.Bind(payload)
// 	if nil != err {
// 		r.JSON(400, models.ErrBadRequest("Bind error: "+err.Error()))
// 		return
// 	}
// 	iot, err := ctrl.iot.GetIOTByAddress(token.Address)
// 	if nil != err {
// 		r.JSON(500, err)
// 		return
// 	}
// 	payload.Id = iot.ID
// 	err = ctrl.operator.SetStatus(payload)
// 	if nil != err {
// 		r.JSON(500, err)
// 	} else {
// 		r.JSON(200, &Empty{})
// 	}
// }

// Create godoc
// @Summary      GetStatus
// @Description  Status set status itself
// @Tags         Operator
// @Accept       json
// @Produce      json
// @Param        iotId					path  		int 				true	"IOT id"
// @Success      200					{object}	models.OpIotStatus
// @Failure      400					{object}	models.Error
// @Failure      404					{object}	models.Error
// @Failure      500					{object}	models.Error
// @Router       /op/status/{iotId}		[get]
func (ctrl *OperatorCtrl) GetStatus(r *gin.Context) {
	iotId, err := strconv.Atoi(r.Param("iotId"))
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Invalid iot id (Must be integer)"))
		return
	}

	metric, err := ctrl.operator.GetStatus(int64(iotId))
	if nil != err {
		r.JSON(400, err)
		return
	}
	r.JSON(200, metric)
}

// Create godoc
// @Summary      ChangeMetric
// @Description  IOT set metric itself
// @Tags         Operator
// @Accept       json
// @Produce      json
// @Param        Authorization		header		string					true	"Authorization token (`Bearer $token`), use sign token"
// @Success      200				{object}	models.OpSensorMetric
// @Failure      400				{object}	models.Error
// @Failure      404				{object}	models.Error
// @Failure      500				{object}	models.Error
// @Router       /op/metrics		[put]
// func (ctrl *OperatorCtrl) ChangeMetric(r *gin.Context) {
// 	var payload = &domain.ROpSetStatus{}
// 	var err = r.Bind(payload)
// 	if nil != err {
// 		r.JSON(400, models.ErrBadRequest("Bind error: "+err.Error()))
// 		return
// 	}
// 	// metric, err := ctrl.operator.ChangeMetrics(&domain.RChangeMetric{})
// 	// err = ctrl.repo.Create(payload)
// 	// if nil != err {
// 	// 	r.JSON(500, err)
// 	// } else {
// 	// 	r.JSON(200, payload)
// 	// }
// }

// Create godoc
// @Summary      GetMetrics
// @Description  Get metrics of iot
// @Tags         Operator
// @Accept       json
// @Produce      json
// @Param        iotId					path		int 				true	"IOT id"
// @Success      200					{object}	domain.RsGetMetrics
// @Failure      400					{object}	models.Error
// @Failure      404					{object}	models.Error
// @Failure      500					{object}	models.Error
// @Router       /op/metrics/{iotId}	[get]
func (ctrl *OperatorCtrl) GetMetrics(r *gin.Context) {
	iotId, err := strconv.Atoi(r.Param("iotId"))
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Invalid iot id (Must be integer)"))
		return
	}

	metric, err := ctrl.operator.GetMetrics(int64(iotId))
	if nil != err {
		r.JSON(400, err)
		return
	}
	r.JSON(200, metric)
}

type Empty struct {
}
