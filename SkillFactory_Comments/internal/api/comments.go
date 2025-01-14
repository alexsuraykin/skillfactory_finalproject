package api

import (
	"net/http"
	"strconv"

	"github.com/alexsuraykin/SkillFactory_Comments/internal/storage/queries"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/gin-gonic/gin"
)

func (api *API) CreateComment(ctx *gin.Context) {
	var req queries.Comment
	if err := ctx.BindJSON(&req); err != nil {
		api.l.Error().Err(err).Msg("failed to unmarshal comment body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse comment object body"})

		return
	}

	args := queries.CreateCommentsParams{
		NewsID:          req.NewsID,
		ParentCommentID: req.ParentCommentID,
		Content:         req.Content,
	}

	if err := api.storage.Queries.CreateComments(ctx.Request.Context(), args); err != nil {
		api.l.Error().Err(err).Msg("failed to add comment to storage")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment to storage"})
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}

func (api *API) GetAllComments(ctx *gin.Context) {
	comments, err := api.storage.Queries.GetAllComments(ctx.Request.Context())
	if err != nil {
		api.l.Error().Err(err).Msg("failed to get comments from storage")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments from storage"})

		return
	}

	res := make([]queries.Comment, len(comments))

	if len(comments) != 0 {

		for idx, val := range comments {
			comment := queries.Comment{
				ID:              val.ID,
				NewsID:          val.NewsID,
				ParentCommentID: val.ParentCommentID,
				Content:         val.Content,
			}

			res[idx] = comment
		}

		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "No comments found"})
	}
}

func (api *API) GetCommentById(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		api.l.Error().Err(err).Msg("Invalid argument")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid argument"})

		return
	}

	comments, err := api.storage.Queries.GetCommentById(ctx.Request.Context(), pgtype.Int4{Int32: int32(id), Valid: true})
	if err != nil {
		api.l.Error().Err(err).Msg("failed to get comments from storage")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments from storage"})

		return
	}

	res := make([]queries.Comment, len(comments))

	for idx, value := range comments {
		res[idx] = queries.Comment{
			ID:              value.ID,
			NewsID:          value.NewsID,
			ParentCommentID: value.ParentCommentID,
			Content:         value.Content,
			CreatedAt:       value.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, res)
}
