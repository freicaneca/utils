package nextcloud

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"utils/logging"

	"github.com/gorilla/schema"
)

type handlerOCS struct {
	url     string
	admin   string
	adminPW string
}

func NewHandlerOCS(
	url string,
	admin string,
	adminPW string,
) (
	*handlerOCS,
	error,
) {

	if url == "" {
		return nil, errors.New("empty url")
	}

	if admin == "" {
		return nil, errors.New("empty admin")
	}

	if adminPW == "" {
		return nil, errors.New("empty admin pw")
	}

	return &handlerOCS{
		url:     url,
		admin:   admin,
		adminPW: adminPW,
	}, nil

}

func (h *handlerOCS) RegisterUser(
	log *logging.Logger,
	ctx context.Context,
	userID string,
	password string,
	//email string,
) error {

	l := log.New()

	req := registerUserReq{
		UserID:   userID,
		Password: password,
		//Email:    email,
	}

	err := req.IsValid()
	if err != nil {
		return fmt.Errorf("invalid req: %v: %w",
			err, ErrBadRequest)
	}

	usersURL := h.url + "/ocs/v2.php/cloud/users"

	l.Debug("full url: %v", usersURL)

	enc := schema.NewEncoder()

	form := url.Values{}

	err = enc.Encode(req, form)
	if err != nil {
		return fmt.Errorf("encoding form values: %v: %w",
			err, ErrInternal)
	}

	hReq, err := http.NewRequest(
		http.MethodPost,
		usersURL,
		bytes.NewBufferString(form.Encode()),
	)
	if err != nil {
		return fmt.Errorf("creating new http req: %v: %w",
			err, ErrInternal)
	}

	hReq.Header.Set("OCS-APIRequest", "true")
	hReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	hReq.SetBasicAuth(h.admin, h.adminPW)

	client := http.Client{}

	resp, err := client.Do(hReq)
	if err != nil {
		return fmt.Errorf("sending req: %v: %w",
			err, ErrInternal)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	l.Debug("raw resp body: %+v", string(body))

	var ocsResp registerUserResp

	err = xml.Unmarshal(body, &ocsResp)
	if err != nil {
		return fmt.Errorf("unmarshalling xml: %v: %w",
			err, ErrInternal)
	}

	l.Debug("unserialized resp: %+v", ocsResp)

	code := ocsResp.Meta.StatusCode
	l.Info("got status code %v", code)
	if code != 100 && code != 200 {
		switch code {
		case 107:
			err = ErrBadPassword
		case 102:
			err = ErrUserAlreadyExists
		default:
			err = ErrInternal
		}
		return fmt.Errorf("bad status: %v: %w",
			ocsResp.Meta.Status, err)
	}

	l.Info("created user %v",
		userID)

	return nil

}

func (h *handlerOCS) RemoveUser(
	log *logging.Logger,
	ctx context.Context,
	userID string,
) error {

	l := log.New()

	if userID == "" {
		return fmt.Errorf("empty user id: %w",
			ErrBadRequest)
	}

	usersURL := h.url + "/ocs/v2.php/cloud/users/" + userID

	l.Debug("full url: %v", usersURL)

	hReq, err := http.NewRequest(
		http.MethodDelete,
		usersURL,
		nil,
	)
	if err != nil {
		return fmt.Errorf("creating new http req: %v: %w",
			err, ErrInternal)
	}

	hReq.Header.Set("OCS-APIRequest", "true")
	hReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	hReq.SetBasicAuth(h.admin, h.adminPW)

	client := http.Client{}

	resp, err := client.Do(hReq)
	if err != nil {
		return fmt.Errorf("sending req: %v: %w",
			err, ErrInternal)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	l.Debug("raw resp body: %+v", string(body))

	var ocsResp removeUserResp

	err = xml.Unmarshal(body, &ocsResp)
	if err != nil {
		return fmt.Errorf("unmarshalling xml: %v: %w",
			err, ErrInternal)
	}

	l.Debug("unserialized resp: %+v", ocsResp)

	code := ocsResp.Meta.StatusCode
	l.Info("got status code %v", code)
	if code != 100 && code != 200 {
		return fmt.Errorf("bad status: %v: %w",
			ocsResp.Meta.Status, ErrInternal)
	}

	l.Info("removed user %v",
		userID)

	return nil

}
