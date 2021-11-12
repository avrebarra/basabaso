package userservice

import (
	"context"
	"fmt"

	"github.com/avrebarra/basabaso/pkg/jsoncast"
	"github.com/avrebarra/basabaso/subsys/userstore"
)

//go:generate moq -stub -out mocks_test.go -pkg userservice_test . Store

type (
	Store userstore.Store
)

// ***

type Service interface {
	FindUsers(ctx context.Context, in FindUsersInput) (out FindUsersOutput, err error)
	PersistUser(ctx context.Context, in User) (err error)
}

type FindUsersInput struct {
	Limit  int    `validate:"-"`
	Offset int    `validate:"-"`
	ID     string `validate:"-"`
}
type FindUsersOutput struct {
	Total int64
	Users []User
}

// ***

type User struct {
	ID             string `validate:"required"`
	Username       string `validate:"-"`
	DisplayName    string `validate:"-"`
	DisplayProfile string `validate:"-"`
}

func (User) FromStoreData(in userstore.User) (out User) {
	if err := jsoncast.CastWithOpts(in, &out, jsoncast.CastOpts{Strict: true}); err != nil {
		panic(fmt.Errorf("bad casting: %w", err))
	}
	return
}

func (in User) ToStoreData() (out userstore.User) {
	if err := jsoncast.CastWithOpts(in, &out, jsoncast.CastOpts{Strict: true}); err != nil {
		panic(fmt.Errorf("bad casting: %w", err))
	}
	return
}
