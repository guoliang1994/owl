package utils

import (
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func TestMap(t *testing.T) {
	a := `[{"k": "xxx", "v": "zzz"}]`
	var b []map[string]string
	err := jsoniter.UnmarshalFromString(a, &b)
	if err != nil {
		return
	}
}
