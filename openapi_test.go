package goopenapiui_test

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daxartio/goopenapiui"
)

func TestOpenapiui(t *testing.T) {
	openapiui := &goopenapiui.OpenapiUI{
		Title:       "Example API",
		Description: "Example API Description",
		Openapi:     []byte(`{"swagger":"2.0"}`),
		OpenapiURL:  "/openapi.json",
		SwaggerURL:  "/docs",
	}

	t.Run("Body", func(t *testing.T) {
		body, err := openapiui.SwaggerUI()
		if err != nil {
			t.Fatal(err)
		}
		_ = body
	})

	t.Run("Handler", func(t *testing.T) {
		handler := openapiui.Handler()

		t.Run("Spec", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
			w := httptest.NewRecorder()
			handler(w, req)

			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}
			if resp.Header.Get("Content-Type") != "application/json" {
				t.Errorf("expected Content-Type %q, got %q", "application/json", resp.Header.Get("Content-Type"))
			}
			// TODO: check body
		})

		t.Run("Docs", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/docs", nil)
			w := httptest.NewRecorder()
			handler(w, req)

			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}
			if resp.Header.Get("Content-Type") != "text/html" {
				t.Errorf("expected Content-Type %q, got %q", "text/html", resp.Header.Get("Content-Type"))
			}
			// TODO: check body
		})
	})
}
