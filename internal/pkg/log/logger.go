package log

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ramadianrp/skeleton-go/internal/config"
)

var (
	sysLogger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	tdrLogger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
)

func New() {
	sysLogger = newLogger("sys")
	tdrLogger = newLogger("tdr")
}

func SystemLogger() *slog.Logger {
	return sysLogger
}

func newLogger(name string) *slog.Logger {
	var writers []io.Writer

	if config.GetBool("logger." + name + ".stdout") {
		writers = append(writers, os.Stdout)
	}

	if location := config.GetString("logger." + name + ".fileLocation"); location != "" {
		if err := os.MkdirAll(filepath.Dir(location), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "create log directory: %v\n", err)
		} else if file, err := os.OpenFile(location, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "open log file: %v\n", err)
		} else {
			writers = append(writers, file)
		}
	}

	if len(writers) == 0 {
		writers = append(writers, os.Stderr)
	}

	return slog.New(slog.NewJSONHandler(io.MultiWriter(writers...), nil)).
		With("app", config.GetString("appName"))
}

func Info(ctx context.Context, title string, messages ...interface{}) {
	sysLogger.LogAttrs(ctx, slog.LevelInfo, title, formatLogs(messages...)...)
}

func Warn(ctx context.Context, title string, messages ...interface{}) {
	sysLogger.LogAttrs(ctx, slog.LevelWarn, title, formatLogs(messages...)...)
}

func Error(ctx context.Context, title string, messages ...interface{}) {
	sysLogger.LogAttrs(ctx, slog.LevelError, title, formatLogs(messages...)...)
}

func Fatal(ctx context.Context, title string, messages ...interface{}) {
	sysLogger.LogAttrs(ctx, slog.LevelError, title, formatLogs(messages...)...)
	os.Exit(1)
}

func T2(ctx context.Context, title string, messages ...interface{}) {
	Info(ctx, title, messages...)
}

func T3(ctx context.Context, title string, startProcessTime time.Time, messages ...interface{}) {
	fields := formatLogs(messages...)
	fields = append(fields, slog.String(
		"_process_time",
		fmt.Sprintf("%d ms", time.Since(startProcessTime).Milliseconds()),
	))
	sysLogger.LogAttrs(ctx, slog.LevelInfo, title, fields...)
}

func TDR(ctx context.Context, request, response []byte) {
	fields := []slog.Attr{
		slog.String("app_name", config.GetString("appName")),
		slog.String("src_ip", GetRequestIPFromContext(ctx)),
		slog.String("request_method", GetRequestMethodFromContext(ctx)),
		slog.String("request_url", GetRequestURLFromContext(ctx)),
		slog.Int64("response_time_ms", time.Since(GetRequestTimeFromContext(ctx)).Milliseconds()),
		slog.Any("header", GetRequestHeaderFromContext(ctx)),
		slog.Int64("request_size_bytes", GetRequestSizeFromContext(ctx)),
		slog.Int64("response_size_bytes", GetResponseSizeFromContext(ctx)),
		slog.Int("http_status", GetResponseStatusFromContext(ctx)),
		slog.Bool("request_body_truncated", GetRequestTruncatedFromContext(ctx)),
		slog.Bool("response_body_truncated", GetResponseTruncatedFromContext(ctx)),
		slog.Any("error", GetErrorMessageFromContext(ctx)),
	}
	if len(request) > 0 {
		fields = append(fields, slog.String("request", safeJSONPayload(request)))
	}
	if len(response) > 0 {
		fields = append(fields, slog.String("response", safeJSONPayload(response)))
	}
	if requestID, ok := ctx.Value(CTX_XID).(string); ok && requestID != "" {
		fields = append(fields, slog.String("request_id", requestID))
	}
	if journeyID, ok := ctx.Value(CTX_JID).(string); ok && journeyID != "" {
		fields = append(fields, slog.String("journey_id", journeyID))
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(response, &responseMap); err == nil {
		if status, ok := responseMap["status"]; ok && status != nil {
			fields = append(fields, slog.String("response_code", fmt.Sprint(status)))
		}
	}
	tdrLogger.LogAttrs(ctx, slog.LevelInfo, "TDR", fields...)
}

func safeJSONPayload(payload []byte) string {
	var value interface{}
	if err := json.Unmarshal(payload, &value); err != nil {
		return "<non-JSON body omitted>"
	}
	encoded, err := json.Marshal(redactSensitiveValue(value))
	if err != nil {
		return "<body omitted>"
	}
	return string(encoded)
}

func redactSensitiveValue(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		redacted := make(map[string]interface{}, len(typed))
		for key, item := range typed {
			if isSensitiveKey(key) {
				redacted[key] = "[REDACTED]"
			} else {
				redacted[key] = redactSensitiveValue(item)
			}
		}
		return redacted
	case []interface{}:
		redacted := make([]interface{}, len(typed))
		for index, item := range typed {
			redacted[index] = redactSensitiveValue(item)
		}
		return redacted
	default:
		return value
	}
}

func isSensitiveKey(key string) bool {
	normalized := strings.ToLower(strings.NewReplacer("_", "", "-", "").Replace(key))
	for _, sensitive := range []string{"password", "secret", "token", "authorization", "cookie", "apikey", "email", "phone", "accountnumber"} {
		if strings.Contains(normalized, sensitive) {
			return true
		}
	}
	return false
}

func formatLogs(messages ...interface{}) []slog.Attr {
	if len(messages) > 0 && len(messages)%2 == 0 {
		fields := make([]slog.Attr, 0, len(messages)/2)
		for index := 0; index < len(messages); index += 2 {
			key, ok := messages[index].(string)
			if !ok || key == "" {
				fields = nil
				break
			}
			fields = append(fields, slog.Any(key, messages[index+1]))
		}
		if fields != nil {
			return fields
		}
	}

	fields := make([]slog.Attr, 0, len(messages))
	for index, message := range messages {
		fields = append(fields, slog.Any(fmt.Sprintf("_message_%d", index), message))
	}
	return fields
}
