package scim

import (
	"net/http"
	"strings"

	"github.com/supabase/auth/internal/api/scim/core"
	"github.com/supabase/auth/internal/api/scim/protocol"
	"github.com/supabase/auth/internal/conf"
	"github.com/supabase/auth/internal/models"
	"github.com/supabase/auth/internal/storage"
)

const BasePath = "/scim/v2"

type Server struct {
	db                    *storage.Connection
	users                 Mapper[*models.User, *core.User]
	serviceProviderConfig *core.ServiceProviderConfig
}

func NewServer(config *conf.GlobalConfiguration, db *storage.Connection) *Server {
	baseURL := strings.TrimRight(config.API.ExternalURL, "/") + BasePath

	return &Server{
		db:    db,
		users: NewUserMapper(baseURL),
		serviceProviderConfig: core.NewServiceProviderConfig(
			baseURL,
			[]core.AuthenticationScheme{core.OAuthBearerToken().AsPrimary()},
		),
	}
}

func (srv *Server) ServiceProviderConfig(w http.ResponseWriter, r *http.Request) error {
	return protocol.Send(w, http.StatusOK, srv.serviceProviderConfig)
}

func (srv *Server) ResourceTypes(w http.ResponseWriter, r *http.Request) error {
	return emptyList(w, r, []any{})
}

func (srv *Server) Schemas(w http.ResponseWriter, r *http.Request) error {
	return emptyList(w, r, []any{})
}

func (srv *Server) NotFound(w http.ResponseWriter, r *http.Request) error {
	return protocol.NewError(http.StatusNotFound, "", "Endpoint or resource does not exist")
}

func (srv *Server) MethodNotAllowed(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Allow", http.MethodGet)
	return protocol.NewError(http.StatusMethodNotAllowed, "", "The request method is not supported by this endpoint")
}

func emptyList[T any](w http.ResponseWriter, r *http.Request, resources []T) error {
	if r.URL.Query().Get("filter") != "" {
		return protocol.NewError(http.StatusForbidden, "", "Filtering is not supported on this endpoint")
	}
	return protocol.Send(w, http.StatusOK, protocol.NewListResponse(resources))
}
