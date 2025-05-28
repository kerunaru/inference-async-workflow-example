package main

import (
	"inference-workflow-example/internal/inference/application"
	infrastructure "inference-workflow-example/internal/shared/infrastructure/server"
	"net/http"

	"github.com/gorilla/websocket"
)

func inference(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}

	// Don't do this on production, please use a proper origin check
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, _ := upgrader.Upgrade(w, r, nil)
	defer ws.Close()

	useCase := application.NewInferenceUseCase()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			return
		}

		if message != nil {
			result := useCase.Execute(message)

			_ = ws.WriteMessage(mt, result)
		}
	}
}

func main() {
	httpServer := infrastructure.NewHttpServer("8081")

	httpServer.RegisterRoute("/ws", inference)

	httpServer.Start()
}
