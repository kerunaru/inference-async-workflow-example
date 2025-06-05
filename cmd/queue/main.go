package main

import (
	"fmt"
	"inference-workflow-example/internal/inference/application"
	"inference-workflow-example/internal/inference/infrastructure"
	redis "inference-workflow-example/internal/shared/infrastructure/persistence"
	server "inference-workflow-example/internal/shared/infrastructure/server"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func inference(w http.ResponseWriter, r *http.Request) {
	redisClient := redis.NewRedisClient(
		fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		"",
		0,
	)

	redisInferenceJobRepository := infrastructure.NewRedisInferenceJobRepository(redisClient)

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	dataCrunchInferenceRequestService := infrastructure.NewDataCrunchInferenceRequestService(httpClient)

	upgrader := websocket.Upgrader{}

	// Don't do this on production, please use a proper origin check
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, _ := upgrader.Upgrade(w, r, nil)
	defer ws.Close()

	usecase := application.NewProcessInferenceUsecase(
		redisInferenceJobRepository,
		dataCrunchInferenceRequestService,
		ws,
	)

	usecase.Execute()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	httpServer := server.NewHttpServer(os.Getenv("QUEUE_PORT"))

	httpServer.RegisterRoute("/ws", inference)

	httpServer.Start()
}
