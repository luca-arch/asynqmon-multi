package server

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

// OptionsProducer defines the HTTP server options.
type OptionsProducer interface {
	Address() string
	AsynqQueues() map[string]Queue
	HTTPHandler(opts asynqmon.Options) *asynqmon.HTTPHandler
	ShutdownDuration() time.Duration
	Template() string
}

// Options defines the configuration options for the asynqmon-multi server. It implements
// the [OptionsProducer] interface.
type Options struct {
	// Address to listen on, defaults to [DefaultAddr].
	Addr string `json:"address"`
	// HTML template to use for the index page, defaults to the content of default_index.tpl.
	HTMLTemplate string `json:"template"`
	// Port to listen on, defaults to [DefaultPort].
	Port uint `json:"port"`
	// Map of arbitrary queue names to their Redis connection options. Required.
	Queues map[string]Queue `json:"queues"`
	// Seconds to wait for the server to gracefully shut down when receiving a termination signal,
	// defaults to [DefaultShutdownDuration].
	ShutdownTimeout int `json:"shutdownTimeout"`
}

// Address returns the address and port to listen on, combining the [Options] Addr and Port fields.
func (o *Options) Address() string {
	var (
		addr = o.Addr
		port = o.Port
	)

	if addr == "" {
		addr = DefaultAddr
	}

	if port == 0 {
		port = DefaultPort
	}

	return fmt.Sprintf("%s:%d", addr, port)
}

// AsynqQueues returns the map of queue names to their Redis connection options.
func (o *Options) AsynqQueues() map[string]Queue {
	return o.Queues
}

// HTTPHandler is invoked by [Serve] to create a HTTPHandler with the given options.
func (o *Options) HTTPHandler(opts asynqmon.Options) *asynqmon.HTTPHandler {
	return asynqmon.New(opts)
}

// ShutdownDuration returns the duration to wait for the server to gracefully shut down.
func (o *Options) ShutdownDuration() time.Duration {
	if o.ShutdownTimeout < 1 {
		return DefaultShutdownDuration
	}

	return time.Duration(o.ShutdownTimeout) * time.Second
}

// Template returns the HTML template to use for the index page.
func (o *Options) Template() string {
	if o.HTMLTemplate == "" {
		return defaultTemplate
	}

	return o.HTMLTemplate
}

// Queue defines the [asynq.RedisClientOpt] values of a specific queue.
type Queue struct {
	ReadOnly       bool                  `json:"readOnly"`
	RedisAddr      string                `json:"redisAddr"`
	RedisClientOpt *asynq.RedisClientOpt `json:"-"`
	RedisDB        int                   `json:"redisDb"`
}

// RedisClientOption generates a [asynq.RedisClientOpt] for the queue.
func (q Queue) RedisClientOption() asynq.RedisClientOpt {
	out := asynq.RedisClientOpt{
		// Fixed options.
		Addr:     q.RedisAddr,
		DB:       q.RedisDB,
		PoolSize: 1,

		// Inherited options.
		DialTimeout:  0,
		Network:      "",
		Password:     "",
		ReadTimeout:  0,
		TLSConfig:    nil,
		WriteTimeout: 0,
		Username:     "",
	}

	if q.RedisClientOpt == nil {
		return out
	}

	out.DialTimeout = q.RedisClientOpt.DialTimeout
	out.Network = q.RedisClientOpt.Network
	out.Password = q.RedisClientOpt.Password
	out.ReadTimeout = q.RedisClientOpt.ReadTimeout
	out.TLSConfig = q.RedisClientOpt.TLSConfig
	out.WriteTimeout = q.RedisClientOpt.WriteTimeout
	out.Username = q.RedisClientOpt.Username

	return out
}
