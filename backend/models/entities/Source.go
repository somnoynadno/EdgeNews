package entities

import "github.com/jinzhu/gorm"

type Source struct {
	gorm.Model
	Name           string
	URL            string
	Color          string
	ScrapperTypeID uint
	ScrapperType   ScrapperType
}
