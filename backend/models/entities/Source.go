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

var DefaultSources = []Source{
	{Name: "Meduza API", ScrapperTypeID: 1, Color: "purple", URL: "https://meduza.io"},       // 1
	{Name: "News API", ScrapperTypeID: 2, Color: "green", URL: "https://newsapi.org"},        // 2
	{Name: "Newscatcher API", ScrapperTypeID: 2, Color: "teal", URL: "https://rapidapi.com"}, // 3
	{Name: "Эхо Москвы", ScrapperTypeID: 3, Color: "gray", URL: "https://echo.msk.ru"},       // 4
	{Name: "OnoMedia", ScrapperTypeID: 1, Color: "yellow", URL: "https://onomedia.today"},    // 5
	{Name: "Медиазона", ScrapperTypeID: 3, Color: "red", URL: "https://zona.media"},          // 6
}
