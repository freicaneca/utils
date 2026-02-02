package nextcloud

import (
	"encoding/xml"
	"fmt"
)

type registerUserReq struct {
	UserID   string `schema:"userid"`
	Password string `schema:"password"`
	//Email    string `schema:"email"`
}

func (h *registerUserReq) IsValid() error {
	if h.UserID == "" {
		return fmt.Errorf("empty user id")
	}

	if h.Password == "" {
		return fmt.Errorf("empty user password")
	}

	// if h.Email == "" {
	// 	return fmt.Errorf("empty user email")
	// }

	return nil
}

type registerUserResp struct {
	OCS xml.Name `xml:"ocs"`
	metaResp
	Data struct {
		ID string `xml:"id"`
	} `xml:"data"`
}

type removeUserResp struct {
	metaResp
}

type metaResp struct {
	Meta struct {
		Status     string `xml:"status"`
		StatusCode int    `xml:"statuscode"`
		Message    string `xml:"message"`
	} `xml:"meta"`
}
