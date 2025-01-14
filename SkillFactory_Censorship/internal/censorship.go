package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alexsuraykin/Skillfactory_Censorship/internal/models"
	"github.com/gin-gonic/gin"
)

const (
	deepPavlovModelURL = "http://deep-pavlov:5000/model"
)

type Censor struct {
	Content []string `json:"x"`
}

func (api *API) Censorship(ctx *gin.Context) {
	var req models.Comments
	if err := ctx.BindJSON(&req); err != nil {
		api.l.Error().Err(err).Msg("failed to unmarshal comment body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse comment object body"})
		return
	}

	censorRequest := Censor{
		Content: []string{req.Content},
	}

	data, err := json.Marshal(censorRequest)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to marshal censor body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to marshal censor body"})
		return
	}

	res, err := api.aiRequest(data)
	if err != nil {
		api.l.Error().Err(err).Msg("failed to send request to AI model")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to send request to censorship service"})
		return
	}

	var censorResult [][]string
	err = json.Unmarshal(res, &censorResult)
	if err != nil {
		api.l.Debug().Err(err).Msgf("Error unmarshalling JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse censorship result"})
		return
	}

	if len(censorResult) > 0 && len(censorResult[0]) > 0 && censorResult[0][0] == "negative" {
		ctx.JSON(http.StatusBadRequest, "bad comment")
		return
	}

	ctx.JSON(http.StatusOK, "comment is okay")
}

func (api *API) aiRequest(data []byte) ([]byte, error) {
	reqURL := deepPavlovModelURL

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Minute)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(data))
	if err != nil {
		api.l.Debug().Err(err).Str("address", reqURL).Msg("Failed to create request")
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		api.l.Debug().Err(err).Str("address", reqURL).Msg("Failed to perform POST request")
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		api.l.Debug().Str("url", reqURL).Msgf("Unexpected status code: %v, response: %s", response.StatusCode, string(bodyBytes))

		return nil, fmt.Errorf("unexpected status code: %v, response: %s", response.StatusCode, string(bodyBytes))
	}

	result, err := io.ReadAll(response.Body)
	if err != nil {
		api.l.Debug().Err(err).Str("url", reqURL).Msg("Failed to read response body")
		return nil, err
	}

	return result, nil
}
