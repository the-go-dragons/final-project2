package usecase

import (
	"time"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type WalletService struct {
	walletRepo  persistence.WalletRepository
	paymentRepo persistence.PaymentRepository
	trxRepo     persistence.TransactionRepository
}

func NewWallet(
	walletRepo persistence.WalletRepository,
	paymentRepo persistence.PaymentRepository,
	trxRepo persistence.TransactionRepository,
) WalletService {
	return WalletService{
		walletRepo:  walletRepo,
		paymentRepo: paymentRepo,
		trxRepo:     trxRepo,
	}
}
func (w WalletService) ChargeRequest(walletId int, amount int64) (uint, error) {
	_, err := w.walletRepo.Get(walletId)
	if err != nil {
		return 0, WallertNotFound{walletId}
	}
	payment := domain.Payment{Amount: amount, WalletID: walletId}
	payment, err = w.paymentRepo.Create(payment)
	if err != nil {
		return 0, err
	}
	return payment.ID, nil
}

func (w WalletService) FinalizeCharge(paymentID int) (uint, error) {
	payment, err := w.paymentRepo.Get(paymentID)
	if err != nil {
		return 0, PaymentNotFound{paymentID}
	}
	if payment.Status == domain.UNPAID {
		return 0, PaymentNotPaid{paymentID}
	}
	if payment.Status == domain.APPLIED {
		return 0, PaymentAlreadyApplied{paymentID}
	}
	if payment.Status != domain.PAID {
		return 0, InvalidPaymentStatus{paymentID, payment.Status}
	}
	err = w.walletRepo.ChargeWallet(payment.WalletID, payment.Amount)
	if err != nil {
		return 0, err
	}
	transaction := domain.Transaction{
		Amount:          payment.Amount,
		WalletID:        payment.WalletID,
		Status:          domain.DEPOSIT,
		TransactionDate: time.Now(),
	}
	_, err = w.trxRepo.Create(transaction)
	if err != nil {
		return 0, err
	}
	payment.Status = domain.APPLIED
	_, err = w.paymentRepo.Update(payment)
	if err != nil {
		return 0, err
	}
	return uint(payment.WalletID), nil
}
