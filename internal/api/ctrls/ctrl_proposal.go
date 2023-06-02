package ctrls

import (
	"net/http"

	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/Dcarbon/iott-cloud/internal/repo"
	"github.com/gin-gonic/gin"
)

type ProposalCtrl struct {
	repo domain.IProposal
}

func NewProposalCtrl(dbUrl string) (*ProposalCtrl, error) {
	var rp, err = repo.NewProposalRepo(dbUrl)
	if nil != err {
		return nil, err
	}

	var ctrl = &ProposalCtrl{
		repo: rp,
	}

	return ctrl, nil
}

// Create godoc
// @Summary      Create proposal
// @Description  Create proposal
// @Tags         Proposal
// @Accept       json
// @Produce      json
// @Param        proposal 	body      	Proposal  true  "Min longitude"
// @Success      200		{array}		Proposal
// @Failure      400		{object}	Error
// @Failure      404  		{object}	Error
// @Failure      500  		{object}	Error
// @Router       /proposals/ 	[post]
func (ctrl *ProposalCtrl) Create(r *gin.Context) {
	r.JSON(http.StatusNotImplemented, models.ErrNotImplement())
}

// Create godoc
// @Summary      ChangeStatus
// @Description  Change proposal status
// @Tags         Proposal
// @Accept       json
// @Produce      json
// @Param        proposal 	body      		Proposal		true  "Min longitude"
// @Success      200		{array}			Proposal
// @Failure      400		{object}		Error
// @Failure      404  		{object}		Error
// @Failure      500  		{object}		Error
// @Router       /proposals/change-status 	[put]
func (ctrl *ProposalCtrl) ChangeStatus(r *gin.Context) {
	r.JSON(http.StatusNotImplemented, models.ErrNotImplement())
}

// Create godoc
// @Summary      GetByList
// @Description  Get proposals
// @Tags         Proposal
// @Accept       json
// @Produce      json
// @Param        skip		query      	number  true  "Skip"
// @Param        limit		query      	number  true  "Limit"
// @Success      200		{array}		Project
// @Failure      400		{object}	Error
// @Failure      404  		{object}	Error
// @Failure      500  		{object}	Error
// @Router       /proposals/ [get]
func (ctrl *ProposalCtrl) GetList(r *gin.Context) {
	r.JSON(http.StatusNotImplemented, models.ErrNotImplement())
}
