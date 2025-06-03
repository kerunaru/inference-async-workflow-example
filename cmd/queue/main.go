package main

import (
	"inference-workflow-example/internal/inference/application"
	redis "inference-workflow-example/internal/shared/infrastructure/persistence"
	server "inference-workflow-example/internal/shared/infrastructure/server"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func inference(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisClient := redis.NewRedisClient(
		"redis:6379",
		"",
		0,
	)

	upgrader := websocket.Upgrader{}

	// Don't do this on production, please use a proper origin check
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, _ := upgrader.Upgrade(w, r, nil)
	defer ws.Close()

	usecase := application.NewProcessInferenceUsecase(redisClient, ws)

	usecase.Execute()
}

func main() {
	httpServer := server.NewHttpServer("8081")

	httpServer.RegisterRoute("/ws", inference)

	httpServer.Start()
}
