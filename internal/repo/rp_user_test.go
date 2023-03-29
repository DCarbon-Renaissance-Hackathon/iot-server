package repo

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var uRepo domain.IUser

var adminAddr = models.EthAddress(utils.StringEnv("ADMIN_ADDRESS", ""))
var adminPrv = utils.StringEnv("ADMIN_PRIVATE", "")

var customPrv = utils.StringEnv("ADMIN_PRIVATE", "")

func init() {
	err := InitRepo(dbUrlTest)
	utils.PanicError("", err)

	uRepo, err = NewUserRepo(dbUrlTest)
	utils.PanicError("", err)
}

func TestLogin(t *testing.T) {
	var now = time.Now().Unix()
	var org = fmt.Sprintf("dcarbon_%d_%s", now, adminAddr)
	var signed, err = esign.SignPersonal(adminPrv, []byte(org))
	utils.PanicError("Login-SignPersonal", err)

	user, err := uRepo.Login(adminAddr, hexutil.Encode(signed), org)
	utils.PanicError("Login-SignPersonal", err)
	utils.Dump("Login payload", map[string]interface{}{
		"address":   adminAddr,
		"now":       now,
		"signature": hexutil.Encode(signed),
	})
	log.Println("User: ", user)
}

func TestUserUpdate(t *testing.T) {}

func TestGenerateLoginSignature(t *testing.T) {
	var now = time.Now().Unix()
	var org = fmt.Sprintf("dcarbon_%d_%s", now, customPrv)
	var signed, err = esign.SignPersonal(customPrv, []byte(org))
	utils.PanicError("Login-SignPersonal", err)

	log.Println("Signature hex: ", hexutil.Encode(signed))
	log.Println("Org: ", string(org))
}
