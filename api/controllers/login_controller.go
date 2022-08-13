package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/serg2013/reading/api/auth"
	"github.com/serg2013/reading/api/models"
	"github.com/serg2013/reading/api/responses"
	"github.com/serg2013/reading/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

// Login procedure gets credentials
// @Summary Checks login data
// @Description Checks user credentials
// @Tags Authorization
// @Accept  json
// @Produce  json
// @Param user body models.Cred true "Authorization"
// @Success 200 {object} models.User
// @Router /login [post]
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
