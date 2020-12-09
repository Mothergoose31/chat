package main

import (
	"regexp"
	"encoding json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Api struct {
	url string
	key string
}

var (
	cookievalid = regexp.MustCompile("^[A-z0-9]{10,64}$")
	api         Api
)
func initApi(url, key string) {
	api = Api{
		url: url,
		key: key,
	}
}


func (a *Api) getUserFromAuthToken(tok string) ([]byte, error) {
	if !cookievalid.MatchString(tok) {
		return nil, fmt.Errorf("api: auth token cookie invalid %s", tok)
	}

	endpoint := a.url + "/auth"
	resp, err := http.PostForm(endpoint, url.Values{
		"authtoken":  {tok},
		"privatekey": {a.key},
	})

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("api: auth token invalid: %s, response code: %d", tok, resp.StatusCode)
	}

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

