package main

import (
	infrastructure "inference-workflow-example/internal/shared/infrastructure/server"
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
	// Implement the logic to save the inference job
	w.Write([]byte("Inference job saved"))
}

func main() {
	httpServer := infrastructure.NewHttpServer("8080")

	httpServer.RegisterRoute("GET /", renderGUI)
	httpServer.RegisterRoute("POST /inference", saveInferenceJob)

	httpServer.Start()
}
