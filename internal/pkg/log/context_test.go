package log

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestSetContextFromEchoRequestRedactsSensitiveHeaders(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set(echo.HeaderAuthorization, "Bearer secret")
	request.Header.Set(echo.HeaderCookie, "session=secret")
	request.Header.Set("X-Request-ID", "request-123")
	context := echo.New().NewContext(request, httptest.NewRecorder())

	ctx := SetContextFromEchoRequest(context)
	headers := GetRequestHeaderFromContext(ctx).(http.Header)
	if headers.Get(echo.HeaderAuthorization) != "" || headers.Get(echo.HeaderCookie) != "" {
		t.Fatalf("sensitive headers were not redacted: %#v", headers)
	}
	if got := GetRequestURLFromContext(ctx); got != "/health" {
		t.Fatalf("expected request URL /health, got %q", got)
	}
	if got, ok := ctx.Value(CTX_XID).(string); !ok || got != "request-123" {
		t.Fatalf("expected request ID request-123, got %q", got)
	}
}
