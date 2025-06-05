package infrastructure

import (
	"inference-workflow-example/internal/inference/domain"
	redis "inference-workflow-example/internal/shared/infrastructure/persistence"
)

type RedisInferenceJobRepository struct {
	client *redis.RedisClient
}

func NewRedisInferenceJobRepository(client *redis.RedisClient) *RedisInferenceJobRepository {
	return &RedisInferenceJobRepository{
		client: client,
	}
}

func (r *RedisInferenceJobRepository) Save(job *domain.InferenceJob) error {
	record, err := job.ToJson()
	if err != nil {
		return err
	}

	return r.client.Set(record)
}

func (r *RedisInferenceJobRepository) GetNextJob() (*domain.InferenceJob, error) {
	record, err := r.client.Get()
	if err != nil {
		return nil, err
	}

	inferenceJob, err := domain.InferenceJobFromJson(record)
	if err != nil {
		return nil, err
	}

	return inferenceJob, nil
}
