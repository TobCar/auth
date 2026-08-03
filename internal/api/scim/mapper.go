package scim

import (
	"github.com/supabase/auth/internal/api/scim/core"
	"github.com/supabase/auth/internal/models"
)

type Mapper[TIn, TOut any] interface {
	MapFrom(in TIn) TOut
}

type UserMapper struct {
	baseURL string
}

func NewUserMapper(baseURL string) UserMapper {
	return UserMapper{baseURL: baseURL}
}

func (m UserMapper) MapFrom(u *models.User) *core.User {
	id, email := u.ID.String(), u.GetEmail()

	user := &core.User{
		Schemas:  []string{core.SchemaUser},
		ID:       id,
		UserName: email,
		Meta: core.Meta{
			ResourceType: core.ResourceTypeUser,
			Created:      u.CreatedAt.UTC(),
			LastModified: u.UpdatedAt.UTC(),
			Location:     m.baseURL + "/Users/" + id,
		},
	}

	if email != "" {
		user.Emails = []core.Email{{Value: email, Primary: true}}
	}

	return user
}
