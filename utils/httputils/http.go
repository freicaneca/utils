package httputils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"utils/dto"
	"utils/logging"
	"utils/utils/contextutils"
	"utils/utils/trackingutils"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type HTTPError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

// Binds HTTP request to a given struct dtoIn.
// If classified, will not log request details.
func BindHTTPRequest(
	log *logging.Logger,
	req *http.Request,
	dtoIn dto.DTO,
	classified bool,
) error {

	l := log.New()

	if classified {
		l.Info("got classified request. will not log")
	} else {
		l.Info("got req: %+v", *req)
	}

	content := []byte{}

	defer req.Body.Close()

	switch req.Method {
	case "GET", "DELETE":
		field := mux.Vars(req)
		reference := make(map[string]interface{})
		for key, value := range field {
			reference[key] = value
		}

		_ = schema.NewDecoder().Decode(dtoIn, req.URL.Query())

		content, _ = json.Marshal(reference)
	case "PUT", "POST", "PATCH":

		_, ok := req.Header["Content-Type"]
		if !ok {
			l.Error("no Content-Type header")
			return errors.New("no Content-Type header")
		}

		switch req.Header["Content-Type"][0] {
		case "application/x-www-form-urlencoded":
			req.ParseForm()
			hCont := make(map[string]string)
			for key, values := range req.Form {
				for _, value := range values {
					hCont[key] = value
				}
			}
			content, _ = json.Marshal(hCont)
		case "application/json":
			content, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(content))

		}
	}

	if !classified {
		l.Info("endpoint: %v %v | raw body: %+v",
			req.Method,
			req.URL,
			string(content),
		)
	} else {
		l.Info("endpoint: %v %v",
			req.Method,
			req.URL,
		)
	}

	err := dtoIn.ToObject(content, dtoIn)
	if err != nil {
		errMsg := ""
		if !classified {
			errMsg = fmt.Sprintf("error unmarshalling %v to %+v: %v",
				string(content), dtoIn, err)
		} else {
			errMsg = "error unmarshalling"
		}
		l.Error(errMsg)
		return errors.New(errMsg)
	}

	err = dtoIn.IsValid()
	if err != nil {
		l.Error("dto not valid! %+v", err)
		return err
	}

	return nil

}

func GetJWTFromHeader(r *http.Request) string {

	token := r.Header.Get("Authorization")

	fields := strings.Split(token, " ")

	if len(fields) != 2 {
		return ""
	}

	return fields[1]
}

func SetupCorsResponse(w http.ResponseWriter) {
	// settings cors allows
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods",
		"POST, GET, OPTIONS, PUT, DELETE, PATCH")
	w.Header().Set("Access-Control-Allow-Headers",
		"*",
	)
	// w.Header().Set("Access-Control-Allow-Headers",
	// 	"X-Requested-With, Content-Type, AccessToken, "+
	// 		"CsrfToken, AppVersion, Device, Authorization, "+
	// 		"Set-Cookie, "+
	// 		stringutils.ImplodeString(", ", headers[:]))

}

func WriteHTTPError(
	w http.ResponseWriter,
	err *HTTPError,
	status int,
) {
	bodyB, _ := json.Marshal(err)

	w.WriteHeader(status)
	w.Write(bodyB)
}

func SetTrackingIDFromHeader(
	r *http.Request,
) *http.Request {

	ctx := r.Context()

	trackingID := r.Header.Get(HTTPHeaderTrackingID)
	if trackingID == "" {
		trackingID = trackingutils.GlobalTrackingNumber.Next()
	}

	ctx = context.WithValue(
		ctx, contextutils.ContextKeyReqTracking, trackingID)

	r = r.WithContext(ctx)

	return r
}
