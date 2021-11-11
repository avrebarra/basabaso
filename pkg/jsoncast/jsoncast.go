package jsoncast

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type CastOpts struct {
	Strict bool
}

func Cast(in interface{}, out interface{}) (err error) {
	return CastWithOpts(in, out, CastOpts{Strict: false})
}

func CastWithOpts(in interface{}, out interface{}, opts CastOpts) (err error) {
	bts, err := json.Marshal(in)
	if err != nil {
		err = fmt.Errorf("serialize failure: %w", err)
		return
	}
	dec := json.NewDecoder(bytes.NewBuffer(bts))
	if opts.Strict {
		dec.DisallowUnknownFields()
	}
	err = dec.Decode(out)
	if err != nil {
		err = fmt.Errorf("deserialize failure: %w", err)
		return
	}
	return
}
