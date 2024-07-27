package json

import (
	"encoding/json"
	"golang.org/x/xerrors"
	"io"
)

func Decode[T any](r io.Reader) (T, error) {
	var v T
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return v, xerrors.Errorf("json.Decoder.Decode: %w", err)
	}
	return v, nil
}
