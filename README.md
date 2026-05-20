# asynqmon-multi

Wrapper for [hibiken/asynqmon](https://github.com/hibiken/asynqmon) to manage multiple asynq servers at once.

## Usage

### Install the binary

```sh
go install github.com/luca-arch/asynqmon-multi@latest
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

	"github.com/hibiken/asynq"
	"github.com/luca-arch/asynqmon-multi/server"
)

func main() {
	// Invoke as main
	server.Main(nil)

	// Or invoke as service
	server.Serve(context.Background(), &server.Options{
		Addr:            "localhost",
		HTMLTemplate:    "", // Use default
		Port:            1337,
		ShutdownTimeout: 60,
		Queues: map[string]server.Queue{
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

The function `New` can also be utilised to instantiate a new [http.Server](https://pkg.go.dev/net/http#Server) without starting it.

## Sample queue config file

The JSON file passed via `-queues-file` should look like this. Note the field `RedisClientOpt` cannot be specified via JSON, so `server.Serve` or `server.New` must be used when, say, a username or password are required to connect to Redis.

```json
{
    "server1": {
        "redisAddr": "redis:6379",
        "redisDb": 1
    },
    "server2": {
        "readOnly": true,
        "redisAddr": "redis:6379",
        "redisDb": 2
    }
}
```

## Licence

`Asynqmon-multi` is free and open-source software licensed under the same [MIT License](./LICENSE) as [hibiken/asynqmon](https://github.com/hibiken/asynqmon).
