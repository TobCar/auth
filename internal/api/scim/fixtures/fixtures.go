// Package fixtures holds the expected SCIM wire payloads shared by tests.
package fixtures

import _ "embed"

//go:embed method_not_allowed.json
var MethodNotAllowed string

//go:embed not_found.json
var NotFound string

//go:embed not_implemented.json
var NotImplemented string
