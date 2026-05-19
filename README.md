# asynqmon-multi

Wrapper for [hibiken/asynqmon](https://github.com/hibiken/asynqmon) to manage multiple asynq servers at once.

## Usage

### Install the binary

```sh
go install github.com/luca-arch/asynqmon-multi
```

Display usage

```sh
asynqmon-multi -h
```

```text
Usage of asynqmon-multi:
  -address string
        Address to listen on, defaults to 0.0.0.0
  -port uint
        Port to listen on, defaults to 8080
  -queues-file string
        Path to the JSON file that lists all queues. Required.
  -shutdown-timeout int
        Graceful shutdown timeout (seconds), defaults to 10s
  -template-file string
        Path to the HTML template file for the home page. Optional.
```

### Import as a Library

Install the package

```sh
go get github.com/luca-arch/asynqmon-multi
```

Invoke `Main` or `Serve`:

```go
package main

import (
	"context"

	asynqmonmulti "github.com/luca-arch/asynqmon-multi"
	"github.com/hibiken/asynq"
)

func main() {
	// Invoke as main
	asynqmonmulti.Main(nil)

	// Or invoke as service
	asynqmonmulti.Serve(context.Background(), &asynqmonmulti.Options{
		Addr:            "localhost",
		HTMLTemplate:    "", // Use default
		Port:            1337,
		ShutdownTimeout: 60,
		Queues: map[string]asynqmonmulti.Queue{
			"server1": {
				RedisAddr: "redis:6379",
				RedisDB:   1,
			},
			"server2": {
				RedisAddr: "redis:6379",
				RedisDB:   2,
			},
			"server3": {
				ReadOnly:  true,
				RedisAddr: "redis:6379",
				RedisDB:   3,
			},
			"server4": {
				RedisAddr: "redis:6379",
				RedisDB:   3,
				RedisClientOpt: &asynq.RedisClientOpt{
					Network:  "redis_network",
					Password: "redis_password",
					Username: "redis_user",
				},
			},
		},
	})
}
```

## Sample queue config file

```json
{
    "server1": {
        "RedisAddr": "redis:6379",
        "RedisDB": 1
    },
    "server2": {
        "ReadOnly": true,
        "RedisAddr": "redis:6379",
        "RedisDB": 2
    }
}
```

## Licence

`Asynqmon-multi` is free and open-source software licensed under the same [MIT License](./LICENSE) as [hibiken/asynqmon](https://github.com/hibiken/asynqmon).
