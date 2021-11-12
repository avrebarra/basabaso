package fireconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/avrebarra/valeed"
	"google.golang.org/api/option"

	firebase "firebase.google.com/go"
)

type Config struct {
	FirebaseConfigJSON []byte `validate:"required"`
	CollecName         string `validate:"required"`
	ContentFieldName   string `validate:"required"`
}

type FireConfig struct {
	config          Config
	firestoreclient *firestore.Client
}

func New(cfg Config) (*FireConfig, error) {
	if err := valeed.Validate(cfg); err != nil {
		return nil, err
	}

	e := &FireConfig{config: cfg}

	if err := e.init(); err != nil {
		err = fmt.Errorf("bad init: %w", err)
		return nil, err
	}

	return e, nil
}

func (e *FireConfig) init() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// initialize firebase app
	firebaseapp, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(e.config.FirebaseConfigJSON))
	if err != nil {
		err = fmt.Errorf("firebase app init failure: %w", err)
		return
	}

	// initialize firestore client
	firestoreclient, err := firebaseapp.Firestore(ctx)
	if err != nil {
		err = fmt.Errorf("firestore init failure: %w", err)
		return
	}
	e.firestoreclient = firestoreclient

	return
}

func (e *FireConfig) close(ctx context.Context) (err error) {
	if e.firestoreclient != nil {
		e.firestoreclient.Close()
	}
	return
}

// ***

func (e *FireConfig) Get(ctx context.Context, id string) (out []byte, err error) {
	// fetch document
	outGet, err := e.firestoreclient.Collection(e.config.CollecName).Doc(id).Get(ctx)
	if err != nil {
		err = fmt.Errorf("doc fetch failure: %w", err)
		return
	}

	// build output
	outData, err := outGet.DataAt(e.config.ContentFieldName)
	if err != nil {
		err = fmt.Errorf("doc field fetch failure: %w", err)
		return
	}
	out, err = json.Marshal(outData)
	if err != nil {
		err = fmt.Errorf("doc serialize failure: %w", err)
		return
	}

	return
}

func (e *FireConfig) Close(ctx context.Context) (err error) {
	return e.close(ctx)
}
