package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type ProjectStatus int

const (
	ProjectStatusReject   ProjectStatus = -1
	ProjectStatusRegister ProjectStatus = 1

	ProjectStatusActived ProjectStatus = 20
)

type Project struct {
	ID        int64                 `json:"id" gorm:"primaryKey"`                         //
	Owner     EthAddress            `json:"owner" gorm:"index"`                           // ETH address
	Status    ProjectStatus         `json:"status"`                                       //
	Location  *Point4326            `json:"location" gorm:"type:geometry(POINT, 4326)"`   //
	Specs     *ProjectSpecs         `json:"specs,omitempty" gorm:"foreignKey:ProjectID"`  //
	Descs     []*ProjectDescription `json:"descs,omitempty" gorm:"foreignKey:ProjectID"`  //
	Images    []*ProjectImage       `json:"images,omitempty" gorm:"foreignKey:ProjectID"` //
	CreatedAt time.Time             `json:"createdAt"`                                    //
	UpdatedAt time.Time             `json:"updatedAt"`                                    //
} //@name Project

func (*Project) TableName() string { return TableNameProject }

type ProjectDescription struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	ProjectID int64     `json:"projectId" gorm:"index:idx_project_desc_lang,unique,priority:1"` //
	Language  string    `json:"language"  gorm:"index:idx_project_desc_lang,unique,priority:2"` //
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
} //@name ProjectDescription

func (*ProjectDescription) TableName() string { return TableNameProjectDesc }

type ProjectSpecs struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	ProjectID int64     `json:"projectId" gorm:"unique"`
	Specs     MapSFloat `json:"specs" gorm:"type:json"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
} //@name ProjectSpec

func (*ProjectSpecs) TableName() string { return TableNameProjectSpecs }

type ProjectImage struct {
	ID        int64     `json:"id"`        //
	ProjectID int64     `json:"projectId"` //
	Image     string    `json:"image"`     // Image path
	CreatedAt time.Time `json:"createdAt"`
}

func (*ProjectImage) TableName() string { return TableNameProjectImage }

type MapSFloat map[string]float64 //@name MapSFloat

func (m *MapSFloat) Scan(value interface{}) error {
	if nil == m {
		m = new(MapSFloat)
	}
	switch vt := value.(type) {
	case string:
		return json.Unmarshal([]byte(vt), m)
	case []byte:
		return json.Unmarshal(vt, m)
	}
	return errors.New("scan value type for MapSFloat invalid")
}

func (m MapSFloat) Value() (driver.Value, error) {
	if nil == m {
		return nil, nil
	}
	return json.Marshal(m)
}
