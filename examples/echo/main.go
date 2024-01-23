package main

import (
	_ "embed"

	"github.com/labstack/echo/v4"

	"github.com/daxartio/goopenapiui"
	echoopenapiui "github.com/daxartio/goopenapiui/echo"
)

//go:embed openapi.json
var openapiJson []byte

func main() {
	openapiui := &goopenapiui.Openapiui{
		Title:       "Example API",
		Description: "Example API Description",
		Openapi:     openapiJson,
		OpenapiUrl:  "/openapi.json",
		SwaggerUrl:  "/docs",
	}
	e := echo.New()
	e.Use(echoopenapiui.New(openapiui))

	println("Documentation served at http://127.0.0.1:8000/docs")
	panic(e.Start(":8000"))
}
