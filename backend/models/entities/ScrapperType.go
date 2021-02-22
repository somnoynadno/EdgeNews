package entities

import "github.com/jinzhu/gorm"

type ScrapperType struct {
	gorm.Model
	Name string
}

var DefaultScrapperTypes = []ScrapperType{
	{Name: "Open API"},            // 1
	{Name: "Commercial API"},      // 2
	{Name: "Static Web-Crawler"},  // 3
	{Name: "Dynamic Web-Crawler"}, // 4
}
