package db

import (
	"log"
	"os"

	"github.com/danielblagy/blog-webapp-server/entity"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetUpConnection() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	isHeroku := os.Getenv("IS_HEROKU")
	if isHeroku != "yes" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	dbURI := os.Getenv("DATABASE_URL")

	database, dbConnectionError := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if dbConnectionError != nil {
		log.Fatal(dbConnectionError)
		return database, dbConnectionError
	}

	// make migrations to the db (will be done only once, if the entities have never been created before)
	database.AutoMigrate(&entity.User{})
	database.AutoMigrate(&entity.Article{})

	return database, dbConnectionError
}
