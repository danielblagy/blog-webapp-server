module github.com/danielblagy/blog-webapp-server

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.7
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.4.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	gorm.io/driver/postgres v1.3.4
	gorm.io/gorm v1.23.4
)

// +heroku goVersion go1.14