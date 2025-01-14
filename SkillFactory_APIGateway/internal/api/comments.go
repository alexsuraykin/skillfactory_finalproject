package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"skillfactory/SkillFactory_finalProject/APIGateway/internal/api/oapi"
	"skillfactory/SkillFactory_finalProject/APIGateway/internal/models"

	"github.com/gin-gonic/gin"
)

const (
	commentsServiceURL   = "http://comments:8882/api/comments"
	censorshipServiceURL = "http://censorship:8884/api/censorship"
)

func (api *API) CreateComment(ctx *gin.Context) {
	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to read request body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}

	check, err := http.Post(censorshipServiceURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to censor comment")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to censor comment"})
		return
	}
	defer check.Body.Close()

	if check.StatusCode != http.StatusOK {
		api.l.Error().Msgf("Unexpected status code: %d", check.StatusCode)
		ctx.JSON(check.StatusCode, gin.H{"error": "The comment contains obscene language."})
		return
	}

	resp, err := http.Post(commentsServiceURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to create comment in comments service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		api.l.Error().Msgf("Unexpected status code: %d", resp.StatusCode)
		ctx.JSON(resp.StatusCode, gin.H{"error": "Failed to get comments from comments service"})
		return
	}

	res, err := oapi.ParseCreateCommentResponse(resp)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to read response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	ctx.Data(resp.StatusCode, "application/json", res.Body)
}

func (api *API) GetAllComments(ctx *gin.Context) {
	resp, err := http.Get(commentsServiceURL)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get comments from comments service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments from comments service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		api.l.Error().Msgf("Unexpected status code: %d", resp.StatusCode)
		ctx.JSON(resp.StatusCode, gin.H{"error": "Failed to get comments from comments service"})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to read response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	res := oapi.GetAllCommentsResponse{}

	var comments []models.Comments

	err = json.Unmarshal(body, &comments)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to unmarshal response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response body"})
		return
	}

	list := make([]oapi.Comment, len(comments))

	for idx, val := range comments {
		comment := oapi.Comment{
			Id:              int(val.ID),
			NewsID:          *val.NewsID,
			ParentCommentID: val.ParentCommentID,
			Content:         val.Content,
			CreatedAt:       *val.CreatedAt,
		}

		list[idx] = comment
	}

	res.JSON200 = &list

	if res.JSON200 != nil {
		ctx.JSON(http.StatusOK, res.JSON200)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "No comments found"})
	}
}
