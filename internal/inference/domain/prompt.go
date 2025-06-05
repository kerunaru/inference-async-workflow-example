package domain

import (
	"errors"
	"strings"
)

type Prompt struct {
	value string
}

func NewPrompt(value string) (*Prompt, error) {
	if len(strings.TrimSpace(value)) == 0 {
		return nil, errors.New("prompt value cannot be empty")
	}

	return &Prompt{value: value}, nil
}
