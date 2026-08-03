package core

import "time"

// Meta is the resource metadata common attribute defined in RFC 7643, Section 3.1.
type Meta struct {
	ResourceType ResourceTypeName `json:"resourceType"`
	Created      time.Time        `json:"created,omitzero"`
	LastModified time.Time        `json:"lastModified,omitzero"`
	Location     string           `json:"location,omitempty"`
}

// NewMeta locates one resource of a collection at baseURL + endpoint + "/" +
// id. An empty id means the resource is its own endpoint, as
// ServiceProviderConfig is.
func NewMeta(baseURL string, resourceType ResourceTypeName, endpoint, id string) Meta {
	location := baseURL + endpoint
	if id != "" {
		location += "/" + id
	}

	return Meta{
		ResourceType: resourceType,
		Location:     location,
	}
}
