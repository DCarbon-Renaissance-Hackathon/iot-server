package domain

import "github.com/Dcarbon/models"

type IProject interface {
	Create(project *models.Project) error

	GetById(id int64) (*models.Project, error)
	GetList(skip, limit int64, owner string) ([]*models.Project, error)
	GetByBB(min, max *models.Point4326, owner string) ([]*models.Project, error)

	ChangeStatus(id string, status models.ProjectStatus) error
}
