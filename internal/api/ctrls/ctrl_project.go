package ctrls

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/libs/sclient"
	"github.com/Dcarbon/iott-cloud/internal/api/mids"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/env"
	"github.com/Dcarbon/iott-cloud/internal/repo"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type ProjectCtrl struct {
	dstTmp     string
	serverHost string
	repo       domain.IProject
	storage    sclient.IStorage
}

func NewProjectCtrl(dbUrl, storageHost, isvToken string) (*ProjectCtrl, error) {
	var projectRepo, err = repo.NewProjectRepo()
	if nil != err {
		return nil, err
	}

	storage, err := sclient.NewStorage(storageHost, isvToken)
	if nil != err {
		return nil, err
	}

	var ctrl = &ProjectCtrl{
		dstTmp:     "./static",
		serverHost: env.ServerScheme + "://" + env.ServerHost,
		repo:       projectRepo,
		storage:    storage,
	}
	return ctrl, nil
}

// Create godoc
// @Summary      Create project
// @Description  Create project
// @Tags         Project
// @Accept       json
// @Produce      json
// @Param        project		body		RProjectCreate  	true	"Project"
// @Param        Authorization	header		string				true	"Authorization token (`Bearer $token`)"
// @Success      200			{array}		Project
// @Failure      400			{object}	Error
// @Failure      404			{object}	Error
// @Failure      500			{object}	Error
// @Router       /projects/ 	[post]
func (ctrl *ProjectCtrl) Create(r *gin.Context) {
	var payload = &domain.RProjectCreate{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("Bind error: "+err.Error()))
		return
	}

	project, err := ctrl.repo.Create(payload)
	if nil != err {
		r.JSON(500, err)
		return
	}

	r.JSON(200, project)
}

// Create godoc
// @Summary      Create project
// @Description  Create project
// @Tags         Project
// @Accept       json
// @Produce      json
// @Param        project				body		RProjectUpdateDesc		true	"Project"
// @Param        Authorization			header		string					true	"Authorization token (`Bearer $token`)"
// @Success      200					{array}		ProjectDescription
// @Failure      400					{object}	Error
// @Failure      404					{object}	Error
// @Failure      500					{object}	Error
// @Router       /projects/update-desc 	[post]
func (ctrl *ProjectCtrl) UpdateDesc(r *gin.Context) {
	var payload = &domain.RProjectUpdateDesc{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("Bind error: "+err.Error()))
		return
	}

	err = ctrl.isProjectOwner(r, payload.ProjectID)
	if nil != err {
		r.JSON(http.StatusForbidden, err)
		return
	}

	project, err := ctrl.repo.UpdateDesc(payload)
	if nil != err {
		r.JSON(500, err)
		return
	}

	r.JSON(200, project)
}

// Create godoc
// @Summary      UpdateSpecs
// @Description  Update specs of project
// @Tags         Project
// @Accept       json
// @Produce      json
// @Param        project					body		RProjectUpdateSpecs	true	"Project"
// @Param        Authorization				header		string				true	"Authorization token (`Bearer $token`)"
// @Success      200						{array}		ProjectSpec
// @Failure      400						{object}	Error
// @Failure      404						{object}	Error
// @Failure      500						{object}	Error
// @Router       /projects/update-specs		[post]
func (ctrl *ProjectCtrl) UpdateSpecs(r *gin.Context) {
	var payload = &domain.RProjectUpdateSpecs{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("Bind error: "+err.Error()))
		return
	}

	err = ctrl.isProjectOwner(r, payload.ProjectID)
	if nil != err {
		r.JSON(http.StatusForbidden, err)
		return
	}

	spec, err := ctrl.repo.UpdateSpecs(payload)
	if nil != err {
		r.JSON(500, err)
		return
	}

	r.JSON(200, spec)
}

