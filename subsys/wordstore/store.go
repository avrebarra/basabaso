package wordstore

import "context"

type Store interface {
	Find(ctx context.Context, in FindInput) (out []Word, total int64, err error)
	Persist(ctx context.Context, in *Word) (err error)
}

type FindInput struct {
	PagingLimit  int `validate:"required"`
	PagingOffset int `validate:"gte=0"`

	ID string `validate:"-"`
}

// ***

type Word struct {
	ID   string `validate:"required"`
	Word string `validate:"-"`
}
