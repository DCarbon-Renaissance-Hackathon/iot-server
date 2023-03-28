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

type ProjectType int

const (
	ProjectTypeNone        ProjectType = 0
	ProjectTypeWindPower   ProjectType = 10
	ProjectTypeSolarPower  ProjectType = 11
	ProjectTypeBurnMethane ProjectType = 20
	ProjectTypeFertilizer  ProjectType = 30
	ProjectTypeTrash       ProjectType = 31
)

type Project struct {
	ID        int64                 `json:"id" gorm:"primaryKey"`                        //
	Owner     string                `json:"owner" gorm:"index"`                          // ETH address
	Pos       *Point4326            `json:"pos" gorm:"index;type:geometry(POINT, 4326)"` //
	Status    ProjectStatus         `json:"status" `                                     //
	Type      ProjectType           `json:"type" `                                       //
	Descs     []*ProjectDescription `json:"descs" gorm:"foreignKey:ProjectID"`
	Specs     *ProjectSpec          `json:"specs" gorm:"foreignKey:ProjectID"`
	CreatedAt time.Time             `json:"createdAt"` //
	UpdatedAt time.Time             `json:"updatedAt"` //
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
