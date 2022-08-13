package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Book struct {
	ID       uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Title    string `gorm:"size:255;not null;unique" json:"title"`
	Content  string `gorm:"size:255;not null;" json:"content"`
	Author   Author `json:"author"`
	AuthorID uint32 `sql:"type:int REFERENCES authors(id) ON UPDATE CASCADE ON DELETE CASCADE" json:"author_id"`
}

func (b *Book) Prepare() {
	b.ID = 0
	b.Title = html.EscapeString(strings.TrimSpace(b.Title))
	b.Content = html.EscapeString(strings.TrimSpace(b.Content))
	b.Author = Author{}
}

func (b *Book) Validate() error {

	if b.Title == "" {
		return errors.New("Required Title")
	}
	if b.Content == "" {
		return errors.New("Required Content")
	}
	if b.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (b *Book) SaveBook(db *gorm.DB) (*Book, error) {
	var err error
	err = db.Debug().Model(&Book{}).Create(&b).Error
	if err != nil {
		return &Book{}, err
	}
	if b.ID != 0 {
		err = db.Debug().Model(&Book{}).Where("id = ?", b.AuthorID).Take(&b.Author).Error
		if err != nil {
			return &Book{}, err
		}
	}
	return b, nil
}

func (b *Book) FindAllBooks(db *gorm.DB) (*[]Book, error) {
	var err error
	books := []Book{}
	err = db.Debug().Model(&Book{}).Limit(100).Find(&books).Error
	if err != nil {
		return &[]Book{}, err
	}
	if len(books) > 0 {
		for i, _ := range books {
			err := db.Debug().Model(&Author{}).Where("id = ?", books[i].AuthorID).Take(&books[i].Author).Error
			if err != nil {
				return &[]Book{}, err
			}
		}
	}
	return &books, nil
}

func (b *Book) FindBookByID(db *gorm.DB, pid uint64) (*Book, error) {
	var err error
	err = db.Debug().Model(&Book{}).Where("id = ?", pid).Take(&b).Error
	if err != nil {
		return &Book{}, err
	}
	if b.ID != 0 {
		err = db.Debug().Model(&Author{}).Where("id = ?", b.AuthorID).Take(&b.Author).Error
		if err != nil {
			return &Book{}, err
		}
	}
	return b, nil
}

func (b *Book) Book(db *gorm.DB) (*Book, error) {

	var err error
	err = db.Debug().Model(&Book{}).Where("id = ?", b.ID).Updates(Book{Title: b.Title, Content: b.Content}).Error
	if err != nil {
		return &Book{}, err
	}
	if b.ID != 0 {
		err = db.Debug().Model(&Book{}).Where("id = ?", b.AuthorID).Take(&b.Author).Error
		if err != nil {
			return &Book{}, err
		}
	}
	return b, nil
}

func (b *Book) UpdateABook(db *gorm.DB) (*Book, error) {

	var err error
	db = db.Debug().Model(&Book{}).Where("id = ?", b.ID).Take(&Book{}).UpdateColumns(
		map[string]interface{}{
			"title":     b.Title,
			"content":   b.Content,
			"author":    b.Author,
			"author_id": b.AuthorID,
		},
	)
	err = db.Debug().Model(&Book{}).Where("id = ?", b.ID).Take(&b).Error
	if err != nil {
		return &Book{}, err
	}
	if b.ID != 0 {
		err = db.Debug().Model(&Book{}).Where("id = ?", b.AuthorID).Take(&b.Author).Error
		if err != nil {
			return &Book{}, err
		}
	}
	err = db.Debug().Model(&Book{}).Where("id = ?", b.ID).Updates(Book{Title: b.Title, Content: b.Content, Author: b.Author, AuthorID: b.AuthorID}).Error
	if err != nil {
		return &Book{}, err
	}
	if b.ID != 0 {
		err = db.Debug().Model(&Book{}).Where("id = ?", b.AuthorID).Take(&b.Author).Error
		if err != nil {
			return &Book{}, err
		}
	}
	return b, nil
}

func (b *Book) DeleteABook(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Book{}).Where("id = ? and author_id = ?", pid, uid).Take(&Book{}).Delete(&Book{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Book not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
