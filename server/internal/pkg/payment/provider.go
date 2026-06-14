package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"gosh/internal/model"
)

const mockSecret = "gosh-mock-secret-key"

type PaymentResult struct {
	TransactionNo string
	PayAmount     int64
	Status        string
	RawData       string
}

type CallbackResult struct {
	TransactionNo string
	OrderNo       string
	Amount        int64
	Status        string
	RawData       string
	SignOk        bool
}

type MockCallbackData struct {
	TransactionNo string `json:"transaction_no"`
	OrderNo       string `json:"order_no"`
	Amount        int64  `json:"amount"`
	Timestamp     int64  `json:"timestamp"`
	Sign          string `json:"sign"`
	Method        string `json:"method"`
}

type Provider interface {
	CreatePayment(order *model.Order, method string) (*PaymentResult, error)
	ProcessCallback(method string, notifyData []byte) (*CallbackResult, error)
}

type mockProvider struct{}

func NewMockProvider() Provider {
	return &mockProvider{}
}

func (p *mockProvider) CreatePayment(order *model.Order, method string) (*PaymentResult, error) {
	txNo := generateMockTxNo()
	now := time.Now()

	data := MockCallbackData{
		TransactionNo: txNo,
		OrderNo:       order.OrderNo,
		Amount:        order.PayAmount,
		Timestamp:     now.Unix(),
		Method:        method,
	}
	sign := signMockData(data)
	data.Sign = sign

	raw, _ := json.Marshal(data)

	return &PaymentResult{
		TransactionNo: txNo,
		PayAmount:     order.PayAmount,
		Status:        model.PaymentStatusSuccess,
		RawData:       string(raw),
	}, nil
}

func (p *mockProvider) ProcessCallback(method string, notifyData []byte) (*CallbackResult, error) {
	var data MockCallbackData
	if err := json.Unmarshal(notifyData, &data); err != nil {
		return nil, fmt.Errorf("invalid callback data: %w", err)
	}

	expectedSign := signMockData(data)
	signOk := hmac.Equal([]byte(data.Sign), []byte(expectedSign))

	return &CallbackResult{
		TransactionNo: data.TransactionNo,
		OrderNo:       data.OrderNo,
		Amount:        data.Amount,
		Status:        model.PaymentStatusSuccess,
		RawData:       string(notifyData),
		SignOk:        signOk,
	}, nil
}

func generateMockTxNo() string {
	return fmt.Sprintf("MOCK%s%04d", time.Now().Format("20060102150405"), time.Now().UnixMilli()%10000)
}

func signMockData(data MockCallbackData) string {
	raw := fmt.Sprintf("%s|%s|%d|%d|%s", data.TransactionNo, data.OrderNo, data.Amount, data.Timestamp, data.Method)
	mac := hmac.New(sha256.New, []byte(mockSecret))
	mac.Write([]byte(raw))
	return hex.EncodeToString(mac.Sum(nil))
}
