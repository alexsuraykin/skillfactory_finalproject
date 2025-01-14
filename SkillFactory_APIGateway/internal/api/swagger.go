package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) GetSwagger(ctx *gin.Context) {
	ctx.String(http.StatusOK, string(api.swaggerFile))
}
