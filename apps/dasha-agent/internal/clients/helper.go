package clients

import (
	"encoding/json"
	"io"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

func fromBody(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	return errors.WithStack(decoder.Decode(&dest))
}

func buildURL(base string, api string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", errors.WithStack(err)
	}

	u.Path = path.Join(u.Path, api)
	return u.String(), nil
}
