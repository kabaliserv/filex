package sql

import (
	"bytes"
	"encoding/json"
	"reflect"
)

type dataCompare map[string]interface{}

func getChangedValue(old dataCompare, new dataCompare) dataCompare {
	var updateValue = make(map[string]interface{})

	for k, v := range new {
		if !reflect.DeepEqual(old[k], v) {
			updateValue[k] = v
		}
	}

	return updateValue
}

type getValueFunc func() map[string]interface{}
type saveFunc func(map[string]interface{}) error

func getSaveChangeFunc(save saveFunc, value getValueFunc) func() error {
	oldVal := value()
	return func() error {
		newVal := value()
		return save(getChangedValue(oldVal, newVal))
	}
}

func structToStruct(from interface{}, to interface{}) error {
	v, err := json.Marshal(from)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(v)
	dec := json.NewDecoder(reader)
	return dec.Decode(to)
}
