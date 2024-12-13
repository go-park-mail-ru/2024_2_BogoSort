package cart_purchase

import (
	"context"
	"fmt"

	cartPurchaseProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrInvalidPurchaseStatus = errors.New("invalid purchase status")
	ErrInvalidPaymentMethod  = errors.New("invalid payment method")
	ErrInvalidDeliveryMethod = errors.New("invalid delivery method")
	ErrPurchaseNotFound      = errors.New("purchase not found")
	ErrCartNotFound          = errors.New("cart not found")
)

type CartPurchaseClient struct {
	client cartPurchaseProto.CartPurchaseServiceClient
	conn   *grpc.ClientConn
}

func NewCartPurchaseClient(addr string) (*CartPurchaseClient, error) {
	//nolint:staticcheck // Suppressing deprecation warning for grpc.Dial
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := cartPurchaseProto.NewCartPurchaseServiceClient(conn)

	_, err = client.Ping(context.Background(), &cartPurchaseProto.NoContent{})
	if err != nil {
		return nil, err
	}

	return &CartPurchaseClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *CartPurchaseClient) Close() error {
	return c.conn.Close()
}

func (c *CartPurchaseClient) AddPurchase(ctx context.Context, req dto.PurchaseRequest) (*dto.PurchaseResponse, error) {
	paymentMethod, err := ConvertDBPaymentMethodToEnum(string(req.PaymentMethod))
	if err != nil {
		return nil, err
	}

	deliveryMethod, err := ConvertDBDeliveryMethodToEnum(string(req.DeliveryMethod))
	if err != nil {
		return nil, err
	}

	protoReq := &cartPurchaseProto.AddPurchaseRequest{
		CartId:         req.CartID.String(),
		Address:        req.Address,
		PaymentMethod:  paymentMethod,
		DeliveryMethod: deliveryMethod,
		UserId:         req.UserID.String(),
	}

	resp, err := c.client.AddPurchase(ctx, protoReq)
	if err != nil {
		return nil, errors.Wrap(ErrPurchaseNotFound, err.Error())
	}

	dtoAdverts := make([]dto.PreviewAdvertCard, 0)
	for _, a := range resp.Adverts {
		dtoAdverts = append(dtoAdverts, convertPreviewAdvertCardFromProto(a))
	}

	fmt.Println("resp.Status", resp.Status)
	fmt.Println("resp.PaymentMethod", resp.PaymentMethod)
	fmt.Println("resp.DeliveryMethod", resp.DeliveryMethod)

	fmt.Println("resp.Status", ConvertPurchaseStatusToDB(resp.Status))
	fmt.Println("resp.PaymentMethod", ConvertPaymentMethodToDB(resp.PaymentMethod))
	fmt.Println("resp.DeliveryMethod", ConvertDeliveryMethodToDB(resp.DeliveryMethod))

	response := &dto.PurchaseResponse{
		ID:             uuid.MustParse(resp.Id),
		SellerID:       uuid.MustParse(resp.SellerId),
		CustomerID:     uuid.MustParse(resp.CustomerId),
		Adverts:        dtoAdverts,
		Address:        resp.Address,
		Status:         dto.PurchaseStatus(ConvertPurchaseStatusToDB(resp.Status)),
		PaymentMethod:  dto.PaymentMethod(ConvertPaymentMethodToDB(resp.PaymentMethod)),
		DeliveryMethod: dto.DeliveryMethod(ConvertDeliveryMethodToDB(resp.DeliveryMethod)),
	}

	return response, nil
}

func (c *CartPurchaseClient) GetPurchasesByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.PurchaseResponse, error) {
	protoReq := &cartPurchaseProto.GetPurchasesByUserIDRequest{
		UserId: userID.String(),
	}

	resp, err := c.client.GetPurchasesByUserID(ctx, protoReq)
	if err != nil {
		return nil, errors.Wrap(ErrPurchaseNotFound, "purchases not found")
	}

	var purchases []*dto.PurchaseResponse
	for _, p := range resp.Purchases {
		dtoAdverts := make([]dto.PreviewAdvertCard, 0)
		for _, a := range p.Adverts {
			dtoAdverts = append(dtoAdverts, convertPreviewAdvertCardFromProto(a))
		}

		response := &dto.PurchaseResponse{
			ID:             uuid.MustParse(p.Id),
			SellerID:       uuid.MustParse(p.SellerId),
			CustomerID:     uuid.MustParse(p.CustomerId),
			Address:        p.Address,
			Adverts:        dtoAdverts,
			Status:         dto.PurchaseStatus(ConvertPurchaseStatusToDB(p.Status)),
			PaymentMethod:  dto.PaymentMethod(ConvertPaymentMethodToDB(p.PaymentMethod)),
			DeliveryMethod: dto.DeliveryMethod(ConvertDeliveryMethodToDB(p.DeliveryMethod)),
		}

		fmt.Println("RESPONSE", response)

		purchases = append(purchases, response)
	}

	return purchases, nil
}

