package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"iter"
	"maps"
	"os"
	"slices"
)

func execTemplate(htmlTemplate string, queues map[string]Queue) ([]byte, error) {
	html := bytes.Buffer{}

	t, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		return nil, errors.Join(ErrInvalidTemplate, err)
	}

	if err := t.Execute(&html, map[string]any{
		"Queues": sortedQueues(queues),
		"Title":  programName,
	}); err != nil {
		return nil, errors.Join(ErrInvalidTemplate, err)
	}

	return html.Bytes(), nil
}

func optionsFromCli(args []string) (*Options, error) {
	if args == nil {
		args = os.Args[1:]
	}

	var (
		fs     = flag.NewFlagSet(programName, flag.ExitOnError)
		queues map[string]Queue
		tpl    string

		addr = fs.String(
			"address", "",
			"Address to listen on, defaults to "+DefaultAddr,
		)
		port = fs.Uint(
			"port", 0,
			fmt.Sprintf("Port to listen on, defaults to %d", DefaultPort),
		)
		qf = fs.String(
			"queues-file", "",
			"Path to the JSON file that lists all queues. Required.",
		)
		st = fs.Int(
			"shutdown-timeout", 0,
			"Graceful shutdown timeout (seconds), defaults to "+DefaultShutdownDuration.String(),
		)
		tf = fs.String(
			"template-file", "",
			"Path to the HTML template file for the home page. Optional.",
		)
	)

	if err := fs.Parse(args); err != nil {
		return nil, errors.Join(ErrCliOptions, err)
	}

	if *qf == "" {
		return nil, errors.Join(ErrCliOptions, ErrNoQueue)
	}

	cf, err := os.ReadFile(*qf)
	if err != nil {
		return nil, errors.Join(ErrCliOptions, err)
	}

	if err := json.Unmarshal(cf, &queues); err != nil {
		return nil, errors.Join(ErrCliOptions, ErrInvalidFile, err)
	}

	if *tf != "" {
		f, err := os.ReadFile(*tf)
		if err != nil {
			return nil, errors.Join(ErrCliOptions, ErrInvalidFile, err)
		}

		tpl = string(f)
	}

	return &Options{
		Addr:            *addr,
		HTMLTemplate:    tpl,
		Port:            *port,
		Queues:          queues,
		ShutdownTimeout: *st,
	}, nil
}

func sortedQueues(queues map[string]Queue) iter.Seq2[string, Queue] {
	return func(yield func(string, Queue) bool) {
		for _, name := range slices.Sorted(maps.Keys(queues)) {
			if !yield(name, queues[name]) {
				return
			}
		}
	}
}

func validateOptions(queues map[string]Queue) error {
	if len(queues) == 0 {
		return fmt.Errorf("%w: no queues defined", ErrInvalidOption)
	}

	seen := make(map[string]string, len(queues))

	for name, q := range sortedQueues(queues) {
		k := fmt.Sprintf("%s|%d", q.RedisAddr, q.RedisDB)

		if existing, ok := seen[k]; ok {
			return fmt.Errorf("%w: duplicate queue specified by %s and %s (address=%s, db=%d)",
				ErrInvalidOption, existing, name, q.RedisAddr, q.RedisDB,
			)
		}

		seen[k] = name
	}

	return nil
}
