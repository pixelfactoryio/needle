package api_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.pixelfactory.io/needle/internal/api"
	"go.pixelfactory.io/needle/internal/api/handlers"
	"go.pixelfactory.io/pkg/observability/log"
)

func TestNewRouter(t *testing.T) {
	t.Parallel()
	is := require.New(t)

	routes := []api.Route{
		{
			Path:    "/",
			Handler: handlers.NewDefaultHandler(),
		},
	}

	r := api.NewRouter(log.New(), routes...)
	is.NotEmpty(r)
	is.Implements((*http.Handler)(nil), r)
}
