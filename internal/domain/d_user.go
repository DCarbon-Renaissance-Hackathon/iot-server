package domain

import "github.com/Dcarbon/iott-cloud/internal/models"

type RLogin struct {
	Address   models.EthAddress `json:"address"`
	Signature string            `json:"signature"`
	Now       int64             `json:"now"`
} //@name RLogin

// type RsLogin struct {

// }

type IUser interface {
	Login(addr models.EthAddress, signedHex, org string) (*models.User, error)
	Update(id int64, name string) (*models.User, error)

	GetUserById(id int64) (*models.User, error)
	GetUserByAddress(addr string) (*models.User, error)
}
