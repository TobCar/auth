package core

const (
	SchemaUser       = "urn:ietf:params:scim:schemas:core:2.0:User"
	ResourceTypeUser = "User"
)

type Email struct {
	Value   string `json:"value"`
	Primary bool   `json:"primary"`
}

// User is the core User resource defined in RFC 7643, Section 4.1.
type User struct {
	Schemas  []string `json:"schemas"`
	ID       string   `json:"id"`
	UserName string   `json:"userName"`
	Emails   []Email  `json:"emails,omitempty"`
	Meta     Meta     `json:"meta"`
}
