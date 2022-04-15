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

	/*host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)*/

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
