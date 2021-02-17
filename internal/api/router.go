package api

import (
	"net/http"

	"go.pixelfactory.io/pkg/observability/log"

	"github.com/gorilla/mux"
	"go.pixelfactory.io/needle/internal/api/middleware"
)

type Route struct {
	Path        string
	HandlerFunc http.HandlerFunc
}

// NewRouter create router, setup routes and middlewares
func NewRouter(logger log.Logger, routes ...Route) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.Logging(logger))

	for _, r := range routes {
		router.PathPrefix(r.Path).HandlerFunc(r.HandlerFunc)
	}

	return router
}
