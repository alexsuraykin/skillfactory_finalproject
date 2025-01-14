package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"skillfactory/SkillFactory_finalProject/APIGateway/internal/api/oapi"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	feedsServiceURL = "http://news:8883/api/feeds"
)

func (api *API) Feeds(ctx *gin.Context, params oapi.FeedsParams) {
	parsedURL, err := url.Parse(feedsServiceURL)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to parse url address feeds service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse url address feeds service"})
		return
	}

	if *params.Page <= 0 {
		*params.Page = 1
	}

	paramsQuery := url.Values{}
	paramsQuery.Add("page", strconv.FormatInt(int64(*params.Page), 10))
	paramsQuery.Add("title", *params.Title)
	paramsQuery.Add("filter", *params.Filter)
	*params.RequestId = ctx.Request.Context().Value(requestIDKey).(string)
	paramsQuery.Add("request_id", *params.RequestId)

	parsedURL.RawQuery = paramsQuery.Encode()

	reqURL := parsedURL.String()

	resp, err := http.Get(reqURL)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get feeds from feeds service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feeds from feeds service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		api.l.Error().Msgf("Unexpected status code: %d", resp.StatusCode)
		ctx.JSON(resp.StatusCode, gin.H{"error": "Failed to get feeds from feeds service"})
		return
	}

	res, err := oapi.ParseFeedsResponse(resp)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to read response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	if res.JSON200 != nil {
		ctx.JSON(http.StatusOK, res.JSON200)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "No feeds found"})
	}
}

func (api *API) FeedsById(ctx *gin.Context, id oapi.ID, params oapi.FeedsByIdParams) {
	if id <= 0 {
		api.l.Error().Msg("Invalid ID argument")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID argument"})
		return
	}

	idStr := strconv.FormatInt(int64(id), 10)

	var feeds oapi.Feeds
	var comments []oapi.Comment
	var res oapi.FeedsByIdResponse

	requestChan := make(chan any, 2)
	errorChan := make(chan error, 2)
	defer close(requestChan)
	defer close(errorChan)

	api.wg.Add(2)
	go api.prepareRequestFeedsById(ctx, idStr, params, requestChan, errorChan)
	go api.prepareRequestCommentById(ctx, idStr, params, requestChan, errorChan)
	api.wg.Wait()

	for i := 0; i < 2; i++ {
		select {
		case res := <-requestChan:
			switch v := res.(type) {
			case oapi.Feeds:
				feeds = v
			case []oapi.Comment:
				comments = v
			}
		case err := <-errorChan:
			api.l.Error().Err(err).Msg("Failed to process request")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	feedsById := oapi.FeedsById{
		Feeds:    feeds,
		Comments: comments,
	}

	list := make([]oapi.FeedsById, 1)
	list[0] = feedsById

	res.JSON200 = &list

	if res.JSON200 != nil {
		ctx.JSON(http.StatusOK, res.JSON200)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "No feeds found"})
	}
}

func (api *API) prepareRequestFeedsById(ctx *gin.Context, id string, params oapi.FeedsByIdParams, requestChan chan<- any, errorChan chan<- error) {
	defer api.wg.Done()

	parsedURL, err := url.Parse(feedsServiceURL + "/" + id)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to parse url address feeds service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse url address feeds service"})
		errorChan <- err
		return
	}

	paramsQuery := url.Values{}
	*params.RequestId = ctx.Request.Context().Value(requestIDKey).(string)
	paramsQuery.Add("request_id", *params.RequestId)

	parsedURL.RawQuery = paramsQuery.Encode()

	reqURL := parsedURL.String()

	resp, err := http.Get(reqURL)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get feeds by id from feeds service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feeds by id from feeds service"})
		errorChan <- err
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		api.l.Error().Msgf("Unexpected status code: %d", resp.StatusCode)
		ctx.JSON(resp.StatusCode, gin.H{"error": "Failed to get feeds by id from feeds service"})
		errorChan <- err
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to read response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		errorChan <- err
		return
	}

	var feeds oapi.Feeds
	if err := json.Unmarshal(body, &feeds); err != nil {
		api.l.Error().Err(err).Msg("Failed to unmarshal feeds response")
		errorChan <- err
		return
	}

	requestChan <- feeds

}

func (api *API) prepareRequestCommentById(ctx *gin.Context, id string, params oapi.FeedsByIdParams, requestChan chan<- any, errorChan chan<- error) {
	defer api.wg.Done()

	parsedURL, err := url.Parse(commentsServiceURL + "/" + id)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to parse url address comments service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse url address comments service"})
		errorChan <- err
		return
	}

	paramsQuery := url.Values{}
	*params.RequestId = ctx.Request.Context().Value(requestIDKey).(string)
	paramsQuery.Add("request_id", *params.RequestId)

	parsedURL.RawQuery = paramsQuery.Encode()

	reqURL := parsedURL.String()

	resp, err := http.Get(reqURL)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get comment by id from comments service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comment by id from comments service"})
		errorChan <- err
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		api.l.Error().Msgf("Unexpected status code: %d", resp.StatusCode)
		ctx.JSON(resp.StatusCode, gin.H{"error": "Failed to get comment by id from comments service"})
		errorChan <- err
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to read response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		errorChan <- err
		return
	}

	var comments []oapi.Comment

	if err := json.Unmarshal(body, &comments); err != nil {
		api.l.Error().Err(err).Msg("Failed to unmarshal comments response")
		errorChan <- err
		return
	}

	requestChan <- comments
}

/*
func (api *API) FeedsById(ctx *gin.Context, id oapi.ID) {
	if id <= 0 {
		api.l.Error().Msg("Invalid ID argument")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID argument"})
		return
	}

	idStr := strconv.FormatInt(int64(id), 10)

	resp, err := http.Get(feedsServiceURL + "/" + idStr)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get feeds by id from feeds service")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feeds by id from feeds service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		api.l.Error().Msgf("Unexpected status code: %d", resp.StatusCode)
		ctx.JSON(resp.StatusCode, gin.H{"error": "Failed to get feeds by id from feeds service"})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to read response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var feed oapi.Feeds
	var res oapi.FeedsByIdResponse

	err = json.Unmarshal(body, &feed)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to unmarshal response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response body"})
		return
	}

	list := make([]oapi.Feeds, 1)
	list[0] = feed

	res.JSON200 = &list

	if res.JSON200 != nil {
		ctx.JSON(http.StatusOK, res.JSON200)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "No feeds found"})
	}
}
*/
