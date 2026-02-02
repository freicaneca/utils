package testutils

import (
	"errors"
	"reflect"
	"testing"
)

func AssertStruct(t *testing.T, got, want interface{}) bool {

	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("\n\nObteve:\n\n\t%+v\n\nDesejado:\n\n\t%+v\n\n",
			got, want)
		return false
	}

	return true
}

func AssertString(t *testing.T, got, want string) bool {

	t.Helper()

	if got != want {
		t.Fatalf("\n\nObteve:\n\n\t%s\n\nDesejado: \n\n\t%s\n\n",
			got, want)
		return false
	}

	return true
}

func AssertInt(t *testing.T, got, want int) bool {

	t.Helper()

	if got != want {
		t.Fatalf("\n\nObteve:\n\n\t%v\n\nDesejado: \n\n\t%v\n\n",
			got, want)
		return false
	}

	return true
}

func AssertBool(t *testing.T, got, want bool) bool {

	t.Helper()

	if got != want {
		t.Fatalf("\n\nObteve:\n\n\t%v\n\nDesejado: \n\n\t%v\n\n",
			got, want)
		return false
	}

	return true
}

func AssertError(t *testing.T, got, want error) bool {

	t.Helper()

	if !errors.Is(got, want) {
		t.Fatalf("\n\nObteve:\n\n\t%v\n\nDesejado: \n\n\t%v\n\n",
			got, want)
		return false
	}

	return true

}

func AssertFloat(t *testing.T, got, want float64) bool {

	t.Helper()

	if got != want {
		t.Fatalf("\n\nObteve:\n\n\t%v\n\nDesejado: \n\n\t%v\n\n",
			got, want)
		return false
	}

	return true
}
