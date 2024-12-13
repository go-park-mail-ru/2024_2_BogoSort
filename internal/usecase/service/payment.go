package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/postgres"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type PaymentService struct {
	paymentSecret string
	paymentShopID string

	paymentRepo repository.PaymentRepository
	advertRepo  repository.AdvertRepository
}

func NewPaymentService(paymentShopID, paymentSecret string,
	paymentRepo repository.PaymentRepository, advertRepo repository.AdvertRepository) *PaymentService {
	return &PaymentService{
		paymentSecret: paymentSecret,
		paymentShopID: paymentShopID,
		paymentRepo:   paymentRepo,
		advertRepo:    advertRepo,
	}
}

func (s *PaymentService) PaymentProcessor(ctx context.Context) {
	process := func(order entity.Order) error {
		url := fmt.Sprintf("https://api.yookassa.ru/v3/payments/%s", order.PaymentID)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request for order %s: %w", order.OrderID, err)
		}

		req.SetBasicAuth(s.paymentShopID, s.paymentSecret)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request for order %s: %w", order.OrderID, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code %d for order %s", resp.StatusCode, order.OrderID)
		}

		var paymentResponse struct {
			Status   string `json:"status"`
			Metadata struct {
				ItemID uuid.UUID `json:"item_id"`
			} `json:"metadata"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&paymentResponse); err != nil {
			return fmt.Errorf("failed to decode response for order %s: %w", order.OrderID, err)
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
			return fmt.Errorf("failed to update status for order %s: %w", order.OrderID, err)
		}

		if status == postgres.OrderStatusCompleted {
			_, err := s.advertRepo.PromoteAdvert(paymentResponse.Metadata.ItemID)
			if err != nil {
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
					log.Println("Failed to fetch orders:", err)
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
	process func(order entity.Order) error) {
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
						log.Println("Failed to process order:", err)
					}
				}
			}
		}()
	}

	wg.Wait()
}

const promotionAmount = "123.00"
const returnURL = "http://5.188.141.136:8008"

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
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", uuid.NewString())
	req.SetBasicAuth(s.paymentShopID, s.paymentSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	_, err = s.paymentRepo.InsertOrder(paymentReq.Metadata.OrderID,
		paymentReq.Amount.Value,
		paymentResp.ID,
		postgres.OrderStatusInProcess)
	if err != nil {
		return nil, err
	}

	return &paymentResp.Confirmation.ConfirmationURL, nil
}