package ctrls

import (
	"strconv"

	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/Dcarbon/iott-cloud/internal/repo"
	"github.com/gin-gonic/gin"
)

type ProjectCtrl struct {
	repo domain.IProject
}

func NewProjectCtrl(dbUrl string) (*ProjectCtrl, error) {
	var projectRepo, err = repo.NewProjectRepo()
	if nil != err {
		return nil, err
	}

	var ctrl = &ProjectCtrl{
		repo: projectRepo,
	}
	return ctrl, nil
}

// Create godoc
// @Summary      Create project
// @Description  Create project
// @Tags         Project
// @Accept       json
// @Produce      json
// @Param        project		body		models.Project  true  "Project"
// @Success      200			{array}		models.Project
// @Failure      400			{object}	models.Error
// @Failure      404			{object}	models.Error
// @Failure      500			{object}	models.Error
// @Router       /projects/ 	[post]
func (ctrl *ProjectCtrl) Create(r *gin.Context) {
	var payload = &models.Project{}
	var err = r.Bind(payload)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Bind error: "+err.Error()))
		return
	}
	err = ctrl.repo.Create(payload)
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(200, payload)
	}
}

// Create godoc
// @Summary      Create project
// @Description  Create project
// @Tags         Project
// @Accept       json
// @Produce      json
// @Param        projectId					path  		string true  "Project id"
// @Success      200						{array}		models.Project
// @Failure      400						{object}	models.Error
// @Failure      404  						{object}	models.Error
// @Failure      500  						{object}	models.Error
// @Router       /projects/{projectId} 		[get]
func (ctrl *ProjectCtrl) GetByID(r *gin.Context) {
	id, err := strconv.ParseInt(r.Param("projectId"), 10, 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("ID must be int64"))
		return
	}

	data, err := ctrl.repo.GetById(id, "vi", true)
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(200, data)
	}
}

// Create godoc
// @Summary      GetList
// @Description  Get list of project by created time
// @Tags         Project
// @Accept       json
// @Produce      json
// @Param        skip		query		integer				true		"Skip"
// @Param        limit		query		integer				true		"Limit"
// @Success      200		{array}		models.Project
// @Failure      400		{object}	models.Error
// @Failure      404		{object}	models.Error
// @Failure      500		{object}	models.Error
// @Router       /projects/ [get]
func (ctrl *ProjectCtrl) GetList(r *gin.Context) {
	skip, err := strconv.ParseInt(r.DefaultQuery("skip", "0"), 10, 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Skip must be int64"))
		return
	}

	limit, err := strconv.ParseInt(r.DefaultQuery("skip", "0"), 10, 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("limit must be int64"))
		return
	}

	owner := r.Query("owner")
	data, err := ctrl.repo.GetList(&domain.RFilterProject{
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
// @Summary      GetByBB
// @Description  Get project by bounding box
// @Tags         Project
// @Accept       json
// @Produce      json
// @Param        minLng				query			number  true  "Min longitude"
// @Param        minLat				query			number  true  "Min latitude"
// @Param        maxLng				query			number  true  "Max longitude"
// @Param        maxLat				query			number  true  "Max latitude"
// @Success      200				{array}			models.Project
// @Failure      400				{object}		models.Error
// @Failure      404				{object}		models.Error
// @Failure      500				{object}		models.Error
// @Router       /projects/by-bb	[get]
func (ctrl *ProjectCtrl) GetByBB(r *gin.Context) {
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
		r.JSON(400, models.ErrBadRequest("Max longitude must be double"))
		return
	}
	maxLat, err := strconv.ParseFloat(r.Query("maxLat"), 64)
	if nil != err {
		r.JSON(400, models.ErrBadRequest("Max latitude must be double"))
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
	data, err := ctrl.repo.GetByBB(min, max, "")
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
// @Param        payload	body						models.Project  true  "Project"
// @Param        projectId	path  						string 			true  "Project id"
// @Success      200		{array}						models.Project
// @Failure      400		{object}					models.Error
// @Failure      404		{object}					models.Error
// @Failure      500		{object}					models.Error
// @Router       /projects/{projectId}/change-status 	[put]
func (ctrl *ProjectCtrl) ChangeStatus(r *gin.Context) {

}
