package entities

import "github.com/jinzhu/gorm"

type ScrapperType struct {
	gorm.Model
	Name string
}

var DefaultScrapperTypes = []ScrapperType{
	{Name: "Open API"},
	{Name: "Commercial API"},
	{Name: "Static Web-Crawler"},
}