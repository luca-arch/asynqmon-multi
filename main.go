package main

import (
	"github.com/luca-arch/asynqmon-multi/server"
)

func main() {
	if err := server.Main(nil); err != nil {
		panic(err)
	}
}
