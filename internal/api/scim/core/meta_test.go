package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMeta(t *testing.T) {
	t.Run("serializes to JSON correctly", func(t *testing.T) {
		body, err := json.Marshal(Meta{
			ResourceType: "ServiceProviderConfig",
			Location:     "http://localhost:9999/scim/v2/ServiceProviderConfig",
		})

		require.NoError(t, err)
		require.JSONEq(t, `{
			"resourceType": "ServiceProviderConfig",
			"location": "http://localhost:9999/scim/v2/ServiceProviderConfig"
		}`, string(body))
	})

	t.Run("omits the location when it is empty", func(t *testing.T) {
		body, err := json.Marshal(Meta{ResourceType: "ServiceProviderConfig"})

		require.NoError(t, err)
		require.JSONEq(t, `{"resourceType": "ServiceProviderConfig"}`, string(body))
	})
}
