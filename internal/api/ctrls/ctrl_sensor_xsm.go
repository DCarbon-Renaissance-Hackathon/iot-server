package ctrls

import (
	"net/http"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/repo"
	"github.com/Dcarbon/iott-cloud/internal/rss"
	"github.com/gin-gonic/gin"
)

type XSMCtrl struct {
	ixsm domain.IXSM
	// sensorPusher *edef.SensorPusher
}

func NewXSMCtrl() (*XSMCtrl, error) {
	ixsm, err := repo.NewXSMImpl(rss.GetDB())
	if nil != err {
		return nil, err
	}
	var ctrl = &XSMCtrl{
		ixsm: ixsm,
	}
	return ctrl, nil
}

// Create godoc
// @Summary      Create
// @Description  create sensor
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        sensor				body		RXSMCreate	true	"Create experiment metric"
// @Success      200				{object}	models.XSMetric
// @Failure      400				{object}	Error
// @Failure      404				{object}	Error
// @Failure      500				{object}	Error
// @Router       /sensors/xsm 		[post]
func (ctrl *XSMCtrl) Create(r *gin.Context) {
	var payload = &domain.RXSMCreate{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest(err.Error()))
		return
	}

	resp, err := ctrl.ixsm.Create(payload)
	if nil != err {
		r.JSON(500, dmodels.ErrBadRequest(err.Error()))
		return
	}

	r.JSON(http.StatusOK, resp)
}

// Create godoc
// @Summary      GetList experiment metric
// @Description  Get list experiment metric
// @Tags         Sensors
// @Accept       json
// @Produce      json
// @Param        sensor				body		RXSMGetList	true	"Payload"
// @Success      200				{array}		models.XSMetric
// @Failure      400				{object}	Error
// @Failure      404				{object}	Error
// @Failure      500				{object}	Error
// @Router       /sensors/xsm 		[get]
func (ctrl *XSMCtrl) GetList(r *gin.Context) {
	var payload = &domain.RXSMGetList{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest(err.Error()))
		return
	}

	resp, err := ctrl.ixsm.GetList(payload)
	if nil != err {
		r.JSON(500, dmodels.ErrBadRequest(err.Error()))
		return
	}

	r.JSON(http.StatusOK, resp)
}
