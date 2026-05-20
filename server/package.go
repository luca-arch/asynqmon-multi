// Package server provides an HTTP server for asynqmon-multi that monitors multiple Asynq queues.
//
// It exports two main functions: [Main] for starting the server with CLI argument parsing and
// signal handling, and [Serve] for running the server with provided options and context.
//
// To instantiate a new HTTP server without starting it, use [New].
package server

import (
	_ "embed"
	"errors"
	"time"
)

// Default values for command line invocation.
const (
	DefaultAddr             = "0.0.0.0"
	DefaultPort             = 8080
	DefaultShutdownDuration = 10 * time.Second

	programName = "asynqmon-multi"
)

//go:embed default_index.tpl
var defaultTemplate string

// Errors returned by this package.
var (
	ErrCliOptions      = errors.New("parsing cli arguments")
	ErrInvalidFile     = errors.New("invalid file provided")
	ErrInvalidOption   = errors.New("invalid option")
	ErrInvalidTemplate = errors.New("invalid HTML template")
	ErrNoQueue         = errors.New("missing --queues-file")

	errRunServer    = errors.New("running server")
	errShuttingDown = errors.New("shutting down server")
)
