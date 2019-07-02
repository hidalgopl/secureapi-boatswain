package testjson

import (
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
)

func ToJSON(t *testing.T, v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(errors.Wrap(err, "cannot marshal given value"))
	}
	return string(b)
}
