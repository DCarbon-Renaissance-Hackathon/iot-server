package repo

import (
	"github.com/Dcarbon/domain"
	"github.com/Dcarbon/libs/dbutils"
	"github.com/Dcarbon/libs/esign"
	"github.com/Dcarbon/models"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(dbUrl string) (domain.IUser, error) {
	var db, err = dbutils.NewDB(dbUrl)
	if nil != err {
		return nil, err
	}
	err = db.AutoMigrate(
		&models.User{},
	)
	if nil != err {
		return nil, err
	}

	var up = &userRepo{
		db: db,
	}
	return up, nil
}

func (up *userRepo) Login(addr string, signedHex, org string,
) (*models.User, error) {
	var signedBytes, err = hexutil.Decode(signedHex)
	if nil != err {
		return nil, models.ErrBadRequest("Invalid sign " + err.Error())
	}

	err = esign.VerifyPersonalSign(addr, []byte(org), signedBytes)
	if nil != err {
		return nil, models.ErrBadRequest("Invalid signed" + err.Error())
	}

	var user = &models.User{
		EAddress: addr,
	}
	err = up.tblUser().
		Where("e_address = ?", addr).
		First(user).Error
	if nil != err {
		if err == gorm.ErrRecordNotFound {
			err = up.tblUser().Create(user).Error
			if nil != err {
				return nil, models.ParsePostgresError("User", err)
			}
		} else {
			return nil, models.ParsePostgresError("User", err)
		}

	}

	return user, nil
}

func (up *userRepo) Update(id int64, name string) (*models.User, error) {
	var user = &models.User{}
	var err = up.tblUser().
		Model(user).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Update("name = ?", name).
		Error
	return user, models.ParsePostgresError("User", err)
}

func (up *userRepo) GetUserById(id int64) (*models.User, error) {
	var user = &models.User{}
	var err = up.tblUser().
		Where("id = ?", id).
		First(user).Error

	return user, models.ParsePostgresError("User", err)
}

func (up *userRepo) GetUserByAddress(addr string) (*models.User, error) {
	var user = &models.User{}
	var err = up.tblUser().
		Where("e_address = ?", addr).
		First(user).Error

	return user, models.ParsePostgresError("User", err)
}

func (up *userRepo) tblUser() *gorm.DB {
	return up.db.Table(models.TableNameUser)
}
