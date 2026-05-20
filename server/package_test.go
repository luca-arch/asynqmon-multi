package server_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/luca-arch/asynqmon-multi/server"
)

const baseURL = "http://127.0.0.1:64242"

func mkTestFile(t *testing.T, content string) string {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "asynqmon-multi.test*")
	if err != nil {
		t.Fatal(err)

		return ""
	}

	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)

		return ""
	}

	return f.Name()
}

func mkTestQueues(t *testing.T) map[string]server.Queue {
	t.Helper()

	return map[string]server.Queue{
		"q0": {
			ReadOnly:  false,
			RedisAddr: "redis:6379",
			RedisDB:   0,
		},
		"q1": {
			ReadOnly:  false,
			RedisAddr: "redis:6379",
			RedisDB:   1,
		},
	}
}

func mkGetResponse(t *testing.T, endpoint string) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, baseURL+endpoint, nil)
	if err != nil {
		t.Fatalf("preparing %s %s: %v", http.MethodGet, endpoint, err)

		return nil
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("performing %s %s: %v", http.MethodGet, endpoint, err)

		return nil
	}

	return res
}
