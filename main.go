package main

import (
	"github.com/serg2013/reading/api"
	_ "github.com/serg2013/reading/docs"
)

// @title reading API
// @version 1.0
// @description This is a sample server API

// @contact.name API Support
// @contact.url github.com/serg2013
// @contact.email raven1901@mail.ru
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @query.collection.format multi

func main() {
	api.Run()
}
