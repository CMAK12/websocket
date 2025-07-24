package main

import "cw/internal/server"

func main() {
	if err := server.MustRun(); err != nil {
		panic(err)
	}
}
