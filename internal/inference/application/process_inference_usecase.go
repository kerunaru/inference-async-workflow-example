package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"inference-workflow-example/internal/inference/domain"
	redis "inference-workflow-example/internal/shared/infrastructure/persistence"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type ProcessInferenceUsecase struct {
	redisClient *redis.RedisClient
	ws          *websocket.Conn
}

func NewProcessInferenceUsecase(redisClient *redis.RedisClient, ws *websocket.Conn) *ProcessInferenceUsecase {
	return &ProcessInferenceUsecase{
		redisClient: redisClient,
		ws:          ws,
	}
}

func (u *ProcessInferenceUsecase) preparePrompt(prompt string) []byte {
	return []byte(`{
    "input": {
        "prompt": "` + prompt + `",
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
}

func (u *ProcessInferenceUsecase) Execute() {
	for {
		var inferenceJob = domain.InferenceJob{}
		var client http.Client
		var req *http.Request
		var res *http.Response
		var resBody []byte
		var imgResult map[string]any

		record, err := u.redisClient.Get()
		if err != nil {
			goto takeanap
		}

		err = json.Unmarshal([]byte(record), &inferenceJob)
		if err != nil {
			log.Print("Error unmarshalling inference job: " + err.Error())
			goto takeanap
		}

		req, err = http.NewRequest(
			http.MethodPost,
			inferenceJob.Url.String(),
			bytes.NewReader(
				u.preparePrompt(inferenceJob.Prompt),
			),
		)
		if err != nil {
			log.Print("Error creating request: " + err.Error())
			goto takeanap
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("DC_AUTHORIZATION_KEY")))

		client = http.Client{
			Timeout: 30 * time.Second,
		}

		res, err = client.Do(req)
		if err != nil {
			log.Print("Error sending request: " + err.Error())
			goto takeanap
		}
		defer res.Body.Close()

		resBody, err = io.ReadAll(res.Body)
		if err != nil {
			log.Print("Error reading response body: " + err.Error())
			goto takeanap
		}

		err = json.Unmarshal(resBody, &imgResult)
		if err != nil {
			log.Print("Error unmarshalling response body: " + err.Error())
			goto takeanap
		}

		if output, ok := imgResult["output"].(map[string]any); ok {
			if outputs, ok := output["outputs"].([]any); ok && len(outputs) > 0 {
				if outputStr, ok := outputs[0].(string); ok {
					_ = u.ws.WriteMessage(websocket.TextMessage, []byte(outputStr))
				}
			}
		}
	takeanap:
		time.Sleep(20 * time.Millisecond)
	}
}
