package api

import (
	"encoding/json"
	"io"
	"strings"
	"testing"
)

func TestGenerateBulkFetchIssuesBody(t *testing.T) {
	issueKeys := []string{"PROJ-1", "PROJ-2"}
	buf, err := generateBulkFetchIssuesBody(issueKeys)
	if err != nil {
		t.Fatalf("failed to generate body: %v", err)
	}
	bodyBytes, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read buffer: %v", err)
	}

	// Check that the buffer contains the expected keys
	body := string(bodyBytes)

	if !strings.Contains(body, "issueIdsOrKeys") || !strings.Contains(body, "fields") {
		t.Errorf("body missing expected keys: %s", body)
	}

	// Try to unmarshal to verify valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		t.Errorf("body is not valid JSON: %v", err)
	}

	// Check that issueIdsOrKeys matches input
	ids, ok := parsed["issueIdsOrKeys"].([]interface{})
	if !ok || len(ids) != 2 {
		t.Errorf("issueIdsOrKeys not correct: %v", parsed["issueIdsOrKeys"])
	}
}
