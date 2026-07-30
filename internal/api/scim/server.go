package scim

import (
	"net/http"

	"github.com/supabase/auth/internal/api/scim/protocol"
	"github.com/supabase/auth/internal/conf"
)

const BasePath = "/scim/v2"

type Server struct {
	config *conf.GlobalConfiguration
}

func NewServer(config *conf.GlobalConfiguration) *Server {
	return &Server{
		config: config,
	}
}

func (srv *Server) ServiceProviderConfig(w http.ResponseWriter, r *http.Request) error {
	return notImplemented()
}

func (srv *Server) ResourceTypes(w http.ResponseWriter, r *http.Request) error {
	return notImplemented()
}

func (srv *Server) Schemas(w http.ResponseWriter, r *http.Request) error {
	return notImplemented()
}

func (srv *Server) NotFound(w http.ResponseWriter, r *http.Request) error {
	return protocol.NewError(http.StatusNotFound, "", "Endpoint or resource does not exist")
}

func (srv *Server) MethodNotAllowed(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Allow", http.MethodGet)
	return protocol.NewError(http.StatusMethodNotAllowed, "", "The request method is not supported by this endpoint")
}

func notImplemented() error {
	return protocol.NewError(http.StatusNotImplemented, "", "The request endpoint is not implemented")
}
