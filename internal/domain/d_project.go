package domain

import (
	"github.com/Dcarbon/iott-cloud/internal/models"
)

type RProjectCreate struct {
	Owner    models.EthAddress            `json:"owner" `   // ETH address
	Location *models.Point4326            `json:"location"` //
	Specs    *models.ProjectSpec          `json:"specs"`    //
	Descs    []*models.ProjectDescription `json:"descs"`    //
} // @name RProjectCreate

type RProjectFilter struct {
	Skip  int    ``
	Limit int    ``
	Owner string ``
} // @name RProjectFilter

type RProjectAddImage struct {
	ProjectID int64  `json:""`
	ImgPath   string `json:""`
} //@name RProjectCreate

type IProject interface {
	Create(req *RProjectCreate) (*models.Project, error)

	UpdateDesc(req *models.ProjectDescription) (*models.ProjectDescription, error)
	UpdateSpec(req *models.ProjectSpec) (*models.ProjectSpec, error)

	GetById(id int64, lang string) (*models.Project, error)
	GetList(filter *RProjectFilter) ([]*models.Project, error)
	GetOwner(projectId int64) (string, error)

	AddImage(*RProjectAddImage) (*models.ProjectImage, error)
	ChangeStatus(id string, status models.ProjectStatus) error

	// GetByBB(min, max *models.Point4326, owner string) ([]*models.Project, error)
}
