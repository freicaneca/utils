package httputils

import (
	"bytes"
	"net/http"
	"testing"
	"utils/dto"
	"utils/logging"
	"utils/utils/testutils"
)

type obj struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string `json:"password,omitempty"`
	dto.JSON
}

func (h obj) IsValid() error {
	return nil
}

func TestHTTP(t *testing.T) {

	l := logging.New()

	t.Run("regular bind", func(t *testing.T) {

		inData := `
		{
			"name": "bababobo",
			"age": 53
		}
		`

		req, err := http.NewRequest(
			http.MethodPost,
			"localhost:8888",
			bytes.NewReader([]byte(inData)),
		)
		testutils.AssertBool(t, err == nil, true)
		req.Header.Set("Content-Type", "application/json")

		got := obj{}

		err = BindHTTPRequest(
			l, req, &got, false,
		)
		testutils.AssertBool(t, err == nil, true)

		want := obj{
			Name: "bababobo",
			Age:  53,
		}

		testutils.AssertStruct(
			t, got, want,
		)

	})

}
