package usecase

type SmsService interface {
	SendSingle(fromNumber string, toNumber string, text string) error
}

type SmsServiceImpl struct {
}

func NewSmsService() SmsService {
	return SmsServiceImpl{}
}

func (s SmsServiceImpl) SendSingle(fromNumber string, toNumber string, text string) error {
	// Todo: must call mock api for sending sms
	// s.sms.send(phoneNumber, text)
	println("Sending Text `", text, "`", " to `", toNumber, "` from number:", fromNumber)
	return nil
}
