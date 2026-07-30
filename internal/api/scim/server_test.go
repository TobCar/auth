package scim

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/supabase/auth/internal/api/scim/protocol"
)

func TestServer(t *testing.T) {
	srv := NewServer(nil)
	require.NotNil(t, srv)

	for _, tc := range []struct {
		path    string
		handler func(http.ResponseWriter, *http.Request) error
	}{
		{"ServiceProviderConfig", srv.ServiceProviderConfig},
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
