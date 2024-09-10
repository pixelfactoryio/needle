package http_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	router "go.pixelfactory.io/needle/internal/infra/http"
	"go.pixelfactory.io/needle/internal/infra/http/handlers"
	"go.pixelfactory.io/pkg/observability/log"
)

func TestNewRouter(t *testing.T) {
	is := require.New(t)

	routes := []router.Route{
		{
			Path:    "/",
			Handler: handlers.NewDefaultHandler(),
		},
	}

	r := router.NewRouter(log.New(), routes...)
	is.NotEmpty(r)
	is.Implements((*http.Handler)(nil), r)
}
