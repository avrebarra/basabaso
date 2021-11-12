package wordservice

import (
	"context"
	"fmt"

	"github.com/avrebarra/basabaso/pkg/jsoncast"
	"github.com/avrebarra/basabaso/subsys/wordstore"
)

//go:generate moq -stub -out mocks_test.go -pkg wordservice_test . Store

type (
	Store wordstore.Store
)

// ***

type Service interface {
	FindWords(ctx context.Context, in FindWordsInput) (out FindWordsOutput, err error)
	PersistWord(ctx context.Context, in Word) (err error)
}

type FindWordsInput struct {
	Limit  int    `validate:"-"`
	Offset int    `validate:"-"`
	ID     string `validate:"-"`
}
type FindWordsOutput struct {
	Total int64
	Users []Word
}

// ***

type Word struct {
	ID   string `validate:"required"`
	Word string `validate:"-"`
}

func (Word) FromStoreData(in wordstore.Word) (out Word) {
	if err := jsoncast.CastWithOpts(in, &out, jsoncast.CastOpts{Strict: true}); err != nil {
		panic(fmt.Errorf("bad casting: %w", err))
	}
	return
}

func (in Word) ToStoreData() (out wordstore.Word) {
	if err := jsoncast.CastWithOpts(in, &out, jsoncast.CastOpts{Strict: true}); err != nil {
		panic(fmt.Errorf("bad casting: %w", err))
	}
	return
}
