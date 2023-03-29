package models

import (
	"time"

	"github.com/Dcarbon/go-shared/libs/dbutils"
)

type ProjectStatus int

const (
	ProjectStatusReject      ProjectStatus = -1
	ProjectStatusRegister    ProjectStatus = 1
	ProjectStatusDescUpdated ProjectStatus = 2
	ProjectStatusSpecUpdated ProjectStatus = 3

	ProjectStatusActived ProjectStatus = 20
)

type Project struct {
	ID        int64                 `json:"id" gorm:"primaryKey"`                        //
	Owner     EthAddress            `json:"owner" gorm:"index"`                          // ETH address
	Pos       *Point4326            `json:"pos" gorm:"index;type:geometry(POINT, 4326)"` //
	Status    ProjectStatus         `json:"status"`                                      //
	Descs     []*ProjectDescription `json:"descs" gorm:"foreignKey:ProjectID"`           //
	Specs     *ProjectSpec          `json:"specs" gorm:"foreignKey:ProjectID"`           //
	CreatedAt time.Time             `json:"createdAt"`                                   //
	UpdatedAt time.Time             `json:"updatedAt"`                                   //
}

func (*Project) TableName() string { return TableNameProject }

type ProjectDescription struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	ProjectID int64     `json:"projectId" gorm:"index:idx_project_desc_lang,unique,priority:1"` //
	Language  string    `json:"language"  gorm:"index:idx_project_desc_lang,unique,priority:2"` //
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (*ProjectDescription) TableName() string { return TableNameProjectDesc }

type ProjectSpec struct {
	ID        int64             `json:"id" gorm:"primaryKey"`
	ProjectID int64             `json:"projectId" gorm:"unique"`
	Specs     dbutils.MapSFloat `json:"specs" gorm:"type:json"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

func (*ProjectSpec) TableName() string { return TableNameProjectSpec }
