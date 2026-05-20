package compliance

import (
	"strings"
	"testing"
)

func TestParseJSONUsesBoundedPrefixDecode(t *testing.T) {
	raw, ok := parseJSON("validator log prefix\n{\"error_count\":2,\"warning_count\":1}\ntrailing log line")
	if !ok {
		t.Fatalf("parseJSON did not recover JSON prefix")
	}
	object, ok := raw.(map[string]any)
	if !ok || numberField(object, "error_count") != 2 || numberField(object, "warning_count") != 1 {
		t.Fatalf("parsed object = %#v, want counts", raw)
	}

	hugeMalformed := "validator log " + strings.Repeat("{", 20000) + strings.Repeat("x", 20000)
	if raw, ok := parseJSON(hugeMalformed); ok || raw != nil {
		t.Fatalf("malformed huge JSON parsed unexpectedly: %#v", raw)
	}
}
