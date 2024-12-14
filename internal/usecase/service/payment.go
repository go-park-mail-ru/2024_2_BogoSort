package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/postgres"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PaymentService struct {
	paymentSecret string
	paymentShopID string

	paymentRepo repository.PaymentRepository
	advertRepo  repository.AdvertRepository
}

func NewPaymentService(paymentShopID, paymentSecret string,
	paymentRepo repository.PaymentRepository, advertRepo repository.AdvertRepository,
) *PaymentService {
	return &PaymentService{
		paymentSecret: paymentSecret,
		paymentShopID: paymentShopID,
		paymentRepo:   paymentRepo,
		advertRepo:    advertRepo,
	}
}

func (s *PaymentService) PaymentProcessor(ctx context.Context) {
	process := func(order entity.Order) error {
		logger := middleware.GetLogger(ctx)

		url := "https://api.yookassa.ru/v3/payments/" + order.PaymentID

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logger.Error("failed to create request", zap.String("order_id", order.OrderID), zap.Error(err))
			return err
		}

		req.SetBasicAuth(s.paymentShopID, s.paymentSecret)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error("failed to send request", zap.String("order_id", order.OrderID), zap.Error(err))
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logger.Error("unexpected status code", zap.Int("status_code", resp.StatusCode), zap.String("order_id", order.OrderID))
			return err
		}

		var paymentResponse struct {
			Status   string `json:"status"`
			Metadata struct {
				ItemID uuid.UUID `json:"item_id"`
			} `json:"metadata"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&paymentResponse); err != nil {
			logger.Error("failed to decode response", zap.String("order_id", order.OrderID), zap.Error(err))
			return err
		}

		if paymentResponse.Status != "succeeded" && paymentResponse.Status != "canceled" {
			return nil
		}

		var status string
		switch paymentResponse.Status {
		case "canceled":
			status = postgres.OrderStatusCanceled
		case "succeeded":
			status = postgres.OrderStatusCompleted
		}

		_, err = s.paymentRepo.UpdateOrderStatus(order.OrderID, status)
		if err != nil {
			logger.Error("failed to update status", zap.String("order_id", order.OrderID), zap.Error(err))
			return err
		}

		if status == postgres.OrderStatusCompleted {
			_, err := s.advertRepo.PromoteAdvert(paymentResponse.Metadata.ItemID)
			if err != nil {
				logger.Error("failed to promote advert", zap.String("item_id", paymentResponse.Metadata.ItemID.String()), zap.Error(err))
				return err
			}
		}

		return nil
	}

	orders := make(chan entity.Order, 25)

	go func() {
		defer close(orders)
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				unprocessedOrders, err := s.paymentRepo.GetOrdersInProcess()
				if err != nil {
					logger := middleware.GetLogger(ctx)
					logger.Error("failed to fetch orders", zap.Error(err))
					continue
				}

				for _, order := range unprocessedOrders {
					select {
					case orders <- order:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	workerPool(ctx, orders, process)
}

func workerPool(ctx context.Context, orders <-chan entity.Order,
	process func(order entity.Order) error,
) {
	wg := new(sync.WaitGroup)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case order, ok := <-orders:
					if !ok {
						return
					}
					err := process(order)
					if err != nil {
						logger := middleware.GetLogger(ctx)
						logger.Error("failed to process order", zap.Error(err))
					}
				}
			}
		}()
	}

	wg.Wait()
}

const (
	promotionAmount = "50.00"
	returnURL       = "http://5.188.141.136:8008"
)

type PaymentRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Capture      bool `json:"capture"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnURL string `json:"return_url"`
	} `json:"confirmation"`
	Description string `json:"description"`
	Metadata    struct {
		OrderID string `json:"order_id"`
		ItemID  string `json:"item_id"`
	} `json:"metadata"`
}

type PaymentResponse struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Confirmation struct {
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
}

func (s *PaymentService) InitPayment(itemId string) (*string, error) {
	url := "https://api.yookassa.ru/v3/payments"

	paymentReq := PaymentRequest{
		Capture:     true,
		Description: "Платное продвижение товара",
	}
	paymentReq.Amount.Value = promotionAmount
	paymentReq.Amount.Currency = "RUB"
	paymentReq.Confirmation.Type = "redirect"
	paymentReq.Confirmation.ReturnURL = returnURL
	paymentReq.Metadata.OrderID = uuid.NewString()
	paymentReq.Metadata.ItemID = itemId

	requestBody, err := json.Marshal(paymentReq)
	if err != nil {
		logger := middleware.GetLogger(context.Background())
		logger.Error("failed to marshal request body", zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		logger := middleware.GetLogger(context.Background())
		logger.Error("failed to create request", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", uuid.NewString())
	req.SetBasicAuth(s.paymentShopID, s.paymentSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger := middleware.GetLogger(context.Background())
		logger.Error("request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger := middleware.GetLogger(context.Background())
		logger.Error("failed to read response body", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		logger := middleware.GetLogger(context.Background())
		logger.Error("unexpected status code", zap.Int("status_code", resp.StatusCode), zap.String("response", string(body)))
		return nil, err
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		logger := middleware.GetLogger(context.Background())
		logger.Error("failed to unmarshal response", zap.Error(err))
		return nil, err
	}

	_, err = s.paymentRepo.InsertOrder(paymentReq.Metadata.OrderID,
		paymentReq.Amount.Value,
		paymentResp.ID,
		postgres.OrderStatusInProcess)
	if err != nil {
		logger := middleware.GetLogger(context.Background())
		logger.Error("failed to insert order", zap.Error(err))
		return nil, err
	}

	return &paymentResp.Confirmation.ConfirmationURL, nil
}
