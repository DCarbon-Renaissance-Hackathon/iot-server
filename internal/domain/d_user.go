package domain

import (
	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/iott-cloud/internal/models"
)

type RLogin struct {
	Address   dmodels.EthAddress `json:"address"`
	Signature string             `json:"signature"`
	Now       int64              `json:"now"`
} //@name RLogin

// type RsLogin struct {

// }

type IUser interface {
	Login(addr dmodels.EthAddress, signedHex, org string) (*models.User, error)
	Update(id int64, name string) (*models.User, error)

	GetUserById(id int64) (*models.User, error)
	GetUserByAddress(addr string) (*models.User, error)
}
