package protocol

import (
	"errors"
	"net/http"
	"strconv"
)

// ScimType is a detail error keyword from RFC 7644, Table 9.
type ScimType string

const (
	ScimTypeInvalidFilter ScimType = "invalidFilter"
	ScimTypeInvalidPath   ScimType = "invalidPath"
	ScimTypeInvalidSyntax ScimType = "invalidSyntax"
	ScimTypeInvalidValue  ScimType = "invalidValue"
	ScimTypeInvalidVers   ScimType = "invalidVers"
	ScimTypeMutability    ScimType = "mutability"
	ScimTypeNoTarget      ScimType = "noTarget"
	ScimTypeSensitive     ScimType = "sensitive"
	ScimTypeTooMany       ScimType = "tooMany"
	ScimTypeUniqueness    ScimType = "uniqueness"
)

// Error is the error message form defined in RFC 7644, Section 3.12.
type Error struct {
	Schemas  []string `json:"schemas"`
	ScimType ScimType `json:"scimType,omitempty"`
	Detail   string   `json:"detail,omitempty"`
	Status   string   `json:"status"`

	// internal is the cause a 500 was raised for. It reaches the log line
	// through Error and never the response body.
	internal error
}

func NewError(status int, scimType ScimType, detail string) *Error {
	return &Error{
		Schemas:  []string{SchemaError},
		ScimType: scimType,
		Detail:   detail,
		Status:   strconv.Itoa(status),
	}
}

// Wrap restates err in the RFC 7644 error form so a SCIM client is never
// answered in another dialect. An err that already is one is returned as is,
// and a nil err stays nil.
func Wrap(err error) error {
	var scimErr *Error
	if err == nil || errors.As(err, &scimErr) {
		return err
	}

	wrapped := NewError(http.StatusInternalServerError, "", "Internal server error")
	wrapped.internal = err

	return wrapped
}

func (e *Error) StatusCode() int {
	status, err := strconv.Atoi(e.Status)
	if err != nil {
		return http.StatusInternalServerError
	}
	return status
}

func (e *Error) Error() string {
	detail := e.Detail
	if detail == "" {
		detail = http.StatusText(e.StatusCode())
	}

	if e.internal != nil {
		detail += ": " + e.internal.Error()
	}
	return e.Status + ": " + detail
}

func (e *Error) Unwrap() error {
	return e.internal
}
