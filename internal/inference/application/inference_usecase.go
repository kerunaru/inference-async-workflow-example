package application

import (
	"inference-workflow-example/internal/inference/domain"
	persistence "inference-workflow-example/internal/shared/infrastructure/persistence"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type InferenceUseCase struct {
	redisClient *persistence.RedisClient
}

func NewInferenceUseCase(redisClient *persistence.RedisClient) *InferenceUseCase {
	return &InferenceUseCase{
		redisClient: redisClient,
	}
}

func (i *InferenceUseCase) Execute(prompt string) (uuid.UUID, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	inferenceJob, err := domain.NewInferenceJob(prompt, os.Getenv("DC_ENDPOINT"))
	if err != nil {
		log.Printf("Error creating inference job: %v", err)
		return uuid.Nil, err
	}

	i.redisClient.Set(inferenceJob.String())

	return inferenceJob.Id, nil
}
