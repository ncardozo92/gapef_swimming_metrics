package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/ncardozo92/gapef_swimming_metrics/logging"
	"github.com/ncardozo92/gapef_swimming_metrics/user"
)

const DEV_FLAG string = "--dev"

var UserHandler user.Handler

func main() {

	if isDevEnvironment() {
		if err := godotenv.Load(".env"); err != nil {
			logging.LogFatal("cannot read .env file, %v", err.Error())
		}
	}

	// setting up the handlers
	UserHandler = user.NewUserHandler(user.NewUserRepository())

	e := echo.New()

	// Login and user CRUD
	e.POST("/login", UserHandler.Login)
	e.GET("/users", UserHandler.GetAllUsers)
	e.POST("/users", UserHandler.Create)

	if launchErr := e.Start(":8080"); launchErr != nil {
		logging.LogFatal("Cannot start application")
	}
}

func isDevEnvironment() bool {
	var result bool

	if len(os.Args) > 1 {
		for _, arg := range os.Args {
			if arg == DEV_FLAG {
				logging.LogInfo("environment set for development")
				result = true
				break
			}
		}
	}

	return result
}
