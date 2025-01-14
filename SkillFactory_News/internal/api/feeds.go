package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	pageSize int = 15
)

func (api *API) Feeds(ctx *gin.Context) {
	pageQuery := ctx.Query("page")
	titleQuery := ctx.Query("title")
	filterQuery := ctx.Query("filter")

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get news from storage")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid page parameter"})
		return
	}

	feeds, err := api.storage.Feeds.Feeds(page, pageSize, titleQuery, filterQuery)
	if err != nil {
		api.l.Error().Err(err).Msg("Failed to get news from storage")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get news from storage"})
		return
	}

	ctx.JSON(http.StatusOK, feeds)
}

func (api *API) FeedsById(ctx *gin.Context) {
	queryId := ctx.Param("id")

	id, err := strconv.Atoi(queryId)
	if err != nil {
		api.l.Error().Err(err).Msg("invalid arguments")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid arguments"})
		return
	}

	feed, err := api.storage.Feeds.FeedById(id)
	if err != nil {
		api.l.Error().Err(err).Msgf("failed to get news from storage:%v", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get news from storage"})
		return
	}

	ctx.JSON(http.StatusOK, feed)
}
