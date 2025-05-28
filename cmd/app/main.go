package main

import infrastructure "inference-workflow-example/internal/shared/infrastructure/server"

func main() {
	httpServer := infrastructure.NewHttpServer("8080")

	httpServer.Start()
}
