package main

import (
	asynqmonmulti "github.com/luca-arch/asynqmon-multi"
)

func main() {
	if err := asynqmonmulti.Main(nil); err != nil {
		panic(err)
	}
}
