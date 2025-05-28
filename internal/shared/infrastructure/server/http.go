package infrastructure

import (
	"fmt"
	"inference-workflow-example/internal/inference/application"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type HttpServer struct {
	server    *http.Server
	serverMux *http.ServeMux
}

func NewHttpServer(port string) *HttpServer {
	newServer := &HttpServer{}

	newServer.server = &http.Server{Addr: fmt.Sprintf(":%s", port)}
	newServer.serverMux = &http.ServeMux{}

	newServer.server.Handler = newServer.serverMux

	return newServer
}

func inference(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}

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

func renderGUI(w http.ResponseWriter, r *http.Request) {
	filePath := "views/index.html"

	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Error al leer el archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
}

func (s *HttpServer) Start() {
	s.serverMux.HandleFunc("/inference", inference)
	s.serverMux.HandleFunc("/", renderGUI)

	log.Printf("Listening on http://localhost%s", s.server.Addr)

	log.Fatal(s.server.ListenAndServe())
}
