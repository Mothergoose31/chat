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