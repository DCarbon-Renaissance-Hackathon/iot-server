package domain

import "github.com/Dcarbon/iott-cloud/models"

type RFilterProject struct {
	Skip  int
	Limit int
	Owner string
}

type IProject interface {
	Create(project *models.Project) error

	UpdateDesc(req *models.ProjectDescription) (*models.ProjectDescription, error)
	UpdateSpec(req *models.ProjectSpec) (*models.ProjectSpec, error)

	GetById(id int64, lang string, withSpec bool) (*models.Project, error)
	GetList(filter *RFilterProject) ([]*models.Project, error)
	GetByBB(min, max *models.Point4326, owner string) ([]*models.Project, error)

	ChangeStatus(id string, status models.ProjectStatus) error
}
