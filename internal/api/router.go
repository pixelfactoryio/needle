// Package api provides HTTP api and Route.
package api

import (
	"net/http"

	"go.pixelfactory.io/pkg/observability/log"

	"github.com/gorilla/mux"
	"go.pixelfactory.io/needle/internal/api/middleware"
)

// Route holds path and http.Handler
type Route struct {
	Path    string
	Handler http.Handler
}

// NewRouter create router, setup routes and middlewares
func NewRouter(logger log.Logger, routes ...Route) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.Logging(logger))

	for _, r := range routes {
		router.PathPrefix(r.Path).Handler(r.Handler)
	}

	return router
}
