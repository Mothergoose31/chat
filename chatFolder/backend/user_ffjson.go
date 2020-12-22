
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
const (
	ffj_t_sessionuserbase = iota
	ffj_t_sessionuserno_such_key

	ffj_t_sessionuser_Features

	ffj_t_sessionuser_UserId

	ffj_t_sessionuser_Username
)

var ffj_key_sessionuser_Features = []byte("features")

var ffj_key_sessionuser_UserId = []byte("userId")

var ffj_key_sessionuser_Username = []byte("username")

func (uj *sessionuser) UnmarshalJSON(input []byte) error {
	fs := fflib.NewFFLexer(input)
	return uj.UnmarshalJSONFFLexer(fs, fflib.FFParse_map_start)
}

func (uj *sessionuser) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error = nil
	currentKey := ffj_t_sessionuserbase
	_ = currentKey
	tok := fflib.FFTok_init
	wantedTok := fflib.FFTok_init

mainparse:
	for {
		tok = fs.Scan()
		//	println(fmt.Sprintf("debug: tok: %v  state: %v", tok, state))
		if tok == fflib.FFTok_error {
			goto tokerror
		}

		switch state {

		case fflib.FFParse_map_start:
			if tok != fflib.FFTok_left_bracket {
				wantedTok = fflib.FFTok_left_bracket
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_key
			continue

		case fflib.FFParse_after_value:
			if tok == fflib.FFTok_comma {
				state = fflib.FFParse_want_key
			} else if tok == fflib.FFTok_right_bracket {
				goto done
			} else {
				wantedTok = fflib.FFTok_comma
				goto wrongtokenerror
			}

		case fflib.FFParse_want_key:
			// json {} ended. goto exit. woo.
			if tok == fflib.FFTok_right_bracket {
				goto done
			}
			if tok != fflib.FFTok_string {
				wantedTok = fflib.FFTok_string
				goto wrongtokenerror
			}

			kn := fs.Output.Bytes()
			if len(kn) <= 0 {
				// "" case. hrm.
				currentKey = ffj_t_sessionuserno_such_key
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'f':

					if bytes.Equal(ffj_key_sessionuser_Features, kn) {
						currentKey = ffj_t_sessionuser_Features
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'u':

					if bytes.Equal(ffj_key_sessionuser_UserId, kn) {
						currentKey = ffj_t_sessionuser_UserId
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffj_key_sessionuser_Username, kn) {
						currentKey = ffj_t_sessionuser_Username
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}
				currentKey = ffj_t_sessionuserno_such_key
				state = fflib.FFParse_want_colon
				goto mainparse
			}

		case fflib.FFParse_want_colon:
			if tok != fflib.FFTok_colon {
				wantedTok = fflib.FFTok_colon
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_value
			continue
		case fflib.FFParse_want_value:

			if tok == fflib.FFTok_left_brace || tok == fflib.FFTok_left_bracket || tok == fflib.FFTok_integer || tok == fflib.FFTok_double || tok == fflib.FFTok_string || tok == fflib.FFTok_bool || tok == fflib.FFTok_null {
				switch currentKey {

				case ffj_t_sessionuser_Features:
					goto handle_Features

				case ffj_t_sessionuser_UserId:
					goto handle_UserId

				case ffj_t_sessionuser_Username:
					goto handle_Username

				case ffj_t_sessionuserno_such_key:
					err = fs.SkipField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}
					state = fflib.FFParse_after_value
					goto mainparse
				}
			} else {
				goto wantedvalue
			}
		}
	}
	handle_Features:

	/* handler: uj.Features type=[]string kind=slice */

	{

		{
			if tok != fflib.FFTok_left_brace && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			uj.Features = nil
		} else {

			uj.Features = make([]string, 0)

		}

		wantVal := true

		for {

			var v string

			tok = fs.Scan()
			if tok == fflib.FFTok_error {
				goto tokerror
			}
			if tok == fflib.FFTok_right_brace {
				break
			}

			if tok == fflib.FFTok_comma {
				if wantVal == true {
					// TODO(pquerna): this isn't an ideal error message, this handles
					// things like [,,,] as an array value.
					return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
				}
				continue
			} else {
				wantVal = true
			}

			/* handler: v type=string kind=string */

			{

				{
					if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
						return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
					}
				}

				if tok == fflib.FFTok_null {

				} else {

					v = string(fs.Output.String())

				}
			}

			uj.Features = append(uj.Features, v)
			wantVal = false
		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

	handle_UserId:

	/* handler: uj.UserId type=string kind=string */

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			uj.UserId = string(fs.Output.String())

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Username:

	/* handler: uj.Username type=string kind=string */

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			uj.Username = string(fs.Output.String())

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

wantedvalue:
	return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
wrongtokenerror:
	return fs.WrapErr(fmt.Errorf("ffjson: wanted token: %v, but got token: %v output=%s", wantedTok, tok, fs.Output.String()))
tokerror:
	if fs.BigError != nil {
		return fs.WrapErr(fs.BigError)
	}
	err = fs.Error.ToError()
	if err != nil {
		return fs.WrapErr(err)
	}
	panic("ffjson-generated: unreachable, please report bug.")
done:
	return nil
}