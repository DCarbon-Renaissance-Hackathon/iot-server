package dbutils

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/jackc/pgx/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MustNewDB :
func MustNewDB(dbURL string) *gorm.DB {
	var db, err = NewDB(dbURL)
	if nil != err {
		panic("Connect database error: " + err.Error())
	}

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

func CreateDB(dbURL string) {
	parsed, err := pgx.ParseConfig(dbURL)
	if nil != err {
		panic(fmt.Sprintf("Parse db_url:%s error: %s", dbURL, err.Error()))
	}

	cmd := exec.Command(
		"createdb",
		"-h", parsed.Host,
		"-p", fmt.Sprintf("%d", parsed.Port),
		"-U", parsed.User,
		parsed.Database,
	)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", parsed.Password))
	err = cmd.Run()
	log.Println("Create db error: ", err)
}
