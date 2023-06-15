package ctrls

import (
	"strconv"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/iott-cloud/internal/domain"
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
// @Summary      GetStatus
// @Description  Status set status itself
// @Tags         Operator
// @Accept       json
// @Produce      json
// @Param        iotId					path  		int 				true	"IOT id"
// @Success      200					{object}	models.OpIotStatus
// @Failure      400					{object}	Error
// @Failure      404					{object}	Error
// @Failure      500					{object}	Error
// @Router       /op/status/{iotId}		[get]
func (ctrl *OperatorCtrl) GetStatus(r *gin.Context) {
	iotId, err := strconv.Atoi(r.Param("iotId"))
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("Invalid iot id (Must be integer)"))
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
// @Summary      GetMetrics
// @Description  Get metrics of iot
// @Tags         Operator
// @Accept       json
// @Produce      json
// @Param        iotId					path		int 				true	"IOT id"
// @Success      200					{object}	domain.RsGetMetrics
// @Failure      400					{object}	Error
// @Failure      404					{object}	Error
// @Failure      500					{object}	Error
// @Router       /op/metrics/{iotId}	[get]
func (ctrl *OperatorCtrl) GetMetrics(r *gin.Context) {
	iotId, err := strconv.Atoi(r.Param("iotId"))
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("Invalid iot id (Must be integer)"))
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
