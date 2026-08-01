package server_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/luca-arch/asynqmon-multi/server"
)

func TestMain(t *testing.T) {
	t.Parallel()

	validQueues := mkTestFile(t, `{"db1": {"redisDb": 1}, "db2": {"redisDb": 2}}`)

	tests := map[string]struct {
		args []string
		want error
	}{
		"no arguments": {
			args: []string{},
			want: server.ErrNoQueue,
		},
		"bad json in --queues-file": {
			args: []string{
				"--queues-file",
				mkTestFile(t, `{"q": false}`),
			},
			want: server.ErrInvalidFile,
		},
		"empty --queues-file": {
			args: []string{
				"--queues-file",
				mkTestFile(t, `{}`),
			},
			want: server.ErrInvalidOption,
		},
		"non-existent --template-file": {
			args: []string{
				"--queues-file",
				validQueues,
				"--template-file",
				"/does/not/exist",
			},
			want: server.ErrCliOptions,
		},
		"bad template in --template-file": {
			args: []string{
				"--queues-file",
				validQueues,
				"--template-file",
				mkTestFile(t, "{{ call .Broken }}"),
			},
			want: server.ErrInvalidTemplate,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := server.Main(tc.args)

			if !errors.Is(got, tc.want) {
				t.Errorf("Main(): got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts *server.Options
		want error
	}{
		"no queues defined": {
			opts: &server.Options{
				Queues: nil,
			},
			want: server.ErrInvalidOption,
		},
		"duplicate queues defined": {
			opts: &server.Options{
				Queues: map[string]server.Queue{
					"q0": {
						ReadOnly:  false,
						RedisAddr: "redis",
						RedisDB:   1,
					},
					"q1": {
						ReadOnly:  true,
						RedisAddr: "redis",
						RedisDB:   1,
					},
				},
			},
			want: server.ErrInvalidOption,
		},
		"bad template": {
			opts: &server.Options{
				Queues:       mkTestQueues(t),
				HTMLTemplate: "{{ .Broken",
			},
			want: server.ErrInvalidTemplate,
		},
		"bad rendering": {
			opts: &server.Options{
				Queues:       mkTestQueues(t),
				HTMLTemplate: "{{ call .Queues }}",
			},
			want: server.ErrInvalidTemplate,
		},
		"ok": {
			opts: &server.Options{
				Queues: map[string]server.Queue{
					"q0": {
						ReadOnly:  false,
						RedisAddr: "redis:6379",
						RedisDB:   1,
					},
					"q1": {
						ReadOnly:  true,
						RedisAddr: "redis:6379",
						RedisDB:   2,
					},
					"q2": {
						ReadOnly:  false,
						RedisAddr: "redis:16379",
						RedisDB:   1,
					},
					"q3": {
						ReadOnly:  true,
						RedisAddr: "redis:16379",
						RedisDB:   2,
					},
				},
			},
			want: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, got := server.New(tc.opts)

			if !errors.Is(got, tc.want) {
				t.Errorf("New(): got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestServe(t *testing.T) {
	t.Parallel()

	startTestServer(t, &server.Options{
		Addr:            "127.0.0.1",
		Port:            64242,
		ShutdownTimeout: 1,
		Queues: map[string]server.Queue{
			"db1": {RedisDB: 1},
			"db5": {RedisDB: 5},
		},
	})

	tests := map[string]struct {
		endpoint   string
		wantHTML   []string
		wantStatus int
	}{
		"index page": {
			endpoint:   "/",
			wantStatus: http.StatusOK,
			wantHTML: []string{
				"<title>asynqmon-multi</title>",
				`<a href="/db1" target="asynqmon">db1 (db1)</a>`,
				`<a href="/db5" target="asynqmon">db5 (db5)</a>`,
				`<iframe name="asynqmon" />`,
			},
		},
		"db1 page": {
			endpoint:   "/db1/",
			wantStatus: http.StatusOK,
			wantHTML: []string{
				"<title>Asynq - Monitoring</title>",
				`<link rel="manifest" href="/db1/manifest.json"/>`,
			},
		},
		"db5 page": {
			endpoint:   "/db5/",
			wantStatus: http.StatusOK,
			wantHTML: []string{
				"<title>Asynq - Monitoring</title>",
				`<link rel="manifest" href="/db5/manifest.json"/>`,
			},
		},
		"db10 page": {
			endpoint:   "/db10/",
			wantStatus: http.StatusNotFound,
			wantHTML: []string{
				`404 page not found`,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := mkGetResponse(t, "http://127.0.0.1:64242"+tc.endpoint)
			defer got.Body.Close()

			if got.StatusCode != tc.wantStatus {
				t.Errorf("got %d, want %d", got.StatusCode, tc.wantStatus)

				return
			}

			if tc.wantHTML == nil {
				return
			}

			html, err := io.ReadAll(got.Body)
			if err != nil {
				t.Error(err)

				return
			}

			for _, tag := range tc.wantHTML {
				if !strings.Contains(string(html), tag) {
					t.Errorf("tag %s not found in %s", tag, html)

					return
				}
			}
		})
	}
}

func TestServeOne(t *testing.T) {
	t.Parallel()

	startTestServer(t, &server.Options{
		Addr:            "127.0.0.1",
		Port:            64343,
		ShutdownTimeout: 1,
		Queues: map[string]server.Queue{
			"only_db": {RedisDB: 1},
		},
	})

	tests := map[string]struct {
		endpoint   string
		wantHeader map[string]string
		wantStatus int
	}{
		"index page": {
			endpoint:   "/",
			wantHeader: map[string]string{"Location": "/only_db"},
			wantStatus: http.StatusFound,
		},
		"only_db page": {
			endpoint:   "/only_db/",
			wantStatus: http.StatusOK,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := mkGetResponse(t, "http://127.0.0.1:64343"+tc.endpoint)
			defer got.Body.Close()

			if got.StatusCode != tc.wantStatus {
				t.Errorf("got %d, want %d", got.StatusCode, tc.wantStatus)

				return
			}

			if tc.wantHeader == nil {
				return
			}

			for h, val := range tc.wantHeader {
				all, ok := got.Header[h]

				switch {
				case !ok:
					t.Errorf("header %s not found in response", h)
				case len(all) != 1:
					t.Errorf("wrong number of headers %s returned: got %d, want 1", h, len(all))
				case all[0] != val:
					t.Errorf("invalid header %s returned: got %s, want %s", h, all[0], val)
				}
			}
		})
	}
}

func TestGracefulShutdown(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 2_500*time.Millisecond)
	defer cancel()

	if err := server.Serve(ctx, &server.Options{
		Addr:   "127.0.0.1",
		Port:   64243,
		Queues: mkTestQueues(t),
	}); err != nil {
		t.Error(err)
	}
}
