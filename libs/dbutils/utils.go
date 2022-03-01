package dbutils

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//MustNewDB :
func MustNewDB(dbURL string) *gorm.DB {
	var db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if nil != err {
		panic("Connect database error: " + err.Error())
	}
	go func() {
		var db2, _ = db.DB()
		for {
			time.Sleep(10 * time.Minute)
			err := db2.Ping()
			if nil != err {
				log.Println("Ping to database error: ", err)
			}
		}
	}()
	return db
}

func NewDB(dbURL string) (*gorm.DB, error) {
	var db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if nil != err {
		return nil, err
	}
	go func() {
		var db2, _ = db.DB()
		for {
			time.Sleep(10 * time.Second)
			err := db2.Ping()
			if nil != err {
				log.Println("Ping to database error: ", err)
			}
		}
	}()
	return db, nil
}
