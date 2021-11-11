package wordservice

import (
	"context"
	"fmt"

	"github.com/avrebarra/basabaso/subsys/wordstore"
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

func (e *Default) FindWords(ctx context.Context, in FindWordsInput) (out FindWordsOutput, err error) {
	defaults := FindWordsInput{
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
	users, tot, err := e.config.Store.Find(ctx, wordstore.FindInput{
		PagingLimit:  in.Limit,
		PagingOffset: in.Offset,
		ID:           in.ID,
	})
	if err != nil {
		err = fmt.Errorf("query failure: %w", err)
		return
	}

	// build output
	out = FindWordsOutput{
		Total: tot,
		Users: []Word{},
	}
	for _, v := range users {
		out.Users = append(out.Users, Word{}.FromStoreData(v))
	}

	return
}

func (e *Default) PersistWord(ctx context.Context, in Word) (err error) {
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
	in = Word{}.FromStoreData(user)

	return
}
