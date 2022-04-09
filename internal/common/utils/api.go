package utils

import (
	"encoding/json"
	"io"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

func FromBody(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	return errors.WithStack(decoder.Decode(&dest))
}

func ToBody(body io.Writer, source interface{}) error {
	encoder := json.NewEncoder(body)
	return errors.WithStack(encoder.Encode(source))
}

func BuildURL(base string, api string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", errors.WithStack(err)
	}

	u.Path = path.Join(u.Path, api)
	return u.String(), nil
}

func Is4xx(status int) bool {
	if 400 <= status && status < 500 {
		return true
	}
	return false
}
func Is5xx(status int) bool {
	if 500 <= status && status < 600 {
		return true
	}
	return false
}
