package userservice

import (
	"context"
	"fmt"

	"github.com/avrebarra/basabaso/subsys/userstore"
	"github.com/avrebarra/valeed"
	"github.com/jinzhu/copier"
)

type Config struct {
	Store Store `validate:"required"`
}

type Default struct {
	config Config
}

func New(cfg Config) (Service, error) {
	if err := valeed.Validate(cfg); err != nil {
		return nil, err
	}
	e := &Default{config: cfg}
	return e, nil
}

func (e *Default) FindUsers(ctx context.Context, in FindUsersInput) (out FindUsersOutput, err error) {
	defaults := FindUsersInput{
		Limit:  10,
		Offset: 0,
	}

	// prep and validate
	if err = valeed.Validate(in); err != nil {
		err = fmt.Errorf("bad input: %w", err)
		return
	}
	if copier.CopyWithOption(&defaults, &in, copier.Option{IgnoreEmpty: true}); err != nil {
		err = fmt.Errorf("value fallback failure: %w", err)
		return
	}
	in = defaults

	// perform operation
	users, tot, err := e.config.Store.Find(ctx, userstore.FindInput{
		PagingLimit:  in.Limit,
		PagingOffset: in.Offset,
		ID:           in.ID,
	})
	if err != nil {
		err = fmt.Errorf("query failure: %w", err)
		return
	}

	// build output
	out = FindUsersOutput{
		Total: tot,
		Users: []User{},
	}
	for _, v := range users {
		out.Users = append(out.Users, User{}.FromStoreData(v))
	}

	return
}

func (e *Default) PersistUser(ctx context.Context, in User) (err error) {
	// prep and validate
	if err = valeed.Validate(in); err != nil {
		err = fmt.Errorf("bad input: %w", err)
		return
	}

	// perform operation
	user := in.ToStoreData()

	err = e.config.Store.Persist(ctx, &user)
	if err != nil {
		err = fmt.Errorf("storing failure: %w", err)
		return
	}

	// build output
	in = User{}.FromStoreData(user)

	return
}
