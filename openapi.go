package goopenapiui

import (
	"bytes"
	_ "embed"
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
		return nil, err
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
		return nil, err
	}

	return buf.Bytes(), nil
}

// Handler sets some defaults and returns a HandlerFunc
func (r *OpenapiUI) Handler() http.HandlerFunc {
	html, err := r.SwaggerUI()
	if err != nil {
		panic(err)
	}

	var openapiURL string
	if r.OpenapiURL == "" {
		openapiURL = OpenapiURL
	} else {
		openapiURL = r.OpenapiURL
	}

	var swaggerURL string
	if r.SwaggerURL == "" {
		swaggerURL = SwaggerURL
	} else {
		swaggerURL = r.SwaggerURL
	}

	openapi := r.Openapi

	return func(w http.ResponseWriter, req *http.Request) {
		method := strings.ToLower(req.Method)
		if method != "get" && method != "head" {
			return
		}

		header := w.Header()
		if req.URL.Path == openapiURL {
			header.Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(openapi)
			return
		}

		if req.URL.Path == swaggerURL {
			header.Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(html)
		}
	}
}
