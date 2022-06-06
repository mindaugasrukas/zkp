package main

import "github.com/mindaugasrukas/zkp_example/server/app"

func main() {
	server := app.NewServer()
	// todo: get server port from ENV
	server.Run("8080")
}
