package ctrls

import (
	"errors"
	"fmt"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/iott-cloud/internal/api/mids"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/Dcarbon/iott-cloud/internal/repo"
	"github.com/gin-gonic/gin"
)

type UserCtrl struct {
	jwtKey        string
	repo          domain.IUser
	tokenDuration int64
}

func NewUserCtrl(jwtKey string, tokenDuration int64) (*UserCtrl, error) {
	userRepo, err := repo.NewUserRepo()
	if nil != err {
		return nil, err
	}

	var ctrl = &UserCtrl{
		jwtKey:        jwtKey,
		repo:          userRepo,
		tokenDuration: tokenDuration,
	}
	return ctrl, nil
}

// Create godoc
// @Summary      Login
// @Description  Login
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        payload	body		RLogin  true  "Login request"
// @Success      200		{object}	RsLogin
// @Failure      400		{object}	Error
// @Failure      404  		{object}	Error
// @Failure      500  		{object}	Error
// @Router       /users/login [post]
func (ctrl *UserCtrl) Login(r *gin.Context) {
	var payload = &domain.RLogin{}
	var err = r.BindJSON(payload)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("Body must be json"))
		return
	}

	var org = fmt.Sprintf("dcarbon_%d_%s", payload.Now, payload.Address)
	user, err := ctrl.repo.Login(payload.Address, payload.Signature, org)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest("Invalid signature"))
		return
	}

	token, err := mids.EncodeJWT(ctrl.jwtKey, user, ctrl.tokenDuration)
	if nil != err {
		r.JSON(500, dmodels.ErrInternal(err))
		return
	}

	r.JSON(200, &RsLogin{
		Token: token,
		User:  user,
	})
}

// Create godoc
// @Summary      Update
// @Description  Update user information
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user   	query      	number  true  "User information"
// @Success      200		{object}	models.User
// @Failure      400		{object}	Error
// @Failure      404  		{object}	Error
// @Failure      500  		{object}	Error
// @Router       /users/ 	[put]
func (ctrl *UserCtrl) Update(r *gin.Context) {
	var current, err = mids.GetAuth(r.Request.Context())
	if nil != err {
		r.JSON(500, dmodels.ErrInternal(errors.New("missing check user")))
	}

	var user = &models.User{}
	err = r.BindJSON(user)
	if nil != err {
		r.JSON(400, dmodels.ErrBadRequest(""))
		return
	}

	user.ID = current.ID
	user.Address = ""
	user, err = ctrl.repo.Update(user.ID, user.Name)
	if nil != err {
		r.JSON(500, err)
	} else {
		r.JSON(200, user)
	}

}

type RsLogin struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
} // @name RsLogin
