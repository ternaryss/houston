package mocks

import (
	"bytes"
	"net/http"
	"net/url"
)

func MockRequest(mth string, dat url.Values) (*http.Request, error) {
	req, err := http.NewRequest(mth, "", bytes.NewBufferString(dat.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}
