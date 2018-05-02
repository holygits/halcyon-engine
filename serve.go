package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/holygits/halcyon-engine/handlers"
)

var (
	log *zap.Logger
)

func setup() {
	log, _ = zap.NewProduction()
	defer log.Sync()
}

func main() {

	// Setup container for global API services
	api, err := handlers.NewAPI()
	if err != nil {
		log.Fatal("Failed to setup server", zap.Error(err))
	}

	// API routes
	router := fasthttprouter.New()
	// Functions
	router.GET("/func/:id", mw(api.FuncInfo))
	router.POST("/func/:id", mw(api.FuncExec))
	router.DELETE("/func/:id", mw(api.FuncDel))
	router.POST("/func", mw(api.FuncCreate))

	// Executions
	router.GET("/func/:id/executions/:id", mw(api.FuncExecInfo))
	router.GET("/pipe/:id/executions/:id", mw(api.PipeExecInfo))

	// Pipelines
	router.GET("/pipe/:id", mw(api.PipeInfo))
	router.POST("/pipe/:id", mw(api.PipeExec))
	router.DELETE("/pipe/:id", mw(api.PipeDel))
	router.POST("/pipe", mw(api.PipeCreate))

	log.Fatal("Failed to start server", zap.Error(fasthttp.ListenAndServe(":7097", router.Handler)))
}

// mw provides middleware for internal handlers
func mw(h handlers.Handler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if err := h(ctx); err != nil {
			// TODO: handle error codes case by case
			log.Error("Request error", zap.Error(err))
		}
	}
}
