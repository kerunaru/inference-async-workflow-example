package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

type InferenceJob struct {
	id     uuid.UUID
	prompt string
	url    *url.URL
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
		id:     uuid.New(),
		prompt: prompt,
		url:    parsedTarget,
	}, nil
}

func (i *InferenceJob) Id() uuid.UUID {
	return i.id
}

func (i *InferenceJob) Prompt() string {
	return i.prompt
}

func (i *InferenceJob) Url() *url.URL {
	return i.url
}

func (i *InferenceJob) ToJson() (string, error) {
	result, err := json.Marshal(struct {
		Id     string `json:"id"`
		Prompt string `json:"prompt"`
		Url    string `json:"url"`
	}{
		Id:     i.Id().String(),
		Prompt: i.Prompt(),
		Url:    i.Url().String(),
	})

	if err != nil {
		return "", errors.New("failed to marshal inference job")
	}

	return string(result), nil
}

func InferenceJobFromJson(jsonString string) (*InferenceJob, error) {
	var data struct {
		Id     string `json:"id"`
		Prompt string `json:"prompt"`
		Url    string `json:"url"`
	}

	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(data.Id)
	if err != nil {
		return nil, err
	}

	parsedUrl, err := url.Parse(data.Url)
	if err != nil {
		return nil, err
	}

	return &InferenceJob{
		id:     id,
		prompt: data.Prompt,
		url:    parsedUrl,
	}, nil
}
