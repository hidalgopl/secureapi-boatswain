package goldenfile

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
)

var update = flag.Bool("test.update", false, "update golden files")

func Get(t *testing.T, actual string, filename string) string {
	expected := ReadFile(t, []byte(actual), filename)
	return string(expected)
}

func ReadFile(t *testing.T, actual []byte, filename string) []byte {
	golden := filepath.Join("testdata", filename)
	if *update {
		err := ioutil.WriteFile(golden, actual, 0644)
		if err != nil {
			t.Fatal(errors.Wrap(err, "cannot update golden file: "+filename))
		}
	}
	expected, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Fatal(errors.Wrap(err, "cannot read golden file: "+filename))
	}
	return expected
}
