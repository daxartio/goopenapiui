package main

import (
	_ "embed"
	"net/http"

	"github.com/daxartio/goopenapiui"
)

//go:embed openapi.json
var openapiJson []byte

func main() {
	openapiui := &goopenapiui.OpenapiUI{
		Title:       "Example API",
		Description: "Example API Description",
		Openapi:     openapiJson,
		OpenapiURL:  "/openapi.json",
		SwaggerURL:  "/docs",
	}

	println("Documentation served at http://127.0.0.1:8000/docs")
	panic(http.ListenAndServe(":8000", openapiui.Handler()))
}
