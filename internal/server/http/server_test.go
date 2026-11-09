package httpserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ramadianrp/skeleton-go/internal/pkg/constants"
	"github.com/ramadianrp/skeleton-go/internal/server/handler"
	"github.com/ramadianrp/skeleton-go/internal/usecase/healthcheck"
	"github.com/ramadianrp/skeleton-go/internal/usecase/user"
)

func TestHealthCheck(t *testing.T) {
	app := New(handler.New(healthcheck.NewService(), nil, "", ""))
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	app.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	var body healthcheck.Response
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Message != "Server up and running" || body.Version == "" || body.ServerTime == "" {
		t.Fatalf("unexpected health response: %#v", body)
	}
}

func TestUserRoutesRequireBasicAuth(t *testing.T) {
	app := New(handler.New(healthcheck.NewService(), testUserService{}, "api", "secret"))
	request := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
	response := httptest.NewRecorder()

	app.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
	if response.Header().Get(echo.HeaderWWWAuthenticate) == "" {
		t.Fatal("expected WWW-Authenticate header")
	}
}

func TestCreateUserRejectsInvalidRequest(t *testing.T) {
	app := New(handler.New(healthcheck.NewService(), testUserService{}, "api", "secret"))
	request := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(`{"name":"ab","email":"bad"}`))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.SetBasicAuth("api", "secret")
	response := httptest.NewRecorder()

	app.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestCaptureRequestBodyPreservesHandlerBody(t *testing.T) {
	const body = `{"name":"long request body"}`
	request := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))

	captured, truncated := captureRequestBody(request, 8)
	forwarded, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatalf("read forwarded request body: %v", err)
	}
	if string(forwarded) != body {
		t.Fatalf("handler body was changed: got %q, want %q", forwarded, body)
	}
	if string(captured) != body[:8] || !truncated {
		t.Fatalf("expected bounded capture, got %q with truncated=%v", captured, truncated)
	}
}

func TestResponseCaptureDoesNotTruncateClientResponse(t *testing.T) {
	response := httptest.NewRecorder()
	capture := &responseCapture{ResponseWriter: response, limit: 4}

	written, err := capture.Write([]byte("full response"))
	if err != nil {
		t.Fatalf("write response: %v", err)
	}
	if written != len("full response") || response.Body.String() != "full response" {
		t.Fatalf("client response was truncated: wrote %d bytes, body %q", written, response.Body.String())
	}
	if capture.body.String() != "full" || !capture.truncated {
		t.Fatalf("expected capped log body, got %q with truncated=%v", capture.body.String(), capture.truncated)
	}
}

type testUserService struct{}

func (testUserService) Create(context.Context, user.CreateRequest) (user.CreateResponse, error) {
	return constants.DefaultResponse{}, nil
}

func (testUserService) FindAll(context.Context, user.FindAllRequest) (user.FindAllResponse, error) {
	return constants.DefaultResponse{}, nil
}
