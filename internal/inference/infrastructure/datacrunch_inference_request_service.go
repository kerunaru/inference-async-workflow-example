package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"inference-workflow-example/internal/inference/domain"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
)

type DataCrunchInferenceRequestService struct {
	client *http.Client
}

func NewDataCrunchInferenceRequestService(client *http.Client) *DataCrunchInferenceRequestService {
	return &DataCrunchInferenceRequestService{
		client: client,
	}
}

func (d *DataCrunchInferenceRequestService) PrepareRequestFromJob(job *domain.InferenceJob) *domain.InferenceRequest {
	payload := []byte(`{
    "input": {
        "prompt": "` + job.Prompt() + `",
        "size": "1024*1024",
        "num_inference_steps": 35,
        "seed": ` + strconv.Itoa(rand.IntN(65535)) + `,
        "guidance_scale": 5,
        "num_images": 1,
        "enable_safety_checker": false,
        "enable_async_output": false,
        "enable_base64_output": false,
        "cache_threshold": 0.1
    }
}`)
	return domain.NewInferenceRequest(
		http.MethodPost,
		string(payload),
		*job.Url(),
	)
}

func (d *DataCrunchInferenceRequestService) DoRequest(request *domain.InferenceRequest) (*domain.InferenceResponse, error) {
	req, err := http.NewRequest(
		request.Method(),
		request.Url().String(),
		bytes.NewReader(
			[]byte(request.Value()),
		),
	)
	if err != nil {
		log.Print("Error creating request: " + err.Error())
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("DC_AUTHORIZATION_KEY")))

	res, err := d.client.Do(req)
	if err != nil {
		log.Print("Error sending request: " + err.Error())
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Print("Error reading response body: " + err.Error())
		return nil, err
	}

	var imgResult map[string]any

	err = json.Unmarshal(resBody, &imgResult)
	if err != nil {
		log.Print("Error unmarshalling response body: " + err.Error())
		return nil, err
	}

	if output, ok := imgResult["output"].(map[string]any); ok {
		if outputs, ok := output["outputs"].([]any); ok && len(outputs) > 0 {
			if outputStr, ok := outputs[0].(string); ok {
				return domain.NewInferenceResponse(outputStr), nil
			}
		}
	}

	return nil, errors.New("No output in response from DataCrunch")
}
