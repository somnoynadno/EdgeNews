package db

import (
	"EdgeNews/backend/models/entities"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"os"
)

var db *gorm.DB

func init() {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	conn, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			dbHost, dbPort, username, dbName, password))
	if err != nil {
		panic(err)
	} else {
		db = conn
		log.Info("[DB] Connected on " + dbHost)
	}

	err = migrateSchema()
	if err != nil {
		log.Fatal("[DB] Error migrating schema: " + err.Error())
	} else {
		log.Info("[DB] Schema migrated successfully")
	}

	err = initDefaults()
	if err != nil {
		log.Fatal("[DB] Initializing error: " + err.Error())
	}
}

func GetDB() *gorm.DB {
	return db
}

func migrateSchema() error {
	err := db.AutoMigrate(
		entities.Message{},
		entities.News{},
		entities.ScrapperType{},
		entities.Source{},
		entities.TextStream{},
	).Error

	return err
}

func initDefaults() error {
	for i, st := range entities.DefaultScrapperTypes {
		var q entities.ScrapperType
		err := db.Where("id = ?", i + 1).First(&q).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Info("[DB] Creating '" + st.Name + "' scrapper type")
				err = db.Create(&st).Error
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	for i, s := range entities.DefaultSources {
		var q entities.Source
		err := db.Where("id = ?", i + 1).First(&q).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Info("[DB] Creating '" + s.Name + "' source")
				err = db.Create(&s).Error
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
