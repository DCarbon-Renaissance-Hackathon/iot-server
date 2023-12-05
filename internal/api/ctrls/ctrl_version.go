package ctrls

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/gin-gonic/gin"
)

type VersionCtrl struct {
	path     string                    // Path for image
	versions map[models.IOTType]string //
}

func NewVersionCtrl() (*VersionCtrl, error) {
	var vCtrl = &VersionCtrl{
		path:     utils.StringEnv("IOT_IMAGE_PATH", "./data"),
		versions: make(map[models.IOTType]string),
	}
	vCtrl.loadVersion(models.IOTTypeBurnMethane)
	vCtrl.loadVersion(models.IOTTypeBurnBiomass)

	return vCtrl, nil
}

// Create godoc
// @Summary      Get version of iot execute
// @Description  Login
// @Tags         Version
// @Accept       json
// @Produce      json
// @Param        iotType			query		models.IOTType		true	"IOT type"
// @Success      200				{object}	RsVersion
// @Failure      400				{object}	Error
// @Failure      404				{object}	Error
// @Failure      500				{object}	Error
// @Router       /version/latest	[get]
func (ctrl *VersionCtrl) GetLatest(r *gin.Context) {
	var iotType, err = strconv.ParseInt(r.Query("iotType"), 10, 64)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("Invalid iot type"))
		return
	}
	log.Println("Version: ", ctrl.versions[models.IOTType(iotType)])
	r.JSON(200, &RsVersion{
		Version: ctrl.versions[models.IOTType(iotType)],
	})
}

// Create godoc
// @Summary      Login
// @Description  download image
// @Tags         Version
// @Accept       json
// @Produce      json
// @Param        iotType	query		models.IOTType		true	"IOT type"
// @Success      200				{object}	RsLogin
// @Failure      400				{object}	Error
// @Failure      404				{object}	Error
// @Failure      500				{object}	Error
// @Router       /version/download	[get]
func (ctrl *VersionCtrl) Download(r *gin.Context) {
	var iotType, err = strconv.ParseInt(r.Query("iotType"), 10, 64)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("Invalid iot type"))
		return
	}

	log.Println("Download file")
	var fileName = fmt.Sprintf("%s/%d/%s", ctrl.path, iotType, ctrl.versions[models.IOTType(iotType)])
	log.Println("File name: ", fileName)

	r.File(fileName)
}

func (ctrl *VersionCtrl) loadVersion(iott models.IOTType) {
	var envKey = fmt.Sprintf("IOT_VERSION_%d", iott)
	ctrl.versions[iott] = os.Getenv(envKey)
}

type RsVersion struct {
	Version string `json:"version"`
}