func (c *CartPurchaseClient) GetCartByID(ctx context.Context, cartID uuid.UUID) (*dto.Cart, error) {
	protoReq := &cartPurchaseProto.GetCartByIDRequest{
		CartId: cartID.String(),
	}

	resp, err := c.client.GetCartByID(ctx, protoReq)
	if err != nil {
		return nil, errors.Wrap(ErrCartNotFound, "cart not found")
	}

	return convertCartFromProto(resp.Cart), nil
}

func (c *CartPurchaseClient) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*dto.Cart, error) {
	protoReq := &cartPurchaseProto.GetCartByUserIDRequest{
		UserId: userID.String(),
	}

	resp, err := c.client.GetCartByUserID(ctx, protoReq)
	if err != nil {
		return nil, errors.Wrap(ErrCartNotFound, "cart not found")
	}

	return convertCartFromProto(resp.Cart), nil
}

func (c *CartPurchaseClient) AddAdvertToCart(ctx context.Context, userID uuid.UUID, advertID uuid.UUID) (string, error) {
	protoReq := &cartPurchaseProto.AddAdvertToCartRequest{
		UserId:   userID.String(),
		AdvertId: advertID.String(),
	}

	resp, err := c.client.AddAdvertToCart(ctx, protoReq)
	if err != nil {
		return "", errors.Wrap(ErrCartNotFound, "cart not found")
	}

	return resp.Message, nil
}

func (c *CartPurchaseClient) DeleteAdvertFromCart(ctx context.Context, cartID uuid.UUID, advertID uuid.UUID) (string, error) {
	protoReq := &cartPurchaseProto.DeleteAdvertFromCartRequest{
		CartId:   cartID.String(),
		AdvertId: advertID.String(),
	}

	resp, err := c.client.DeleteAdvertFromCart(ctx, protoReq)
	if err != nil {
		return "", err
	}

	return resp.Message, nil
}

func (c *CartPurchaseClient) CheckCartExists(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	protoReq := &cartPurchaseProto.CheckCartExistsRequest{
		UserId: userID.String(),
	}

	resp, err := c.client.CheckCartExists(ctx, protoReq)
	if err != nil {
		return uuid.Nil, errors.Wrap(ErrCartNotFound, "cart not found")
	}

	return uuid.MustParse(resp.CartId), nil
}

func (c *CartPurchaseClient) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &cartPurchaseProto.NoContent{})
	return err
}

func convertPreviewAdvertCardToProto(advert dto.PreviewAdvertCard) *cartPurchaseProto.PreviewAdvertCard {
	return &cartPurchaseProto.PreviewAdvertCard{
		Preview: &cartPurchaseProto.PreviewAdvert{
			Id:          advert.Preview.ID.String(),
			Title:       advert.Preview.Title,
			Price:       uint64(advert.Preview.Price),
			CategoryId:  advert.Preview.CategoryId.String(),
			ImageId:     advert.Preview.ImageId.String(),
			SellerId:    advert.Preview.SellerId.String(),
			Status:      cartPurchaseProto.AdvertStatus(cartPurchaseProto.AdvertStatus_value[string(advert.Preview.Status)]),
			Location:    advert.Preview.Location,
			HasDelivery: advert.Preview.HasDelivery,
		},
		IsSaved:  advert.IsSaved,
		IsViewed: advert.IsViewed,
	}
}

func convertPreviewAdvertCardFromProto(advert *cartPurchaseProto.PreviewAdvertCard) dto.PreviewAdvertCard {
	return dto.PreviewAdvertCard{
		Preview: dto.PreviewAdvert{
			ID:          uuid.MustParse(advert.Preview.Id),
			Title:       advert.Preview.Title,
			Price:       uint(advert.Preview.Price),
			ImageId:     uuid.MustParse(advert.Preview.ImageId),
			CategoryId:  uuid.MustParse(advert.Preview.CategoryId),
			SellerId:    uuid.MustParse(advert.Preview.SellerId),
			Status:      dto.AdvertStatus(ConvertAdvertStatusToDB(advert.Preview.Status)),
			Location:    advert.Preview.Location,
			HasDelivery: advert.Preview.HasDelivery,
		},
		IsSaved:  advert.IsSaved,
		IsViewed: advert.IsViewed,
	}
}

func convertCartFromProto(protoCart *cartPurchaseProto.Cart) *dto.Cart {
	cart := &dto.Cart{
		ID:            uuid.MustParse(protoCart.Id),
		UserID:        uuid.MustParse(protoCart.UserId),
		Status:        entity.CartStatus(ConvertCartStatusToDB(protoCart.Status)),
		CartPurchases: make([]dto.CartPurchase, 0),
	}

	for _, purchase := range protoCart.CartPurchases {
		cartPurchase := dto.CartPurchase{
			SellerID: uuid.MustParse(purchase.SellerId),
			Adverts:  make([]dto.PreviewAdvertCard, 0),
		}

		for _, advert := range purchase.Adverts {
			cartPurchase.Adverts = append(cartPurchase.Adverts, convertPreviewAdvertCardFromProto(advert))
		}
		cart.CartPurchases = append(cart.CartPurchases, cartPurchase)
	}

	return cart
}
