package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServiceProviderConfig(t *testing.T) {
	t.Run("advertises the schemes the caller declares", func(t *testing.T) {
		schemes := []AuthenticationScheme{OAuthBearerToken().AsPrimary()}

		config := NewServiceProviderConfig("", schemes)

		require.Equal(t, []string{SchemaServiceProviderConfig}, config.Schemas)
		require.Equal(t, schemes, config.AuthenticationSchemes)
	})

	t.Run("identifies itself with resource metadata", func(t *testing.T) {
		baseURL := "http://localhost:9999/scim/v2"

		config := NewServiceProviderConfig(baseURL, nil)

		require.Equal(t, "ServiceProviderConfig", config.Meta.ResourceType)
		require.Equal(t, baseURL+"/ServiceProviderConfig", config.Meta.Location)
	})

	t.Run("supports none of the optional protocol features", func(t *testing.T) {
		config := NewServiceProviderConfig("", nil)

		assert.False(t, config.Patch.Supported)
		assert.False(t, config.Bulk.Supported)
		assert.False(t, config.Filter.Supported)
		assert.False(t, config.ChangePassword.Supported)
		assert.False(t, config.Sort.Supported)
		assert.False(t, config.ETag.Supported)
	})

	t.Run("serializes authenticationSchemes as an array", func(t *testing.T) {
		body, err := json.Marshal(NewServiceProviderConfig("", nil))

		require.NoError(t, err)
		require.Contains(t, string(body), `"authenticationSchemes":[]`)
	})

	t.Run("does not alias the caller's slice", func(t *testing.T) {
		schemes := []AuthenticationScheme{OAuthBearerToken()}

		config := NewServiceProviderConfig("", schemes)
		schemes[0].Name = "mutated"

		require.Equal(t, "OAuth Bearer Token", config.AuthenticationSchemes[0].Name)
	})
}

func TestAuthenticationScheme(t *testing.T) {
	t.Run("OAuthBearerToken", func(t *testing.T) {
		scheme := OAuthBearerToken()

		assert.Equal(t, "oauthbearertoken", scheme.Type)
		assert.Equal(t, "OAuth Bearer Token", scheme.Name)
		assert.Equal(t, "Authentication scheme using the OAuth Bearer Token Standard", scheme.Description)
		assert.Equal(t, "http://www.rfc-editor.org/info/rfc6750", scheme.SpecURI)
	})

	t.Run("AsPrimary returns a copy", func(t *testing.T) {
		scheme := OAuthBearerToken()

		assert.False(t, scheme.Primary)
		assert.True(t, scheme.AsPrimary().Primary)
		assert.False(t, scheme.Primary)
	})
}
