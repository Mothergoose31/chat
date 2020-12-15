package main

import (
	"encoding/json"
	"errors"
	"strings"
)

func Unpack(data string) (string, []byte, error) {
	result := strings.SplitN(data, " ", 2)
	if len(result) != 2 {
		return "", nil, errors.New("Unable to extract event name from data.")
	}
	return result[0], []byte(result[1]), nil
}
