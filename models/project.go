package models

type ProjectStatus int

const (
	ProjectStatusReject   ProjectStatus = -1
	ProjectStatusRegister ProjectStatus = 0
	ProjectStatusActived  ProjectStatus = 1
)

type Project struct {
	ID     int64         `json:"id" gorm:"primary_key"`
	Owner  string        `json:"owner" gorm:"index"` // ETH address
	Name   string        `json:"name" `
	Desc   string        `json:"desc" `
	Pos    *Point4326    `json:"pos" gorm:"index;type:geometry(POINT, 4326)"`
	Status ProjectStatus `json:"status" `
}

func (*Project) TableName() string { return TableNameProject }

// func IsValidProjectStatus(status ProjectStatus) error{
// 	return
//  }
