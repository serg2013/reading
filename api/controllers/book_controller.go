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

// CreateBook func creates a new book
// @Description Create new book
// @Summary Create new book
// @Tags Books
// @Accept json
// @Produce json
// @Security ApiKeyAuth
//Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body models.Book true "book data"
// @Success 200 {object} models.Book
// @Router /books [post]
func (server *Server) CreateBook(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	book := models.Book{}
	err = json.Unmarshal(body, &book)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	book.Prepare()
	err = book.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != book.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	postCreated, err := book.SaveBook(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	fmt.Println(r.Host + r.URL.Path)
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))

	responses.JSON(w, http.StatusCreated, postCreated)
}

// GetBooks func gets all exists books.
// @Description Get all exists books.
// @Summary get all exists books
// @Tags Books
// @Accept json
// @Produce json
// @Success 200 {array} models.Book
// @Router /books [get]
func (server *Server) GetBooks(w http.ResponseWriter, r *http.Request) {

	book := models.Book{}

	books, err := book.FindAllBooks(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, books)
}

// GetBook func gets book by given ID or 404 error.
// @Description Get book by given ID.
// @Summary get book by given ID
// @Tags Books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book
// @Router /books/{id} [get]
func (server *Server) GetBook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	book := models.Book{}

	bookReceived, err := book.FindBookByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, bookReceived)
}

// UpdateBook func updates existing book
// @Description Update existing book
// @Summary Update existing book
// @Tags Books
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Book ID"
// @Param data body models.Book true "book data"
// @Success 200 {object} models.Book
// @Router /books/{id} [put]
func (server *Server) UpdateBook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//Check if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	book := models.Book{}
	err = server.DB.Debug().Model(models.Book{}).Where("id = ?", pid).Take(&book).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("post not found"))
		return
	}

	// If a user attempt to update a book not belonging to him
	if uid != book.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	bookUpdate := models.Book{}
	err = json.Unmarshal(body, &bookUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != bookUpdate.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	bookUpdate.Prepare()
	err = bookUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	bookUpdate.ID = book.ID

	bookUpdated, err := bookUpdate.UpdateABook(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, bookUpdated)
}

// DeleteBook func deletes book by given ID or 404 error.
// @Description Deletes a book by given ID.
// @Summary deletes a book by given ID
// @Tags Books
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book
// @Router /books/{id} [delete]
func (server *Server) DeleteBook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid book id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the book exists
	book := models.Book{}
	err = server.DB.Debug().Model(models.Book{}).Where("id = ?", pid).Take(&book).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated author, the owner of this book?
	if uid != book.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	_, err = book.DeleteABook(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
