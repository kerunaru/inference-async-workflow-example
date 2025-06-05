package main

import (
	"encoding/json"
	"fmt"
	"inference-workflow-example/internal/inference/application"
	"inference-workflow-example/internal/inference/infrastructure"
	persistence "inference-workflow-example/internal/shared/infrastructure/persistence"
	server "inference-workflow-example/internal/shared/infrastructure/server"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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
	redisClient := persistence.NewRedisClient(
		fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		"",
		0,
	)

	redisRepository := infrastructure.NewRedisInferenceJobRepository(redisClient)

	usecase := application.NewInferenceUseCase(redisRepository)

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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	httpServer := server.NewHttpServer(os.Getenv("SERVER_PORT"))

	httpServer.RegisterRoute("GET /", renderGUI)
	httpServer.RegisterRoute("POST /inference", saveInferenceJob)

	httpServer.Start()
}
