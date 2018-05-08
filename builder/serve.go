package builder

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/holygits/halcyon-engine/broker"
	"github.com/holygits/halcyon-engine/handlers"
)

func main() {

	// TODO: setup services {datastore, broker, build container}
	api, _ := handlers.NewAPI()
	go api.Broker.Poll()
	// TODO: start container builder channel / service

	router := fasthttprouter.New()
	router.GET("/images/:id", mw())
	router.DELETE("/images/:id", mw())

	log.Fatal("Failed to start server", zap.Error(fasthttp.ListenAndServe(":7098", router.Handler)))
}

// mw provides middleware for internal handlers
func mw(h handler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if err := h(ctx); err != nil {
			// TODO: handle error codes case by case
			log.Error("Request error", zap.Error(err))
		}
	}
}

func getImageHandler(ctx *fasthttp.RequestCtx) error {
	id, _ := ctx.UserValue("id").(string)
	// TODO: load image from db
}

func deleteImageHandler(ctx *fasthttp.RequestCtx) error {
	id, _ := ctx.UserValue("id").(string)
	// TODO: remove image from db
	return nil
}
