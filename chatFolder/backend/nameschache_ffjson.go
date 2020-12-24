
import (
	"bytes"
	"encoding/json"
	"fmt"
	fflib "github.com/pquerna/ffjson/fflib/v1"
)

func (mj *NamesOut) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	buf.Grow(128)
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (mj *NamesOut) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	var err error
	var obj []byte
	var scratch fflib.FormatBitsScratch
	_ = obj
	_ = err
	buf.WriteString(`{"connectioncount":`)
	fflib.FormatBits(&scratch, buf, uint64(mj.Connections), 10, false)
	buf.WriteString(`, "users":`)
	if mj.Users != nil {
		buf.WriteString(`[`)
		for i, v := range mj.Users {
			if i != 0 {
				buf.WriteString(`,`)
			}
			if v != nil {
				/* Falling back. type=main.SimplifiedUser kind=struct */
				obj, err = json.Marshal(v)
				if err != nil {
					return err
				}
				buf.Write(obj)
			} else {
				buf.WriteString(`null`)
			}
		}
		buf.WriteString(`]`)
	} else {
		buf.WriteString(`null`)
	}
	buf.WriteString(`, `)
	buf.Rewind(2)
	buf.WriteByte('}')
	return nil
}

const (
	ffj_t_NamesOutbase = iota
	ffj_t_NamesOutno_such_key

	ffj_t_NamesOut_Connections

	ffj_t_NamesOut_Users
)
var ffj_key_NamesOut_Connections = []byte("connectioncount")

var ffj_key_NamesOut_Users = []byte("users")

func (uj *NamesOut) UnmarshalJSON(input []byte) error {
	fs := fflib.NewFFLexer(input)
	return uj.UnmarshalJSONFFLexer(fs, fflib.FFParse_map_start)
}

func (uj *NamesOut) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error = nil
	currentKey := ffj_t_NamesOutbase
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
				currentKey = ffj_t_NamesOutno_such_key
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'c':

					if bytes.Equal(ffj_key_NamesOut_Connections, kn) {
						currentKey = ffj_t_NamesOut_Connections
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'u':

					if bytes.Equal(ffj_key_NamesOut_Users, kn) {
						currentKey = ffj_t_NamesOut_Users
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}
				currentKey = ffj_t_NamesOutno_such_key
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

				case ffj_t_NamesOut_Connections:
					goto handle_Connections

				case ffj_t_NamesOut_Users:
					goto handle_Users

				case ffj_t_NamesOutno_such_key:
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