package scim

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/supabase/auth/internal/api/scim/protocol"
	"github.com/supabase/auth/internal/models"
)

func (srv *Server) UserByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		return userNotFound()
	}

	user, err := models.FindUserByIDAndSSOProviderID(srv.db.WithContext(ctx), id, providerKey.From(ctx).ID)
	if err != nil {
		if models.IsNotFoundError(err) {
			return userNotFound()
		}
		return protocol.Wrap(err)
	}

	return protocol.Send(w, http.StatusOK, srv.users.MapFrom(user))
}

func userNotFound() error {
	return protocol.NewError(http.StatusNotFound, "", "Resource not found")
}
