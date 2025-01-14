package api

import (
	"context"
	"net/http"
	"skillfactory/SkillFactory_finalProject/APIGateway/internal/api/oapi"
	"sync"
	"time"

	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	FileUploadBufferSize       = 512e+6 //512MB for now
	ServerShutdownDefaultDelay = 5 * time.Second
)

type Opts struct {
	Addr        string
	Log         zerolog.Logger
	WG          *sync.WaitGroup
	RequestChan chan []byte
	SwaggerFile []byte
}

type API struct {
	l           zerolog.Logger
	server      *http.Server
	router      *gin.Engine
	swaggerFile []byte
	wg          *sync.WaitGroup
	requestChan chan []byte
}

func NewAPI(opts *Opts) (*API, error) {
	router := gin.Default()

	swagger, err := oapi.GetSwagger()
	if err != nil {
		return nil, err
	}

	oapiOpts := &middleware.Options{
		Options: openapi3filter.Options{
			ExcludeRequestBody: true,
		},
	}

	router.MaxMultipartMemory = FileUploadBufferSize

	wg := &sync.WaitGroup{}

	api := &API{
		l: opts.Log,
		server: &http.Server{
			Addr:    opts.Addr,
			Handler: router,
		},
		router:      router,
		swaggerFile: opts.SwaggerFile,
		wg:          wg,
		requestChan: opts.RequestChan,
	}

	router.Use(requestIDMiddleware, loggingMiddleware)

	router.Use(middleware.OapiRequestValidatorWithOptions(swagger, oapiOpts))

	oapi.RegisterHandlersWithOptions(router, api, oapi.GinServerOptions{
		BaseURL: "/api",
	})

	return api, nil
}

func (api *API) Serve() error {
	if err := api.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		api.l.Error().Err(err).Msg("failed to start api server")
		return err
	}
	return nil
}

func (api *API) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), ServerShutdownDefaultDelay)
	defer cancel()

	if err := api.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		api.l.Error().Err(err).Msg("failed to stop api server")
	}
}
