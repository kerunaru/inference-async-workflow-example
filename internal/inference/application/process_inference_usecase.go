package application

import (
	"inference-workflow-example/internal/inference/domain"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type ProcessInferenceUsecase struct {
	repository       domain.InferenceJobRepository
	inferenceService domain.ExternalInferenceService
	ws               *websocket.Conn
}

func NewProcessInferenceUsecase(
	repository domain.InferenceJobRepository,
	inferenceService domain.ExternalInferenceService,
	ws *websocket.Conn,
) *ProcessInferenceUsecase {
	return &ProcessInferenceUsecase{
		repository:       repository,
		inferenceService: inferenceService,
		ws:               ws,
	}
}

func (u *ProcessInferenceUsecase) Execute() {
	for {
		var inferenceResponse *domain.InferenceResponse

		job, err := u.repository.GetNextJob()
		if err != nil {
			goto takeanap
		}

		inferenceResponse, err = u.inferenceService.DoRequest(
			u.inferenceService.PrepareRequestFromJob(job),
		)
		if err != nil {
			log.Print("Error sending request: " + err.Error())
			goto takeanap
		}

		_ = u.ws.WriteMessage(websocket.TextMessage, []byte(inferenceResponse.Value()))
	takeanap:
		time.Sleep(20 * time.Millisecond)
	}
}
