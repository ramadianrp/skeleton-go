package log

import (
	"context"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	CTX_XID                = "ctx-xid"
	CTX_JID                = "ctx-jid"
	CTX_REQUEST_IP         = "ctx-request-ip"
	CTX_REQUEST_URL        = "ctx-request-url"
	CTX_REQUEST_METHOD     = "ctx-request-method"
	CTX_REQUEST_SIZE       = "ctx-request-size"
	CTX_REQUEST_TIME       = "ctx-request-time"
	CTX_REQUEST_HEADER     = "ctx-request-header"
	CTX_ERROR_MESSAGE      = "ctx-error-message"
	CTX_RESPONSE_STATUS    = "ctx-response-status"
	CTX_RESPONSE_SIZE      = "ctx-response-size"
	CTX_REQUEST_TRUNCATED  = "ctx-request-truncated"
	CTX_RESPONSE_TRUNCATED = "ctx-response-truncated"
)

func SetContextFromEchoRequest(c echo.Context) context.Context {
	ctx := c.Request().Context()
	headers := c.Request().Header.Clone()
	for key := range headers {
		switch strings.ToLower(key) {
		case "authorization", "proxy-authorization", "cookie", "set-cookie", "x-api-key":
			headers.Del(key)
		}
	}

	ctx = context.WithValue(ctx, CTX_REQUEST_IP, c.RealIP())
	ctx = context.WithValue(ctx, CTX_REQUEST_URL, c.Request().URL.Path)
	ctx = context.WithValue(ctx, CTX_REQUEST_METHOD, c.Request().Method)
	ctx = context.WithValue(ctx, CTX_REQUEST_SIZE, c.Request().ContentLength)
	ctx = context.WithValue(ctx, CTX_XID, c.Request().Header.Get("X-Request-ID"))
	ctx = context.WithValue(ctx, CTX_JID, c.Request().Header.Get("X-Journey-ID"))
	ctx = context.WithValue(ctx, CTX_REQUEST_TIME, time.Now())
	ctx = context.WithValue(ctx, CTX_REQUEST_HEADER, headers)

	return ctx
}

func SetResponseMetadata(ctx context.Context, status int, size int64, requestTruncated, responseTruncated bool) context.Context {
	ctx = context.WithValue(ctx, CTX_RESPONSE_STATUS, status)
	ctx = context.WithValue(ctx, CTX_RESPONSE_SIZE, size)
	ctx = context.WithValue(ctx, CTX_REQUEST_TRUNCATED, requestTruncated)
	return context.WithValue(ctx, CTX_RESPONSE_TRUNCATED, responseTruncated)
}

func SetErrorMessageFromEchoContext(c echo.Context, errMessage string) context.Context {
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, CTX_ERROR_MESSAGE, errMessage)
	return ctx
}

func SetErrorMessage(ctx context.Context, errMessage string) context.Context {
	ctx = context.WithValue(ctx, CTX_ERROR_MESSAGE, errMessage)
	return ctx
}

func GetRequestIPFromContext(ctx context.Context) string {
	s, ok := ctx.Value(CTX_REQUEST_IP).(string)
	if !ok {
		return ""
	}

	return s
}

func GetRequestURLFromContext(ctx context.Context) string {
	value, _ := ctx.Value(CTX_REQUEST_URL).(string)
	return value
}

func GetRequestMethodFromContext(ctx context.Context) string {
	value, _ := ctx.Value(CTX_REQUEST_METHOD).(string)
	return value
}

func GetRequestTimeFromContext(ctx context.Context) time.Time {
	s, ok := ctx.Value(CTX_REQUEST_TIME).(time.Time)
	if !ok {
		return time.Now()
	}

	return s
}

func GetRequestHeaderFromContext(ctx context.Context) interface{} {
	return ctx.Value(CTX_REQUEST_HEADER)
}

func GetRequestSizeFromContext(ctx context.Context) int64 {
	value, _ := ctx.Value(CTX_REQUEST_SIZE).(int64)
	return value
}

func GetResponseStatusFromContext(ctx context.Context) int {
	value, _ := ctx.Value(CTX_RESPONSE_STATUS).(int)
	return value
}

func GetResponseSizeFromContext(ctx context.Context) int64 {
	value, _ := ctx.Value(CTX_RESPONSE_SIZE).(int64)
	return value
}

func GetRequestTruncatedFromContext(ctx context.Context) bool {
	value, _ := ctx.Value(CTX_REQUEST_TRUNCATED).(bool)
	return value
}

func GetResponseTruncatedFromContext(ctx context.Context) bool {
	value, _ := ctx.Value(CTX_RESPONSE_TRUNCATED).(bool)
	return value
}

func GetErrorMessageFromContext(ctx context.Context) string {
	s, ok := ctx.Value(CTX_ERROR_MESSAGE).(string)
	if !ok {
		return ""
	}

	return s
}
