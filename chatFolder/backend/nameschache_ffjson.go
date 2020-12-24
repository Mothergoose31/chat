
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
