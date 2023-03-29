package domain

import "github.com/Dcarbon/iott-cloud/internal/models"

type IUser interface {
	Login(addr models.EthAddress, signedHex, org string) (*models.User, error)
	Update(id int64, name string) (*models.User, error)

	GetUserById(id int64) (*models.User, error)
	GetUserByAddress(addr string) (*models.User, error)
}
