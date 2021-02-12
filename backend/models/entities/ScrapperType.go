package entities

import "github.com/jinzhu/gorm"

type ScrapperType struct {
	gorm.Model
	Name string
}
