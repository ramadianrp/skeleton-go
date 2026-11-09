package log

import (
	"strings"
	"testing"
)

func TestFormatLogsUsesKeyValuePairs(t *testing.T) {
	fields := formatLogs("method", "GET", "status", 200)
	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}
	if fields[0].Key != "method" || fields[0].Value.String() != "GET" {
		t.Fatalf("unexpected first field: %#v", fields[0])
	}
	if fields[1].Key != "status" || fields[1].Value.Int64() != 200 {
		t.Fatalf("unexpected second field: %#v", fields[1])
	}
}

func TestFormatLogsFallsBackForPositionalValues(t *testing.T) {
	fields := formatLogs("message", "extra", 42)
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
	wantKeys := []string{"_message_0", "_message_1", "_message_2"}
	for index, field := range fields {
		if field.Key != wantKeys[index] {
			t.Fatalf("unexpected field key at index %d: %q", index, field.Key)
		}
	}
}

func TestSafeJSONPayloadRedactsSensitiveFields(t *testing.T) {
	got := safeJSONPayload([]byte(`{"name":"Rani","email":"private@example.test","access_token":"secret-token","nested":{"password":"secret-password"}}`))
	for _, sensitive := range []string{"private@example.test", "secret-token", "secret-password"} {
		if strings.Contains(got, sensitive) {
			t.Fatalf("sensitive value %q was not redacted: %s", sensitive, got)
		}
	}
	if !strings.Contains(got, "[REDACTED]") {
		t.Fatalf("expected redaction markers, got %s", got)
	}
}
