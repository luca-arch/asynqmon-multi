package server_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/luca-arch/asynqmon-multi/server"
)

const (
	delayAfterStart = time.Second
	maxReqTime      = time.Second
)

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

func mkGetResponse(t *testing.T, url string) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("preparing %s %s: %v", http.MethodGet, url, err)

		return nil
	}

	c := &http.Client{
		Timeout: maxReqTime,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := c.Do(req)
	if err != nil {
		t.Fatalf("performing %s %s: %v", http.MethodGet, url, err)

		return nil
	}

	return res
}

func startTestServer(t *testing.T, opts *server.Options) {
	t.Helper()

	var err error

	go func() {
		err = server.Serve(t.Context(), opts)
	}()

	time.Sleep(delayAfterStart)

	if err != nil {
		t.Fatalf("could not start test server: %v", err)
	}
}
