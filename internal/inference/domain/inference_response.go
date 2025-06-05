package domain

type InferenceResponse struct {
	value string
}

func NewInferenceResponse(value string) *InferenceResponse {
	return &InferenceResponse{value: value}
}

func (ir *InferenceResponse) Value() string {
	return ir.value
}
