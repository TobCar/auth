// Package core implements the SCIM 2.0 core schema defined in RFC 7643.
package core

// SchemaURI identifies a SCIM schema, e.g.
// "urn:ietf:params:scim:schemas:core:2.0:User".
type SchemaURI string

// ResourceTypeName names a resource type, e.g. "User". It is the value of
// meta.resourceType on every resource of that type.
type ResourceTypeName string
