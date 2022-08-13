package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/serg2013/reading/api/auth"
	"github.com/serg2013/reading/api/models"
	"github.com/serg2013/reading/api/responses"
	"github.com/serg2013/reading/api/utils/formaterror"
)

// CreateAuthor func creates a new author
// @Description Creates a new author
// @Summary Creates new author
// @Tags Authors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body models.Author true "book data"
// @Success 200 {object} models.Author
// @Router /authors [post]
func (server *Server) CreateAuthor(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	author := models.Author{}
	err = json.Unmarshal(body, &author)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	author.Prepare()
	err = author.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	authorCreated, err := author.SaveAuthor(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, authorCreated.ID))
	responses.JSON(w, http.StatusCreated, authorCreated)
}

// GetAuthors func gets all existing authors.
// @Description Gets all existing authors.
// @Summary Gets all existing authors
// @Tags Authors
// @Accept json
// @Produce json
// @Success 200 {array} models.Author
// @Router /authors [get]
func (server *Server) GetAuthors(w http.ResponseWriter, r *http.Request) {

	author := models.Author{}

	authors, err := author.FindAllAuthors(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, authors)
}

// GetAuthor func gets author by given ID or 404 error.
// @Description Get author by given ID.
// @Summary get author by given ID
// @Tags Authors
// @Accept json
// @Produce json
// @Param id path string true "Author ID"
// @Success 200 {object} models.Author
// @Router /authors/{id} [get]
func (server *Server) GetAuthor(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	author := models.Author{}
	authorGotten, err := author.FindAuthorByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, authorGotten)
}

// UpdateAuthor func updates existing author
// @Description Updates existing author
// @Summary Updates existing author
// @Tags Authors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Author ID"
// @Param data body models.Author true "book data"
// @Success 200 {object} models.Author
// @Router /authors/{id} [put]
func (server *Server) UpdateAuthor(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	author := models.Author{}
	err = json.Unmarshal(body, &author)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	author.Prepare()
	err = author.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedAuthor, err := author.UpdateAuthor(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedAuthor)
}

// DeleteAuthor func deletes author by given ID or 404 error.
// @Description Deletes an author by given ID.
// @Summary deletes an author by given ID
// @Tags Authors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Author ID"
// @Success 200 {object} models.Author
// @Router /authors/{id} [delete]
func (server *Server) DeleteAuthor(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	author := models.Author{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = author.DeleteAuthor(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uint32(uid)))
	responses.JSON(w, http.StatusNoContent, "")
}
