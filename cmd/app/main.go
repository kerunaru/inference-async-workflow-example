package main

import (
	"encoding/json"
	"fmt"
	"inference-workflow-example/internal/inference/application"
	persistence "inference-workflow-example/internal/shared/infrastructure/persistence"
	server "inference-workflow-example/internal/shared/infrastructure/server"
	"io"
	"net/http"
	"os"
)

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

func saveInferenceJob(w http.ResponseWriter, r *http.Request) {
	redis := persistence.NewRedisClient(
		"redis:6379",
		"",
		0,
	)
	usecase := application.NewInferenceUseCase(redis)

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var dat map[string]any
	err = json.Unmarshal(requestBody, &dat)
	if err != nil {
		http.Error(w, "Error al descodificar el JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jobId, err := usecase.Execute(dat["prompt"].(string))
	if err != nil {
		http.Error(w, "Error al crear el trabajo de inferencia: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(fmt.Appendf(nil, `{"id": "%s"}`, jobId.String()))
}

func main() {
	httpServer := server.NewHttpServer("8080")

	httpServer.RegisterRoute("GET /", renderGUI)
	httpServer.RegisterRoute("POST /inference", saveInferenceJob)

	httpServer.Start()
}
