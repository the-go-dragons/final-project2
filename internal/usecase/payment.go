package usecase

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type CallbackResult struct {
	Successful bool
	PaymentID  int
}
type PaymentService struct {
	paymentRepo persistence.PaymentRepository
}

func (PaymentService) GetGateway(bank Bank) Gateway {
	return BANKS[bank]()
}
func NewPayment(
	paymentRepo persistence.PaymentRepository,
) PaymentService {
	return PaymentService{paymentRepo: paymentRepo}
}

func (p PaymentService) GetPaymentPage(paymentID int, bankName string) (PaymentPage, error) {
	bank, err := validateAndGetBank(bankName)
	if err != nil {
		return PaymentPage{}, InvalidBankName{bankName}
	}
	payment, err := p.paymentRepo.Get(paymentID)
	if err != nil {
		return PaymentPage{}, PaymentNotFound{paymentID}
	}
	token := p.GetGateway(bank).GetToken(payment)
	url := p.GetGateway(bank).GetPaymentPage(token)
	return url, nil
}

func (p PaymentService) Callback(data map[string][]string, bankName string) (CallbackResult, error) {
	bank, err := validateAndGetBank(bankName)
	if err != nil {
		return CallbackResult{Successful: false}, InvalidBankName{bankName}
	}
	verifyPayment, err := p.GetGateway(bank).VerifyPayment(data)
	if err != nil {
		return CallbackResult{Successful: false}, nil
	}

	payment, err := p.paymentRepo.Get(int(verifyPayment.ID))
	if err != nil {
		return CallbackResult{Successful: false}, nil
	}

	if verifyPayment.Amount != payment.Amount {
		return CallbackResult{Successful: false}, nil
	}

	payment.Status = domain.PAID
	p.paymentRepo.Update(payment)
	return CallbackResult{Successful: true, PaymentID: int(payment.ID)}, nil
}
