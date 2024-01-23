package goopenapiui

import (
	"bytes"
	_ "embed"
	"net/http"
	"strings"
	"text/template"
)

const OpenapiUrl = "/openapi.json"
const SwaggerUrl = "/docs"
const SwaggerJSUrl = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"
const SwaggerCSSUrl = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.9.0/swagger-ui.css"
const SwaggerFaviconUrl = "https://static1.smartbear.co/swagger/media/assets/swagger_fav.png"

//go:embed assets/swagger-ui.html
var SwaggerHTML string

type Openapiui struct {
	Title             string
	Description       string
	Openapi           []byte
	OpenapiUrl        string
	SwaggerUrl        string
	SwaggerJSUrl      string
	SwaggerCSSUrl     string
	SwaggerFaviconUrl string
}

func (o *Openapiui) SwaggerUI() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	tpl, err := template.New("swaggerui").Parse(SwaggerHTML)
	if err != nil {
		return nil, err
	}

	if o.SwaggerJSUrl == "" {
		o.SwaggerJSUrl = SwaggerJSUrl
	}

	if o.SwaggerCSSUrl == "" {
		o.SwaggerCSSUrl = SwaggerCSSUrl
	}

	if o.SwaggerFaviconUrl == "" {
		o.SwaggerFaviconUrl = SwaggerFaviconUrl
	}

	if err = tpl.Execute(buf, map[string]string{
		"title":       o.Title,
		"description": o.Description,
		"openapiurl":  o.OpenapiUrl,
		"jsurl":       o.SwaggerJSUrl,
		"cssurl":      o.SwaggerCSSUrl,
		"iconurl":     o.SwaggerFaviconUrl,
	}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Handler sets some defaults and returns a HandlerFunc
func (r *Openapiui) Handler() http.HandlerFunc {
	html, err := r.SwaggerUI()
	if err != nil {
		panic(err)
	}

	var openapiUrl string
	if r.OpenapiUrl == "" {
		openapiUrl = OpenapiUrl
	} else {
		openapiUrl = r.OpenapiUrl
	}

	var swaggerUrl string
	if r.SwaggerUrl == "" {
		swaggerUrl = SwaggerUrl
	} else {
		swaggerUrl = r.SwaggerUrl
	}

	openapi := r.Openapi

	return func(w http.ResponseWriter, req *http.Request) {
		method := strings.ToLower(req.Method)
		if method != "get" && method != "head" {
			return
		}

		header := w.Header()
		if req.URL.Path == openapiUrl {
			header.Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(openapi)
			return
		}

		if req.URL.Path == swaggerUrl {
			header.Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(html)
		}
	}
}
