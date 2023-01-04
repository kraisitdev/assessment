package main

import (
	"github.com/kraisitdev/assessment/app/installer"
	"github.com/labstack/echo/v4"
)

func init() {
	installer.SetupLogging()
}

func main() {
	e := echo.New()

	installer.SetupMiddleware(e)
	installer.SetupEndPoint(e)
	installer.SetupServer(e)
}
