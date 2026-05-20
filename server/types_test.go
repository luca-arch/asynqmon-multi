package server_test

import (
	_ "embed"
	"maps"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/luca-arch/asynqmon-multi/server"
)

//go:embed default_index.tpl
var defaultTemplate string

func TestOptions(t *testing.T) {
	t.Parallel()

	type want struct {
		addr     string
		queues   map[string]server.Queue
		shutdown time.Duration
		template string
	}

	tests := map[string]struct {
		opts server.Options
		want want
	}{
		"default values": {
			opts: server.Options{
				Addr:            "",
				HTMLTemplate:    "",
				Port:            0,
				Queues:          mkTestQueues(t),
				ShutdownTimeout: 0,
			},
			want: want{
				addr:     "0.0.0.0:8080",
				queues:   mkTestQueues(t),
				shutdown: 10 * time.Second,
				template: defaultTemplate,
			},
		},
		"custom values": {
			opts: server.Options{
				Addr:            "localhost",
				HTMLTemplate:    "<div></div>",
				Port:            9090,
				Queues:          mkTestQueues(t),
				ShutdownTimeout: 60,
			},
			want: want{
				addr:     "localhost:9090",
				queues:   mkTestQueues(t),
				shutdown: time.Minute,
				template: "<div></div>",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tc.opts.Address() != tc.want.addr {
				t.Errorf("Address(): got %s, want %s", tc.opts.Address(), tc.want.addr)
			}

			if !maps.Equal(tc.opts.AsynqQueues(), tc.want.queues) {
				t.Errorf("AsynqQueues(): got %v, want %v", tc.opts.AsynqQueues(), tc.want.queues)
			}

			if tc.opts.ShutdownDuration() != tc.want.shutdown {
				t.Errorf("ShutdownDuration(): got %s, want %s", tc.opts.ShutdownDuration(), tc.want.shutdown)
			}

			if tc.opts.Template() != tc.want.template {
				t.Errorf("Template(): got %s, want %s", tc.opts.Template(), tc.want.template)
			}

			if tc.opts.HTTPHandler(asynqmon.Options{RedisConnOpt: asynq.RedisClientOpt{}}) == nil {
				t.Error("HTTPHandler(): got nil")
			}
		})
	}
}

func TestRedisClientOption(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		queue server.Queue
		want  asynq.RedisClientOpt
	}{
		"with inherited options": {
			queue: server.Queue{
				ReadOnly:  true,
				RedisAddr: "localhost:16379",
				RedisDB:   2,
				RedisClientOpt: &asynq.RedisClientOpt{
					DialTimeout:  time.Second,
					Network:      "test-networw",
					Password:     "password",
					ReadTimeout:  2 * time.Second,
					WriteTimeout: 4 * time.Second,
					Username:     "user",
				},
			},
			want: asynq.RedisClientOpt{
				Addr:         "localhost:16379",
				DB:           2,
				DialTimeout:  time.Second,
				Network:      "test-networw",
				Password:     "password",
				PoolSize:     1,
				ReadTimeout:  2 * time.Second,
				WriteTimeout: 4 * time.Second,
				Username:     "user",
			},
		},
		"without inherited options": {
			queue: server.Queue{
				ReadOnly:       true,
				RedisAddr:      "localhost:26379",
				RedisDB:        3,
				RedisClientOpt: nil,
			},
			want: asynq.RedisClientOpt{
				Addr:         "localhost:26379",
				DB:           3,
				DialTimeout:  0,
				Network:      "",
				Password:     "",
				PoolSize:     1,
				ReadTimeout:  0,
				TLSConfig:    nil,
				WriteTimeout: 0,
				Username:     "",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tc.queue.RedisClientOption() != tc.want {
				t.Errorf("RedisClientOption(): got %v, want %v", tc.queue.RedisClientOption(), tc.want)
			}
		})
	}
}
