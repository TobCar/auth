package scim

import (
	"context"
	"net/http"
	"strings"

	"github.com/supabase/auth/internal/api/scim/protocol"
	"github.com/supabase/auth/internal/models"
	"github.com/supabase/auth/internal/observability"
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
		return nil, protocol.Wrap(err)
	}

	if !provider.IsEnabled() {
		return nil, protocol.NewError(http.StatusForbidden, "", "SCIM is not available for this provider")
	}

	observability.LogEntrySetField(r, "sso_provider_id", provider.ID.String())

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

	token := strings.TrimSpace(rest)
	if token == "" || strings.ContainsAny(token, " \t\r\n\v\f") {
		return "", false
	}

	return token, true
}

func unauthorized(w http.ResponseWriter) error {
	w.Header().Set("WWW-Authenticate", "Bearer")
	return protocol.NewError(http.StatusUnauthorized, "", "A valid SCIM bearer token is required")
}
