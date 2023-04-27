package domain

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestGenerateToken(t *testing.T) {
	generateSign("", "")
}

func TestVerifyToken(t *testing.T) {

}

func generateSign(addr string, pk string) (*SignedToken, error) {
	var signedAt = time.Now().Unix() - 1
	addr = strings.ToLower(addr)

	var org = fmt.Sprintf("dcarbon_%d_%s", signedAt, addr)
	var signed, err = esign.SignPersonal(pk, []byte(org))
	if nil != err {
		return nil, err
	}
	return &SignedToken{
		SignedAt: signedAt,
		Address:  models.EthAddress(addr),
		Signed:   hexutil.Encode(signed),
	}, nil
}
