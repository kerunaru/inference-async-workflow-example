package domain

import "net/url"

type InferenceRequest struct {
	method string
	value  string
	url    url.URL
}

func NewInferenceRequest(
	method string,
	value string,
	url url.URL,
) *InferenceRequest {
	return &InferenceRequest{
		method: method,
		value:  value,
		url:    url,
	}
}

func (ir *InferenceRequest) Method() string {
	return ir.method
}

func (ir *InferenceRequest) Value() string {
	return ir.value
}

func (ir *InferenceRequest) Url() *url.URL {
	return &ir.url
}
