package models

import (
	"time"

	"github.com/Dcarbon/libs/dbutils"
)

type ProposalType int8
type ProposalStatus int8

type Proposal struct {
	ID         int64           `json:"id" gorm:"primary_key"`
	Type       int32           `json:"type"`
	Status     ProposalStatus  `json:""`
	Url        string          `json:"url"`
	Title      string          `json:"title"`
	Summary    string          `json:"summary"`
	ProjectId  int64           `json:"projectId" gorm:"index"`
	IOTId      int64           `json:"iotId" gorm:"index"`
	Attachment dbutils.Strings `json:"attachment" gorm:"type:json"`
	CreatedAt  time.Time       `json:"createdAt" gorm:"index"`
}

func (*Proposal) TableName() string { return TableNameProposal }
