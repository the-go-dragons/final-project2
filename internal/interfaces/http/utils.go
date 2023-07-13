package http

import (
	"errors"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Response struct {
	Message string `json:"message"`
}

func CheckTheNumberFormat(number string) error {
	if !govalidator.Matches(number, `^\+98\d{10}$`) {
		return errors.New("")
	}
	return nil
}

func ValidateSingleSMSBody(sender string, reciever string, content string) error {
	if CheckTheNumberFormat(sender) != nil {
		return errors.New("invalid sender number")
	}
	if CheckTheNumberFormat(reciever) != nil {
		return errors.New("invalid receiver number")
	}
	if strings.EqualFold(reciever, sender) {
		return errors.New("impossible to send sms to yourself")
	}
	if len(strings.Trim(content, " ")) == 0 {
		return errors.New("invalid content")
	}
	return nil
}
