package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
)

func BindJSON(file string, v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("`v` must be a pointer")
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, v)
	return err
}
