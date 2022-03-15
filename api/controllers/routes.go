package controllers

import "github.com/serg2013/reading/api/middlewares"

func (s *Server) initializeRoutes() {

	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	s.Router.HandleFunc("/authors", middlewares.SetMiddlewareJSON(s.CreateAuthor)).Methods("POST")
	s.Router.HandleFunc("/authors", middlewares.SetMiddlewareJSON(s.GetAuthors)).Methods("GET")
	s.Router.HandleFunc("/authors/{id}", middlewares.SetMiddlewareJSON(s.GetAuthor)).Methods("GET")
	s.Router.HandleFunc("/authors/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateAuthor))).Methods("PUT")
	s.Router.HandleFunc("/authors/{id}", middlewares.SetMiddlewareAuthentication(s.UpdateAuthor)).Methods("DELETE")

	s.Router.HandleFunc("/books", middlewares.SetMiddlewareJSON(s.CreateBook)).Methods("POST")
	s.Router.HandleFunc("/books", middlewares.SetMiddlewareJSON(s.GetBooks)).Methods("GET")
	s.Router.HandleFunc("/books/{id}", middlewares.SetMiddlewareJSON(s.GetBook)).Methods("GET")
	s.Router.HandleFunc("/books/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateBook))).Methods("PUT")
	s.Router.HandleFunc("/books/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteBook)).Methods("DELETE")
}
