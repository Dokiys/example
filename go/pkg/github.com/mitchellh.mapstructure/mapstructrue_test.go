package mitchellh_mapstructure

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestMapstructure(t *testing.T) {
	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
	}

	// This input can come from anywhere, but typically comes from
	// something like decoding JSON where we're not quite sure of the
	// struct initially.
	input := map[string]interface{}{
		"name":   "Mitchell",
		"age":    91,
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"twitter": "mitchellh",
		},
	}

	var result Person
	err := mapstructure.Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", result)
}

func TestMapstructure_Tags(t *testing.T) {
	type Family struct {
		LastName string
	}
	type Location struct {
		City string
	}
	type Person struct {
		FirstName string
		Family    `mapstructure:",squash"`
		*Location `mapstructure:",omitempty"`
		Other     map[string]interface{} `mapstructure:",remain"`
	}

	input := map[string]interface{}{
		"FirstName":          "Mitchell",
		"LastName":           "Hashimoto",
		"location_omitempty": "",
		"email":              "mitchell@example.com",
	}

	var result Person
	err := mapstructure.Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", result)
}

type Family struct {
	LastName string `json:"last_name"`
}
type Person struct {
	FirstName string `json:"first_name"`
	Family    Family `json:"family"`
}

func (f *Family) Unmarshal(src []byte) error {
	return json.Unmarshal(src, &f)
}

func TestMapstructure_JSON(t *testing.T) {
	var input = map[string]interface{}{
		"first_name": "San",
		"family":     "{\"last_name\":\"Zhang\"}",
	}

	var dest = &Person{}
	config := &mapstructure.DecoderConfig{
		Result:  dest,
		TagName: "json",
		DecodeHook: func(from reflect.Type, to reflect.Type, src interface{}) (dest interface{}, err error) {
			if from.Kind() != reflect.String {
				return src, nil
			}

			vi := reflect.New(to).Interface()
			if vv, ok := vi.(interface{ Unmarshal([]byte) error }); ok {
				if err := vv.Unmarshal([]byte((src.(string)))); err != nil {
					return nil, err
				}

				return vv, nil
			}

			return src, nil
		},
	}

	decoder, err := mapstructure.NewDecoder(config)
	assert.NoError(t, err)
	err = decoder.Decode(input)
	assert.NoError(t, err)

	t.Logf("%+v", dest)
}
