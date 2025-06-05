package application

import (
	"inference-workflow-example/internal/inference/domain"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type InferenceUseCase struct {
	inferenceJobRepository domain.InferenceJobRepository
}

func NewInferenceUseCase(inferenceJobRepository domain.InferenceJobRepository) *InferenceUseCase {
	return &InferenceUseCase{
		inferenceJobRepository: inferenceJobRepository,
	}
}

func (i *InferenceUseCase) Execute(prompt string) (uuid.UUID, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}

	inferenceJob, err := domain.NewInferenceJob(prompt, os.Getenv("DC_ENDPOINT"))
	if err != nil {
		log.Printf("Error creating inference job: %v", err)
		return uuid.Nil, err
	}

	i.inferenceJobRepository.Save(inferenceJob)

	return inferenceJob.Id(), nil
}
