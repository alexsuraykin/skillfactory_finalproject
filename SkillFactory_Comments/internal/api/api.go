package api

import (
	"context"
	"net/http"
	"time"

	"github.com/alexsuraykin/SkillFactory_Comments/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	FileUploadBufferSize       = 512e+6 //512MB for now
	ServerShutdownDefaultDelay = 5 * time.Second
)

type Opts struct {
	Addr    string
	Log     zerolog.Logger
	Storage *storage.Storage
}

type API struct {
	l       zerolog.Logger
	server  *http.Server
	router  *gin.Engine
	storage *storage.Storage
}

func NewAPI(opts *Opts) (*API, error) {
	router := gin.Default()

	router.MaxMultipartMemory = FileUploadBufferSize

	api := &API{
		l: opts.Log,
		server: &http.Server{
			Addr:    opts.Addr,
			Handler: router,
		},
		router:  router,
		storage: opts.Storage,
	}

	router.Use(requestIDMiddleware, loggingMiddleware)

	api.setupEndpoints()

	return api, nil
}

func (api *API) setupEndpoints() {
	api.router.GET("api/comments", api.GetAllComments)
	api.router.GET("api/comments/:id", api.GetCommentById)
	api.router.POST("api/comments", api.CreateComment)
}

func (api *API) Serve() error {
	if err := api.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		api.l.Error().Err(err).Msg("failed to start api server")
		return err
	}
	return nil
}

func (hdl *API) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), ServerShutdownDefaultDelay)
	defer cancel()

	if err := hdl.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		hdl.l.Error().Err(err).Msg("failed to stop api server")
	}
}
