package userstore

import "context"

type Store interface {
	Find(ctx context.Context, in FindInput) (out []User, total int64, err error)
	Persist(ctx context.Context, in *User) (err error)
}

type FindInput struct {
	PagingLimit  int `validate:"required"`
	PagingOffset int `validate:"gte=0"`

	ID string `validate:"-"`
}

// ***

type User struct {
	ID             string `validate:"required"`
	Username       string `validate:"-"`
	DisplayName    string `validate:"-"`
	DisplayProfile string `validate:"-"`
}
