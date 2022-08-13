package models

import (
	"errors"
	"html"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
)

type Author struct {
	ID       uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name     string `gorm:"size:255;not null;unique" json:"name"`
	Lastname string `gorm:"size:255;not null;unique" json:"lastname"`
	Email    string `gorm:"size:100;not null;unique" json:"email"`
}

func (a *Author) Prepare() {
	a.ID = 0
	a.Name = html.EscapeString(strings.TrimSpace(a.Name))
	a.Lastname = html.EscapeString(strings.TrimSpace(a.Lastname))
	a.Email = html.EscapeString(strings.TrimSpace(a.Email))
}

func (a *Author) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if a.Name == "" {
			return errors.New("Required Name")
		}
		if a.Lastname == "" {
			return errors.New("Required Lastname")
		}
		if a.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(a.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if a.Name == "" {
			return errors.New("Required Name")
		}
		if a.Lastname == "" {
			return errors.New("Required Lastname")
		}
		if a.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(a.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (a *Author) SaveAuthor(db *gorm.DB) (*Author, error) {
	var err error
	err = db.Debug().Create(&a).Error
	if err != nil {
		return &Author{}, err
	}
	return a, nil
}

func (a *Author) FindAllAuthors(db *gorm.DB) (*[]Author, error) {
	var err error
	authors := []Author{}
	err = db.Debug().Model(&Author{}).Limit(100).Find(&authors).Error
	if err != nil {
		return &[]Author{}, err
	}
	return &authors, err
}

func (a *Author) FindAuthorByID(db *gorm.DB, uid uint32) (*Author, error) {
	var err error
	err = db.Debug().Model(Author{}).Where("id = ?", uid).Take(&a).Error
	if err != nil {
		return &Author{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Author{}, errors.New("User Not Found")
	}
	return a, err
}

func (a *Author) UpdateAuthor(db *gorm.DB, uid uint32) (*Author, error) {
	db = db.Debug().Model(&Author{}).Where("id = ?", uid).Take(&Author{}).UpdateColumns(
		map[string]interface{}{
			"name":     a.Name,
			"lastname": a.Lastname,
			"email":    a.Email,
		},
	)
	if db.Error != nil {
		return &Author{}, db.Error
	}
	// This is the display the updated user
	err := db.Debug().Model(&Author{}).Where("id = ?", uid).Take(&a).Error
	if err != nil {
		return &Author{}, err
	}
	return a, nil
}

func (a *Author) DeleteAuthor(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Author{}).Where("id = ?", uid).Take(&Author{}).Delete(&Author{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Author not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
