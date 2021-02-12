package entities

import (
	"github.com/jinzhu/gorm"
)

type News struct {
	gorm.Model
	Title       string `gorm:"index:title"`
	Description *string
	Date        *string
	Body        *string
	URL         *string
	Author      *string
	Rights      *string
	Tag         *string
	SourceID    uint
	Source      *Source `json:",omitempty"`
}
