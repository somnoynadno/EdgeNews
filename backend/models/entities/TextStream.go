package entities

import (
	"github.com/jinzhu/gorm"
	"time"
)

type TextStream struct {
	gorm.Model
	Name             string `gorm:"index:name;"`
	URL              string
	IsActive         bool `gorm:"not null;default:false;"`
	LastStreamUpdate *time.Time
	SourceID         uint
	Source           *Source `json:",omitempty"`
}
