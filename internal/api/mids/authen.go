package mids

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

var ctxKey = new(int)

var rolesTable = map[string]map[string]bool{
	"admin": {
		"": true,
	},
}

type customClaim struct {
	jwt.StandardClaims
	*ClaimModel
}

type ClaimModel struct {
	ID         int64  `json:"id,omitempty"`
	Role       string `json:"role,omitempty"`
	Name       string `json:"name,omitempty"`
	EthAddress string `json:"eth,omitempty"`
}

// Check authen and permission
type A2M struct {
	jwtKey string
	perm   string
}

func NewA2(jwtKey string, perm string) *A2M {
	var a2 = &A2M{
		jwtKey: jwtKey,
		perm:   perm,
	}
	return a2
}

func (a2 *A2M) HandlerFunc(r *gin.Context) {
	var authToken = r.GetHeader("Authorization")
	var idx = strings.Index(authToken, "Bearer ")
	if idx != 0 && len(authToken) < 10 {
		r.AbortWithError(http.StatusUnauthorized, dmodels.ErrorUnauthorized)
		return
	}

	var tokenStr = authToken[7:]
	var user, err = DecodeJWT(a2.jwtKey, tokenStr)
	if nil != err {
		fmt.Println("Decode jwt error: ", err)
		r.AbortWithError(http.StatusUnauthorized, dmodels.ErrorUnauthorized)
		return
	}

	err = hasPerm(user.Role, a2.perm)
	if nil != err {
		r.AbortWithError(http.StatusForbidden, dmodels.ErrorPermissionDenied)
		return
	}

	var ctx = context.WithValue(r.Request.Context(), ctxKey, user)
	r.Request = r.Request.WithContext(ctx)
}

// DecodeJWT :
func DecodeJWT(key string, token string) (*ClaimModel, error) {
	var claim = &customClaim{}
	jtoken, err := jwt.ParseWithClaims(
		token,
		claim,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		},
	)
	if nil != err {
		return nil, dmodels.NewError(dmodels.ECodeUnauthorized, err.Error())
	}

	if !jtoken.Valid {
		return nil, dmodels.NewError(dmodels.ECodeUnauthorized, "token is invalid")
	}

	return claim.ClaimModel, nil
}

// EncodeJWT :
func EncodeJWT(key string, user *models.User, duration int64) (string, error) {
	var claim = &customClaim{
		ClaimModel: &ClaimModel{
			ID:         user.ID,
			Role:       user.Role,
			Name:       user.Name,
			EthAddress: string(user.Address),
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + duration,
		},
	}
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(key))
}

func GetAuth(ctx context.Context) (*ClaimModel, error) {
	var user = ctx.Value(ctxKey).(*ClaimModel)
	if nil == user {
		return nil, dmodels.ErrorUnauthorized
	}
	return user, nil
}

func hasPerm(role string, perm string) error {
	if perm == "" {
		return nil
	}

	if role == "super-admin" {
		return nil
	}

	var tbl = rolesTable[role]
	if nil != tbl && tbl[perm] {
		return nil
	}
	return dmodels.ErrorPermissionDenied
}
