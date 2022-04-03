package clients

import (
	"encoding/json"
	"io"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

var (
	BadResponseStatusCode = errors.New("bad response status code")
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

func is4xx(status int) bool {
	if 400 <= status && status < 500 {
		return true
	}
	return false
}
func is5xx(status int) bool {
	if 500 <= status && status < 600 {
		return true
	}
	return false
}