// Create godoc
// @Summary      Add image
// @Description  Add image for project
// @Tags         Project
// @Accept       mpfd
// @Produce      json
// @Param        projectId				formData	int64			true	"Project id"
// @Param        image					formData	file			true	"Project image (*.png, *.jpg)"
// @Param        Authorization			header		string			true	"Authorization token (`Bearer $token`)"
// @Success      200					{array}		Project
// @Failure      400					{object}	Error
// @Failure      404					{object}	Error
// @Failure      500					{object}	Error
// @Router       /projects/add-image	[post]
func (ctrl *ProjectCtrl) AddImage(r *gin.Context) {
	projectId, err := strconv.ParseInt(r.PostForm("projectId"), 10, 64)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("Project id must be int"))
		return
	}

	err = ctrl.isProjectOwner(r, projectId)
	if nil != err {
		r.JSON(400, err)
		return
	}

	img, err := r.FormFile("image")
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("Missing image"))
		return
	}
	os.MkdirAll(ctrl.dstTmp, 0777)

	fileName := filepath.Base(img.Filename)
	ext := filepath.Ext(fileName)
	fileName = ctrl.dstTmp + "/" + uuid.NewV4().String() + ext
	log.Println("Save file ")

	err = r.SaveUploadedFile(img, fileName)
	if nil != err {
		r.JSON(400, dmodels.ErrInternal(errors.New("missing image")))
		return
	}
	defer os.Remove(fileName)

	path, err := ctrl.storage.UploadToProject(fileName, projectId)
	if nil != err {
		r.JSON(400, err)
		return
	}

	pimg, err := ctrl.repo.AddImage(&domain.RProjectAddImage{
		ProjectID: projectId,
		ImgPath:   path,
	})
	if nil != err {
		r.JSON(500, err)
		return
	}

	r.JSON(200, pimg)
}

// Create godoc
// @Summary      GetByID
// @Description  Get project by id
// @Tags         Project
// @Accept       json
// @Produce      json
// @Param        projectId					path  		string		true	"Project id"
// @Success      200						{array}		Project
// @Failure      400						{object}	Error
// @Failure      404  						{object}	Error
// @Failure      500  						{object}	Error
// @Router       /projects/{projectId} 		[get]
func (ctrl *ProjectCtrl) GetByID(r *gin.Context) {
	id, err := strconv.ParseInt(r.Param("projectId"), 10, 64)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("projectId must be int64"))
		return
	}

	data, err := ctrl.repo.GetById(id, "vi")
	if nil != err {
		r.JSON(500, err)
		return
	}

	for _, v := range data.Images {
		v.Image = ctrl.serverHost + v.Image
	}

	r.JSON(200, data)
}

// Create godoc
// @Summary      GetList
// @Description  Get list of project by created time
// @Tags         Project
// @Accept       json
// @Produce      json
// @Param        skip		query		integer				true		"Skip"
// @Param        limit		query		integer				true		"Limit"
// @Success      200		{array}		Project
// @Failure      400		{object}	Error
// @Failure      404		{object}	Error
// @Failure      500		{object}	Error
// @Router       /projects/ [get]
func (ctrl *ProjectCtrl) GetList(r *gin.Context) {
	skip, err := strconv.ParseInt(r.DefaultQuery("skip", "0"), 10, 64)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("Skip must be int64"))
		return
	}

	limit, err := strconv.ParseInt(r.DefaultQuery("skip", "0"), 10, 64)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("limit must be int64"))
		return
	}

	owner := r.Query("owner")
	data, err := ctrl.repo.GetList(&domain.RProjectFilter{
		Skip:  int(skip),
		Limit: int(limit),
		Owner: owner,
	})
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(200, data)
	}

}

// Create godoc
// @Summary      ChangeStatus
// @Description  Change project status
// @Tags         Project
// @Accept       json
// @Produce      json
// @Param        payload	body						Project			true	"Project"
// @Param        projectId	path  						string 			true	"Project id"
// @Success      200		{array}						Project
// @Failure      400		{object}					dmodels.Error
// @Failure      404		{object}					dmodels.Error
// @Failure      500		{object}					dmodels.Error
// @Router       /projects/{projectId}/change-status 	[put]
// func (ctrl *ProjectCtrl) ChangeStatus(r *gin.Context) {
// }

func (ctrl *ProjectCtrl) isProjectOwner(r *gin.Context, projectId int64,
) error {
	user, err := mids.GetAuth(r.Request.Context())
	if nil != err {
		return dmodels.ErrInternal(errors.New("missing check authen in project add image"))
	}

	if user.Role == "super-admin" {
		return nil
	}

	owner, err := ctrl.repo.GetOwner(projectId)
	if nil != err {
		return err
	}

	if dmodels.EthAddress(user.EthAddress) != dmodels.EthAddress(owner) {
		return dmodels.ErrorPermissionDenied
	}
	return nil
}
