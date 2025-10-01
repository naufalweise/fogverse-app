// main.go
// This is the main entry point for our application.
// It sets up the Gin web server and defines the API endpoint.

package main

import "gitlab.cs.ui.ac.id/fogverse/app/internal/backend"

func main() {
	backend.StartServer()
}
