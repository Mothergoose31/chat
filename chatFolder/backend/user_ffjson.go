package main

import (
	"bytes"
	"fmt"
	fflib "github.com/pquerna/ffjson/fflib/v1"
)

func (mj *sessionuser) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	buf.Grow(256)
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (mj *sessionuser) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"features":`)
	if mj.Features != nil {
		buf.WriteString(`[`)
		for i, v := range mj.Features {
			if i != 0 {
				buf.WriteString(`,`)
			}
			fflib.WriteJsonString(buf, string(v))
		}
		buf.WriteString(`]`)
	} else {
		buf.WriteString(`null`)
	}
	buf.WriteString(`, "userId":`)
	fflib.WriteJsonString(buf, string(mj.UserId))
	buf.WriteString(`, "username":`)
	fflib.WriteJsonString(buf, string(mj.Username))
	buf.WriteString(`, `)
	buf.Rewind(2)
	buf.WriteByte('}')
	return nil
}
