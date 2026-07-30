package scim

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/supabase/auth/internal/api/scim/fixtures"
	"github.com/supabase/auth/internal/api/scim/protocol"
	"github.com/supabase/auth/internal/conf"
)

func newServerFor(externalURL string) *Server {
	return NewServer(&conf.GlobalConfiguration{
		API: conf.APIConfiguration{ExternalURL: externalURL},
	})
}

func TestServer(t *testing.T) {
	srv := newServerFor("http://localhost:9999")
	require.NotNil(t, srv)

	t.Run("NewServer trims a trailing slash from the external URL", func(t *testing.T) {
		location := newServerFor("https://auth.example.com/").serviceProviderConfig.Meta.Location

		require.Equal(t, "https://auth.example.com"+BasePath+"/ServiceProviderConfig", location)
	})

	t.Run("ServiceProviderConfig", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, BasePath+"/ServiceProviderConfig", nil)
		w := httptest.NewRecorder()

		require.NoError(t, srv.ServiceProviderConfig(w, r))

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, protocol.MediaType, w.Header().Get("Content-Type"))
		require.JSONEq(t, fixtures.ServiceProviderConfig, w.Body.String())
	})

	for _, tc := range []struct {
		path    string
		handler func(http.ResponseWriter, *http.Request) error
	}{
		{"ResourceTypes", srv.ResourceTypes},
		{"Schemas", srv.Schemas},
	} {
		t.Run(tc.path, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, BasePath+"/"+tc.path, nil)
			w := httptest.NewRecorder()

			var scimErr *protocol.Error
			require.ErrorAs(t, tc.handler(w, r), &scimErr)

			require.Equal(t, http.StatusNotImplemented, scimErr.StatusCode())
		})
	}

	t.Run("NotFound", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, BasePath+"/Unknown", nil)
		w := httptest.NewRecorder()

		var scimErr *protocol.Error
		require.ErrorAs(t, srv.NotFound(w, r), &scimErr)

		require.Equal(t, http.StatusNotFound, scimErr.StatusCode())
		require.Equal(t, "Endpoint or resource does not exist", scimErr.Detail)
	})

	t.Run("NotAllowed", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, BasePath+"/Schemas", nil)
		w := httptest.NewRecorder()

		var scimErr *protocol.Error
		require.ErrorAs(t, srv.MethodNotAllowed(w, r), &scimErr)

		require.Equal(t, http.StatusMethodNotAllowed, scimErr.StatusCode())
		require.Equal(t, http.MethodGet, w.Header().Get("Allow"))
	})
}
