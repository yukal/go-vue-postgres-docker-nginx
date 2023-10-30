package main

import (
	"fmt"
	"os"
	"time"
	"yu/goweb/internal/app/Repository/postgres"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.DebugMode)
	// gin.SetMode(gin.ReleaseMode)

	var dbUser, dbPass, appHost, appPort, postgresConnStr string
	var server *gin.Engine

	// ctx := context.Background()
	if appHost = os.Getenv("HOSTNAME"); appHost == "" {
		appHost = "localhost"
	}
	if appPort = os.Getenv("API_PORT"); appPort == "" {
		appPort = "50598"
	}

	if DB_USER_FILE := os.Getenv("DB_USER_FILE"); DB_USER_FILE != "" {
		dbUser = MustReadFile(DB_USER_FILE)
	}

	if DB_PASS_FILE := os.Getenv("DB_PASS_FILE"); DB_PASS_FILE != "" {
		dbPass = MustReadFile(DB_PASS_FILE)
	}

	postgresConnStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	if gin.Mode() == gin.DebugMode {

		server = gin.Default()

		server.SetTrustedProxies(nil)
		server.ForwardedByClientIP = true

		// (!) Dont use in production mode
		server.Static("/img", "../data/img")

		corsConfig := cors.Config{
			AllowAllOrigins: true,
			// AllowOrigins: []string{
			// 	"http://localhost:8080/",
			// 	"http://127.0.0.1:8080/",
			// },
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders: []string{
				"Origin",
				"Content-Length",
				"Content-Type",
				"X-Get-Fields",
				"X-Data-Type",
				"X-Jointly",
				"X-Role",
			},
			ExposeHeaders: []string{
				"Content-Length",
			},
			AllowCredentials: true,

			MaxAge: 12 * time.Hour,
		}

		server.Use(cors.New(corsConfig))

	} else {

		server = gin.New()
		server.SetTrustedProxies([]string{"127.0.0.1", "localhost"})

	}

	postgresDb := postgres.MustOpen(postgresConnStr)
	// postgresDb := postgres.MustOpen(
	// 	settings.postgres.UrlWithCreds("user", "pass"))

	InitRoutes(server, postgresDb)

	host := appHost + ":" + appPort
	server.Run(host)
}
