package domain

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

type InferenceJob struct {
	Id     uuid.UUID `json:"id"`
	Prompt string    `json:"prompt"`
	Url    *url.URL  `json:"url"`
}

func NewInferenceJob(prompt string, target string) (*InferenceJob, error) {
	if len(prompt) == 0 {
		return nil, fmt.Errorf("prompt cannot be empty")
	}

	parsedTarget, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("url cannot be empty")
	}

	return &InferenceJob{
		Id:     uuid.New(),
		Prompt: prompt,
		Url:    parsedTarget,
	}, nil
}

func (i *InferenceJob) String() string {
	data, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf("InferenceJob{Id: %s, error marshaling: %v}", i.Id, err)
	}
	return string(data)
}
