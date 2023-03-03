package repo

import (
	"github.com/Dcarbon/iott-cloud/libs/dbutils"
	"github.com/Dcarbon/iott-cloud/models"
	"gorm.io/gorm"
)

var singDB *gorm.DB

var errInit = models.NewError(models.ECodeInternal, "must call infra.InitRepo first")

func InitRepo(dbUrl string) error {
	dbutils.CreateDB(dbUrl)
	var err error

	if nil == singDB {
		singDB, err = dbutils.NewDB(dbUrl)
		if nil != err {
			return err
		}
		// singDB.Logger = logger.New(
		// 	log.New(os.Stdout, "\r\n", log.LstdFlags),
		// 	logger.Config{
		// 		LogLevel: logger.Info,
		// 	},
		// )
	}

	return nil
}

func getSingletonDB() (*gorm.DB, error) {
	if nil == singDB {
		return nil, errInit
	}
	return singDB, nil
}
