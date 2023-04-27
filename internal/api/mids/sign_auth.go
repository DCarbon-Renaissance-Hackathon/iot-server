package mids

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type SignedAuth struct {
	verifier domain.ISignerVerifier
}

func NewSignedAuth() *SignedAuth {
	return &SignedAuth{
		verifier: domain.NewVerifier(),
	}
}

func (sa *SignedAuth) HandlerFunc(r *gin.Context) {
	var authToken = r.GetHeader("Authorization")
	var idx = strings.Index(authToken, "Bearer ")
	if idx != 0 && len(authToken) < 10 {
		r.AbortWithError(http.StatusUnauthorized, models.ErrorUnauthorized)
		return
	}

	var tokenStr = authToken[7:]
	var token = &domain.SignedToken{}
	var err = json.Unmarshal([]byte(tokenStr), token)
	if nil != err {
		r.AbortWithError(
			http.StatusUnauthorized,
			models.NewError(models.ECodeUnauthorized, "Invalid sign token. It must be sign verify"),
		)
		return
	}

	var ctx = context.WithValue(r.Request.Context(), ctxKey, token)
	r.Request = r.Request.WithContext(ctx)
}

func GetSignAuth(ctx context.Context) (*domain.SignedToken, error) {
	var auth = ctx.Value(ctxKey).(*domain.SignedToken)
	if nil == auth {
		return nil, models.ErrorUnauthorized
	}
	return auth, nil
}
