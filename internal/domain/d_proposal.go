package domain

import "github.com/Dcarbon/iott-cloud/internal/models"

type IProposal interface {
	Create(*models.Proposal) error
	GetList(skip, limit, iotID, projectId int64) ([]*models.Proposal, error)
}
