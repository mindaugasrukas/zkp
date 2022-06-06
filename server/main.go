package main

func main() {
	server := NewServer()
	// todo: get server port from ENV
	server.Run("8080")
}
