// Package handlers API HTTP routes and logic
package handlers

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/holygits/halcyon-engine/broker"
	"github.com/holygits/halcyon-engine/datastore"
)

// API container struct for global services
type API struct {
	Broker    broker.Broker
	Datastore datastore.Datastore
	Logger    *zap.Logger
}

// Handler defines an API HTTP handler which may return an error
type Handler func(*fasthttp.RequestCtx) error

// NewAPI constructs a new API instance
func NewAPI() (*API, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &API{
		Logger: logger,
	}, nil
}
