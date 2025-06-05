package domain

type ExternalInferenceService interface {
	PrepareRequestFromJob(job *InferenceJob) *InferenceRequest
	DoRequest(request *InferenceRequest) (*InferenceResponse, error)
}
