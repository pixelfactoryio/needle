package api

import (
	"go.pixelfactory.io/pkg/observability/log"

	"github.com/gorilla/mux"
	"go.pixelfactory.io/needle/internal/api/handlers"
	"go.pixelfactory.io/needle/internal/api/middleware"
)

// NewRouter create router, setup routes and middlewares
func NewRouter(logger log.Logger, caFile string) *mux.Router {
	handler := handlers.NewPixelHandler(caFile)
	router := mux.NewRouter()
	router.Use(middleware.Logging(logger))
	router.PathPrefix("/install-root-ca").HandlerFunc(handler.GetRootCA)
	router.PathPrefix("/").HandlerFunc(handler.GetPixel)
	return router
}
