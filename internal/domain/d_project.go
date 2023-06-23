package domain

import (
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/iott-cloud/internal/models"
)

type IProject interface {
	Create(req *RProjectCreate) (*models.Project, error)

	UpdateDesc(req *RProjectUpdateDesc) (*models.ProjectDescription, error)
	UpdateSpecs(req *RProjectUpdateSpecs) (*models.ProjectSpecs, error)

	GetById(id int64, lang string) (*models.Project, error)
	GetList(filter *RProjectFilter) ([]*models.Project, error)
	GetOwner(projectId int64) (string, error)

	AddImage(*RProjectAddImage) (*models.ProjectImage, error)
	ChangeStatus(id string, status models.ProjectStatus) error
}

type RProjectCreate struct {
	Owner        dmodels.EthAddress    `json:"owner" binding:"required"`    // ETH address
	Location     *models.Point4326     `json:"location" binding:"required"` //
	Specs        *RProjectUpdateSpecs  `json:"specs" binding:"required"`    //
	Descs        []*RProjectUpdateDesc `json:"descs" binding:"required"`    //
	Area         float64               `json:"area"`
	LocationName string                `json:"locationName"`
} // @name RProjectCreate

type RProjectUpdateDesc struct {
	ProjectID int64  `json:"projectId"` //
	Language  string `json:"language" ` //
	Name      string `json:"name"`
	Desc      string `json:"desc"`
} // @name RProjectUpdateDesc

type RProjectUpdateSpecs struct {
	ProjectID int64              `json:"projectId"`
	Specs     map[string]float64 `json:"specs"`
} //@name RProjectUpdateSpecs

type RProjectFilter struct {
	Skip  int    `json:"skip" form:"skip"`
	Limit int    `json:"limit" form:"limit;max=50"`
	Owner string `json:"owner" form:"owner"`
} // @name RProjectFilter

type RProjectAddImage struct {
	ProjectID int64  `json:"projectID"`
	ImgPath   string `json:"imgPath"`
} //@name RProjectAddImage

func (rproject *RProjectCreate) ToProject() *models.Project {
	var project = &models.Project{
		ID:        0,
		Status:    models.ProjectStatusRegister,
		Owner:     rproject.Owner,
		Location:  rproject.Location,
		Specs:     rproject.Specs.ToProjectSpecs(),
		Descs:     make([]*models.ProjectDescription, len(rproject.Descs)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for i, desc := range rproject.Descs {
		project.Descs[i] = desc.ToProjectDesc()
	}

	return project
}

func (rdesc *RProjectUpdateDesc) ToProjectDesc() *models.ProjectDescription {
	return &models.ProjectDescription{
		ID:        0,
		ProjectID: rdesc.ProjectID,
		Language:  rdesc.Language,
		Name:      rdesc.Name,
		Desc:      rdesc.Desc,
	}
}

func (rspec *RProjectUpdateSpecs) ToProjectSpecs() *models.ProjectSpecs {
	return &models.ProjectSpecs{
		ID:        0,
		ProjectID: rspec.ProjectID,
		Specs:     rspec.Specs,
	}
}
