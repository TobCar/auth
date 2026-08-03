package scim

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
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
	resourceTypes         []*core.ResourceType
	schemas               []*core.Schema
}

func NewServer(config *conf.GlobalConfiguration, db *storage.Connection) *Server {
	baseURL := strings.TrimRight(config.API.ExternalURL, "/") + BasePath
	userSchema := newUserSchema(baseURL)

	return &Server{
		db:    db,
		users: NewUserMapper(baseURL),
		serviceProviderConfig: core.NewServiceProviderConfig(
			baseURL,
			[]core.AuthenticationScheme{core.OAuthBearerToken().AsPrimary()},
		),
		resourceTypes: []*core.ResourceType{core.NewResourceType(baseURL, userSchema, core.EndpointUsers)},
		schemas:       []*core.Schema{userSchema},
	}
}

func (srv *Server) ServiceProviderConfig(w http.ResponseWriter, r *http.Request) error {
	return protocol.Send(w, http.StatusOK, srv.serviceProviderConfig)
}

func (srv *Server) ResourceTypes(w http.ResponseWriter, r *http.Request) error {
	return list(w, r, srv.resourceTypes)
}

func (srv *Server) ResourceTypeByID(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	for _, resourceType := range srv.resourceTypes {
		if resourceType.ID == core.ResourceTypeName(id) {
			return protocol.Send(w, http.StatusOK, resourceType)
		}
	}
	return srv.NotFound(w, r)
}

func (srv *Server) Schemas(w http.ResponseWriter, r *http.Request) error {
	return list(w, r, srv.schemas)
}

func (srv *Server) SchemaByID(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	for _, schema := range srv.schemas {
		if schema.ID == core.SchemaURI(id) {
			return protocol.Send(w, http.StatusOK, schema)
		}
	}
	return srv.NotFound(w, r)
}

func (srv *Server) NotFound(w http.ResponseWriter, r *http.Request) error {
	return protocol.NewError(http.StatusNotFound, "", "Endpoint or resource does not exist")
}

func (srv *Server) MethodNotAllowed(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Allow", http.MethodGet)
	return protocol.NewError(http.StatusMethodNotAllowed, "", "The request method is not supported by this endpoint")
}

// RFC 7644, Section 4: filtering SHALL be ignored on these endpoints, and a
// supplied filter SHOULD be answered with a 403 so clients cannot assume the
// conditions were applied.
func list[T any](w http.ResponseWriter, r *http.Request, resources []T) error {
	if r.URL.Query().Get("filter") != "" {
		return protocol.NewError(http.StatusForbidden, "", "Filtering is not supported on this endpoint")
	}
	return protocol.Send(w, http.StatusOK, protocol.NewListResponse(resources))
}
