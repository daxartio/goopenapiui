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
		Title:             "Example API",
		Description:       "Example API Description",
		Openapi:           []byte(`{"swagger":"2.0"}`),
		OpenapiURL:        "/openapi.json",
		SwaggerURL:        "/docs",
		SwaggerjsURL:      goopenapiui.SwaggerjsURL,
		SwaggercssURL:     goopenapiui.SwaggercssURL,
		SwaggerFaviconURL: goopenapiui.SwaggerFaviconURL,
	}

	t.Run("Handler", func(t *testing.T) {
		handler := openapiui.Handler()

		t.Run("Spec", func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
			recorder := httptest.NewRecorder()
			handler(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}
			if resp.Header.Get("Content-Type") != "application/json" {
				t.Errorf("expected Content-Type %q, got %q", "application/json", resp.Header.Get("Content-Type"))
			}

			respBody := recorder.Body.String()
			if respBody != `{"swagger":"2.0"}` {
				t.Errorf("expected body %q, got %q", `{"swagger":"2.0"}`, respBody)
			}
		})

		t.Run("Docs", func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/docs", nil)
			recorder := httptest.NewRecorder()
			handler(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}
			if resp.Header.Get("Content-Type") != "text/html" {
				t.Errorf("expected Content-Type %q, got %q", "text/html", resp.Header.Get("Content-Type"))
			}

			respBody := recorder.Body.String()
			if respBody == "" {
				t.Errorf("expected body, got %q", respBody)
			}
		})
	})
}
