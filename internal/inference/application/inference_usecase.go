package application

import (
	"log"
	"math/rand"
	"time"
)

type InferenceUseCase struct {
}

func NewInferenceUseCase() *InferenceUseCase {
	return &InferenceUseCase{}
}

func (i *InferenceUseCase) Execute(message []byte) []byte {
	log.Printf("Message received: %s", message)

	// Call to inference service synchronously
	time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)

	result := make([]byte, len("Inference response"))
	copy(result, "Inference response")

	log.Printf("Message sent: %s", result)

	return result
}
