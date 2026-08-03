package scim

import (
	"context"
	"net/http"
	"strings"

	"github.com/supabase/auth/internal/api/scim/protocol"
	"github.com/supabase/auth/internal/models"
)

var providerKey = NewKey[*models.SSOProvider]("sso_provider")

func (srv *Server) Authenticate(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	ctx := r.Context()

	token, ok := parseBearerToken(r.Header.Get("Authorization"))
	if !ok {
		return nil, unauthorized(w)
	}

	provider, err := models.FindSSOProviderBySCIMToken(srv.db.WithContext(ctx), token)
	if err != nil {
		if models.IsNotFoundError(err) {
			return nil, unauthorized(w)
		}
		return nil, err
	}

	if !provider.IsEnabled() {
		return nil, protocol.NewError(http.StatusForbidden, "", "SCIM is not available for this provider")
	}

	return providerKey.With(ctx, provider), nil
}

// parseBearerToken returns the token68 of an RFC 6750 Authorization header.
// Per RFC 7235, Section 2.1 the scheme is case insensitive and is separated
// from the token by one or more spaces.
func parseBearerToken(header string) (string, bool) {
	scheme, rest, found := strings.Cut(header, " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}

	fields := strings.Fields(rest)
	if len(fields) != 1 {
		return "", false
	}

	return fields[0], true
}

func unauthorized(w http.ResponseWriter) error {
	w.Header().Set("WWW-Authenticate", "Bearer")
	return protocol.NewError(http.StatusUnauthorized, "", "A valid SCIM bearer token is required")
}
