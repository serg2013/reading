package controllers

import (
	"net/http"

	"github.com/serg2013/reading/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Testing API")
}
