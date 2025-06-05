package domain

type InferenceJobRepository interface {
	Save(job *InferenceJob) error
	GetNextJob() (*InferenceJob, error)
}
