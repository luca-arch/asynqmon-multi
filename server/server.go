package server

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynqmon"
	"golang.org/x/sync/errgroup"
)

// Main is the entry point for the asynqmon-multi server.
// It sets up signal handling for graceful shutdown and calls [Serve].
//
// The program [Options] are read from the provided command line arguments, or defaulted to
// [os.Args] if nil.
func Main(args []string) error {
	opts, err := optionsFromCli(args)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	return Serve(ctx, opts)
}

// New sets up a new [http.Server] without starting it.
func New(opts OptionsProducer) (*http.Server, error) {
	if err := validateOptions(opts.AsynqQueues()); err != nil {
		return nil, err
	}

	html, err := execTemplate(opts.Template(), opts.AsynqQueues())
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	for name, q := range sortedQueues(opts.AsynqQueues()) {
		h := opts.HTTPHandler(asynqmon.Options{
			PayloadFormatter:  nil,
			PrometheusAddress: "",
			ReadOnly:          q.ReadOnly,
			RedisConnOpt:      q.RedisClientOption(),
			ResultFormatter:   nil,
			RootPath:          "/" + name,
		})

		mux.Handle(h.RootPath()+"/", h)
	}

	mux.HandleFunc("/{$}", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		if _, err := w.Write(html); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return &http.Server{ //nolint:exhaustruct // Default values are ok.
		Addr:              opts.Address(),
		Handler:           mux,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		ReadTimeout:       time.Minute,
		WriteTimeout:      time.Minute,
	}, nil
}

// Serve starts the HTTP server for asynqmon-multi with the given context and options.
// It sets up the necessary routes for each queue and handles graceful shutdown when the context
// is cancelled.
func Serve(ctx context.Context, opts OptionsProducer) error {
	srv, err := New(opts)
	if err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(srv.ListenAndServe)

	g.Go(func() error {
		<-gCtx.Done()

		quitCtx, cancel := context.WithTimeout(context.Background(), opts.ShutdownDuration())
		defer cancel()

		//nolint:contextcheck // Parent context is already cancelled here.
		if err := srv.Shutdown(quitCtx); err != nil {
			return errors.Join(errShuttingDown, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Join(errRunServer, err)
	}

	return nil
}
