package domain

import "github.com/Dcarbon/models"

type IProposal interface {
	Create(*models.Proposal) error
	GetList(skip, limit, iotID, projectId int64) ([]*models.Proposal, error)
}
