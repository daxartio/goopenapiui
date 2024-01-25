package goopenapiui

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

const (
	OpenapiURL        = "/openapi.json"
	SwaggerURL        = "/docs"
	SwaggerjsURL      = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"
	SwaggercssURL     = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.9.0/swagger-ui.css"
	SwaggerFaviconURL = "https://static1.smartbear.co/swagger/media/assets/swagger_fav.png"
)

//go:embed assets/swagger-ui.html
var SwaggerHTML string

var (
	ErrRenderSwaggerUI = errors.New("render swagger ui failed")
	ErrParseSwaggerUI  = errors.New("parse swagger ui failed")
)

type OpenapiUI struct {
	Title             string
	Description       string
	Openapi           []byte
	OpenapiURL        string
	SwaggerURL        string
	SwaggerjsURL      string
	SwaggercssURL     string
	SwaggerFaviconURL string
}

func (o *OpenapiUI) SwaggerUI() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	tpl, err := template.New("swaggerui").Parse(SwaggerHTML)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParseSwaggerUI, err)
	}

	if o.SwaggerjsURL == "" {
		o.SwaggerjsURL = SwaggerjsURL
	}

	if o.SwaggercssURL == "" {
		o.SwaggercssURL = SwaggercssURL
	}

	if o.SwaggerFaviconURL == "" {
		o.SwaggerFaviconURL = SwaggerFaviconURL
	}

	if err = tpl.Execute(buf, map[string]string{
		"title":       o.Title,
		"description": o.Description,
		"openapiurl":  o.OpenapiURL,
		"jsurl":       o.SwaggerjsURL,
		"cssurl":      o.SwaggercssURL,
		"iconurl":     o.SwaggerFaviconURL,
	}); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRenderSwaggerUI, err)
	}

	return buf.Bytes(), nil
}

// Handler sets some defaults and returns a HandlerFunc.
func (o *OpenapiUI) Handler() http.HandlerFunc {
	html, err := o.SwaggerUI()
	if err != nil {
		panic(err)
	}

	var openapiURL string
	if o.OpenapiURL == "" {
		openapiURL = OpenapiURL
	} else {
		openapiURL = o.OpenapiURL
	}

	var swaggerURL string
	if o.SwaggerURL == "" {
		swaggerURL = SwaggerURL
	} else {
		swaggerURL = o.SwaggerURL
	}

	openapi := o.Openapi

	return func(writer http.ResponseWriter, req *http.Request) {
		method := strings.ToLower(req.Method)
		if method != "get" && method != "head" {
			return
		}

		header := writer.Header()
		if req.URL.Path == openapiURL {
			header.Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write(openapi)

			return
		}

		if req.URL.Path == swaggerURL {
			header.Set("Content-Type", "text/html")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write(html)
		}
	}
}
