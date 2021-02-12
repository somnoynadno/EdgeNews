package entities

import "github.com/jinzhu/gorm"

type Message struct {
	gorm.Model
	Body         string `gorm:"index:body;"`
	Title        *string
	Time         *string
	TextStreamID uint
	TextStream   *TextStream `json:",omitempty"`
}
