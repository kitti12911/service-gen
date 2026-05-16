package server

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRequest(target string) *http.Request {
	return httptest.NewRequestWithContext(context.Background(), http.MethodGet, target, http.NoBody)
}

func newTestHTTPServer() *HTTPServer {
	return NewHTTPServer(0, "___NAME___")
}

func TestHTTPServerHealth(t *testing.T) {
	srv := newTestHTTPServer()
	req := newRequest("/health")
	rec := httptest.NewRecorder()

	srv.server.Handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"ok"`)
}

func TestHTTPServerGzip(t *testing.T) {
	srv := newTestHTTPServer()
	req := newRequest("/openapi.json")
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()

	srv.server.Handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "gzip", rec.Header().Get("Content-Encoding"))
	assert.Contains(t, rec.Header().Values("Vary"), "Accept-Encoding")

	gz, err := gzip.NewReader(rec.Body)
	require.NoError(t, err)
	defer gz.Close()

	body, err := io.ReadAll(gz)
	require.NoError(t, err)
	assert.Contains(t, string(body), `"title":"___NAME___"`)
}

func TestHTTPServerOpenAPI(t *testing.T) {
	srv := newTestHTTPServer()
	req := newRequest("/openapi.json")
	rec := httptest.NewRecorder()

	srv.server.Handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"title":"___NAME___"`)
	assert.Contains(t, rec.Body.String(), `"/health"`)
}

func TestDocsServesSwaggerUIOfflineAndAllowsDownloads(t *testing.T) {
	srv := newTestHTTPServer()
	req := newRequest("/docs")
	rec := httptest.NewRecorder()

	srv.server.Handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	csp := rec.Header().Get("Content-Security-Policy")
	assert.Contains(t, csp, "allow-downloads")
	assert.NotContains(t, csp, "unpkg.com")

	body := rec.Body.String()
	assert.NotContains(t, body, "unpkg.com")
	assert.Contains(t, body, `href="/assets/swagger-ui/swagger-ui.css"`)
	assert.Contains(t, body, `href="/assets/swagger-ui/docs-overrides.css"`)
	assert.Contains(t, body, `src="/assets/swagger-ui/swagger-ui-bundle.js"`)
	assert.Contains(t, body, `src="/assets/swagger-ui/swagger-initializer.js"`)
	assert.Contains(t, body, `data-url="/openapi.json"`)
}

func TestSwaggerUIAssetsServed(t *testing.T) {
	srv := newTestHTTPServer()

	tests := []struct {
		path      string
		minLength int
	}{
		{path: "/assets/swagger-ui/swagger-ui.css", minLength: 100_000},
		{path: "/assets/swagger-ui/docs-overrides.css", minLength: 100},
		{path: "/assets/swagger-ui/swagger-ui-bundle.js", minLength: 1_000_000},
		{path: "/assets/swagger-ui/swagger-initializer.js", minLength: 100},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := newRequest(tt.path)
			rec := httptest.NewRecorder()

			srv.server.Handler.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.GreaterOrEqual(t, rec.Body.Len(), tt.minLength)
		})
	}
}

func TestOpenAPIDownload(t *testing.T) {
	srv := newTestHTTPServer()

	tests := []struct {
		path        string
		contentType string
		filename    string
		bodyMustHas string
	}{
		{
			path:        "/openapi.json/download",
			contentType: "application/openapi+json",
			filename:    "openapi.json",
			bodyMustHas: `"openapi"`,
		},
		{
			path:        "/openapi.yaml/download",
			contentType: "application/openapi+yaml",
			filename:    "openapi.yaml",
			bodyMustHas: "openapi:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := newRequest(tt.path)
			rec := httptest.NewRecorder()

			srv.server.Handler.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tt.contentType, rec.Header().Get("Content-Type"))
			assert.Equal(t, `attachment; filename="`+tt.filename+`"`, rec.Header().Get("Content-Disposition"))
			assert.True(t, strings.Contains(rec.Body.String(), tt.bodyMustHas))
		})
	}
}
