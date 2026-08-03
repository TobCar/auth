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
	created, lastModified := u.CreatedAt.UTC(), u.UpdatedAt.UTC()

	user := &core.User{
		Schemas:  []string{core.SchemaUser},
		ID:       u.ID.String(),
		UserName: u.GetEmail(),
		Meta: core.Meta{
			ResourceType: core.ResourceTypeUser,
			Created:      &created,
			LastModified: &lastModified,
			Location:     m.baseURL + "/Users/" + u.ID.String(),
		},
	}

	if email := u.GetEmail(); email != "" {
		user.Emails = []core.Email{{Value: email, Primary: true}}
	}

	return user
}
