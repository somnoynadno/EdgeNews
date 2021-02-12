package entities

import "github.com/jinzhu/gorm"

type TextStream struct {
	gorm.Model
	Name           string `gorm:"index:name;"`
	URL            string
	IsActive       bool   `gorm:"not null;default:false;"`
	SourceID       uint
	Source         *Source `json:",omitempty"`
}
