package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/serg2013/reading/api/models"
)

var users = []models.User{
	{
		Nickname: "user1",
		Email:    "user1@gmail.com",
		Password: "password",
	},
	{
		Nickname: "user2",
		Email:    "user2@gmail.com",
		Password: "password",
	},
}

var authors = []models.Author{
	{
		Name:     "Peter",
		Lastname: "Sidorov",
		Email:    "p.sidorov@mail.ru",
	},
	{
		Name:     "Vasiliy",
		Lastname: "Rogov",
		Email:    "v.rogov@mail.ru",
	},
	{
		Name:     "Valera",
		Lastname: "Antonov",
		Email:    "v.antonov@mail.ru",
	},
}

var books = []models.Book{
	{
		Title:   "Book 1",
		Content: "Plot 1",
	},
	{
		Title:   "Book 2",
		Content: "Plot 2",
	},
	{
		Title:   "Book 3",
		Content: "Plot 3",
	},
	{
		Title:   "Book 4",
		Content: "Plot 4",
	},
	{
		Title:   "Book 5",
		Content: "Plot 5",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Book{}, &models.Author{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Author{}, &models.Book{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}
	err = db.Debug().Model(&models.Book{}).AddForeignKey("author_id", "authors(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

	for i, _ := range authors {
		err = db.Debug().Model(&models.Author{}).Create(&authors[i]).Error
		if err != nil {
			log.Fatalf("cannot seed authors table: %v", err)
		}

		books[i].AuthorID = authors[i].ID

		err = db.Debug().Model(&models.Book{}).Create(&books[i]).Error
		if err != nil {
			log.Fatalf("cannot seed books table: %v", err)
		}

	}
}
