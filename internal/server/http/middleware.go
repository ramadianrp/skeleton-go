package httpserver

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	appLog "github.com/ramadianrp/skeleton-go/internal/pkg/log"
	"github.com/ramadianrp/skeleton-go/internal/server/handler"
)

const maxLoggedBodyBytes = 32 * 1024

func SetupMiddleware(app *echo.Echo) {
	app.Use(Recover, handler.RequestContext, RequestLogger)
}

func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestBody, requestTruncated := captureRequestBody(c.Request(), maxLoggedBodyBytes)
		originalWriter := c.Response().Writer
		responseBody := &responseCapture{ResponseWriter: originalWriter, limit: maxLoggedBodyBytes}
		c.Response().Writer = responseBody

		err := next(c)
		c.Response().Writer = originalWriter
		ctx := c.Request().Context()
		status := c.Response().Status
		if err != nil {
			ctx = appLog.SetErrorMessage(ctx, err.Error())
			var httpError *echo.HTTPError
			if errors.As(err, &httpError) {
				status = httpError.Code
			} else if status < http.StatusBadRequest {
				status = http.StatusInternalServerError
			}
		}
		ctx = appLog.SetResponseMetadata(ctx, status, c.Response().Size, requestTruncated, responseBody.truncated)
		appLog.TDR(ctx, requestBody, responseBody.body.Bytes())
		return err
	}
}

func captureRequestBody(request *http.Request, limit int) ([]byte, bool) {
	if request.Body == nil {
		return nil, false
	}

	originalBody := request.Body
	captured, err := io.ReadAll(io.LimitReader(originalBody, int64(limit+1)))
	request.Body = &replayBody{
		Reader: io.MultiReader(bytes.NewReader(captured), originalBody),
		Closer: originalBody,
	}
	truncated := len(captured) > limit || err != nil
	if len(captured) > limit {
		captured = captured[:limit]
	}
	return captured, truncated
}

type replayBody struct {
	io.Reader
	io.Closer
}

type responseCapture struct {
	http.ResponseWriter
	body      bytes.Buffer
	limit     int
	truncated bool
}

func (capture *responseCapture) Write(payload []byte) (int, error) {
	written, err := capture.ResponseWriter.Write(payload)
	bytesToCapture := written
	remaining := capture.limit - capture.body.Len()
	if bytesToCapture > remaining {
		capture.truncated = true
		bytesToCapture = remaining
	}
	if bytesToCapture > 0 {
		_, _ = capture.body.Write(payload[:bytesToCapture])
	}
	return written, err
}

func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				slog.Error("panic recovered", "panic", recovered)
				err = echo.NewHTTPError(http.StatusInternalServerError)
			}
		}()
		return next(c)
	}
}
