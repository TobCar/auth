package core

type SupportedFeature struct {
	Supported bool `json:"supported"`
}

type BulkFeature struct {
	Supported      bool `json:"supported"`
	MaxOperations  int  `json:"maxOperations"`
	MaxPayloadSize int  `json:"maxPayloadSize"`
}

type FilterFeature struct {
	Supported  bool `json:"supported"`
	MaxResults int  `json:"maxResults"`
}

type AuthenticationScheme struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SpecURI     string `json:"specUri,omitempty"`
	Primary     bool   `json:"primary"`
}

// OAuthBearerToken is the scheme described in RFC 7643, Section 8.5.
func OAuthBearerToken() AuthenticationScheme {
	return AuthenticationScheme{
		Type:        "oauthbearertoken",
		Name:        "OAuth Bearer Token",
		Description: "Authentication scheme using the OAuth Bearer Token Standard",
		SpecURI:     "http://www.rfc-editor.org/info/rfc6750",
	}
}

func (scheme AuthenticationScheme) AsPrimary() AuthenticationScheme {
	scheme.Primary = true
	return scheme
}

// ServiceProviderConfig is the schema defined in RFC 7643, Section 5.
type ServiceProviderConfig struct {
	Schemas               []SchemaURI            `json:"schemas"`
	Patch                 SupportedFeature       `json:"patch"`
	Bulk                  BulkFeature            `json:"bulk"`
	Filter                FilterFeature          `json:"filter"`
	ChangePassword        SupportedFeature       `json:"changePassword"`
	Sort                  SupportedFeature       `json:"sort"`
	ETag                  SupportedFeature       `json:"etag"`
	AuthenticationSchemes []AuthenticationScheme `json:"authenticationSchemes"`
	Meta                  Meta                   `json:"meta"`
}

func NewServiceProviderConfig(baseURL string, schemes []AuthenticationScheme) *ServiceProviderConfig {
	return &ServiceProviderConfig{
		Schemas:               []SchemaURI{SchemaServiceProviderConfig},
		AuthenticationSchemes: append(make([]AuthenticationScheme, 0, len(schemes)), schemes...),
		Meta:                  NewMeta(baseURL, ResourceTypeServiceProviderConfig, EndpointServiceProviderConfig, ""),
	}
}

func (config *ServiceProviderConfig) SupportsPatch() *ServiceProviderConfig {
	config.Patch.Supported = true
	return config
}

func (config *ServiceProviderConfig) SupportsBulk(maxOperations, maxPayloadSize int) *ServiceProviderConfig {
	config.Bulk = BulkFeature{Supported: true, MaxOperations: maxOperations, MaxPayloadSize: maxPayloadSize}
	return config
}

func (config *ServiceProviderConfig) SupportsFilter(maxResults int) *ServiceProviderConfig {
	config.Filter = FilterFeature{Supported: true, MaxResults: maxResults}
	return config
}

func (config *ServiceProviderConfig) SupportsChangePassword() *ServiceProviderConfig {
	config.ChangePassword.Supported = true
	return config
}

func (config *ServiceProviderConfig) SupportsSort() *ServiceProviderConfig {
	config.Sort.Supported = true
	return config
}

func (config *ServiceProviderConfig) SupportsETag() *ServiceProviderConfig {
	config.ETag.Supported = true
	return config
}
