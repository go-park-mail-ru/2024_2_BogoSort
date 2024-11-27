package cart_purchase

// import (
// 	"context"

// 	cartPurchaseProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase/proto"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/metrics"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/interceptors"
// 	"github.com/google/uuid"
// 	"github.com/pkg/errors"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// var (
// 	ErrInvalidPurchaseStatus = errors.New("invalid purchase status")
// 	ErrInvalidPaymentMethod  = errors.New("invalid payment method")
// 	ErrInvalidDeliveryMethod = errors.New("invalid delivery method")
// 	ErrPurchaseNotFound      = errors.New("purchase not found")
// 	ErrCartNotFound          = errors.New("cart not found")
// )

// type CartPurchaseClient struct {
// 	client cartPurchaseProto.CartPurchaseServiceClient
// 	conn   *grpc.ClientConn
// }

// func NewCartPurchaseClient(addr string) (*CartPurchaseClient, error) {
// 	metrics, err := metrics.NewGRPCMetrics("cart_purchase")
// 	if err != nil {
// 		return nil, err
// 	}
// 	metricsInterceptor := interceptors.NewMetricsInterceptor(*metrics)
// 	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(metricsInterceptor.ServeMetricsClientInterceptor))
// 	if err != nil {
// 		return nil, err
// 	}

// 	client := cartPurchaseProto.NewCartPurchaseServiceClient(conn)

// 	_, err = client.Ping(context.Background(), &cartPurchaseProto.NoContent{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &CartPurchaseClient{
// 		client: client,
// 		conn:   conn,
// 	}, nil
// }

// func (c *CartPurchaseClient) Close() error {
// 	return c.conn.Close()
// }

// func (c *CartPurchaseClient) AddPurchase(ctx context.Context, req dto.PurchaseRequest) (*dto.PurchaseResponse, error) {
// 	protoReq := &cartPurchaseProto.AddPurchaseRequest{
// 		CartId:         req.CartID.String(),
// 		Address:        req.Address,
// 		PaymentMethod:  cartPurchaseProto.PaymentMethod(cartPurchaseProto.PaymentMethod_value[string(req.PaymentMethod)]),
// 		DeliveryMethod: cartPurchaseProto.DeliveryMethod(cartPurchaseProto.DeliveryMethod_value[string(req.DeliveryMethod)]),
// 		UserId:         req.UserID.String(),
// 	}

// 	resp, err := c.client.AddPurchase(ctx, protoReq)
// 	if err != nil {
// 		return nil, errors.Wrap(ErrPurchaseNotFound, err.Error())
// 	}

// 	purchaseStatus:= ConvertPurchaseStatusToDB(resp.Status)
// 	paymentMethod := ConvertPaymentMethodToDB(resp.PaymentMethod)
// 	deliveryMethod := ConvertDeliveryMethodToDB(resp.DeliveryMethod)

// 	response := &dto.PurchaseResponse{
// 		ID:             uuid.MustParse(resp.Id),
// 		CartID:         uuid.MustParse(resp.CartId),
// 		Address:        resp.Address,
// 		Status:         dto.PurchaseStatus(purchaseStatus),
// 		PaymentMethod:  dto.PaymentMethod(paymentMethod),
// 		DeliveryMethod: dto.DeliveryMethod(deliveryMethod),
// 	}

// 	return response, nil
// }

// func (c *CartPurchaseClient) GetPurchasesByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.PurchaseResponse, error) {
// 	protoReq := &cartPurchaseProto.GetPurchasesByUserIDRequest{
// 		UserId: userID.String(),
// 	}

// 	resp, err := c.client.GetPurchasesByUserID(ctx, protoReq)
// 	if err != nil {
// 		return nil, errors.Wrap(ErrPurchaseNotFound, "purchases not found")
// 	}

// 	var purchases []*dto.PurchaseResponse
// 	for _, p := range resp.Purchases {
// 		purchaseStatus := ConvertPurchaseStatusToDB(p.Status)
// 		paymentMethod := ConvertPaymentMethodToDB(p.PaymentMethod)
// 		deliveryMethod := ConvertDeliveryMethodToDB(p.DeliveryMethod)

// 		purchases = append(purchases, &dto.PurchaseResponse{
// 			ID:             uuid.MustParse(p.Id),
// 			CartID:         uuid.MustParse(p.CartId),
// 			Address:        p.Address,
// 			Status:         dto.PurchaseStatus(purchaseStatus),
// 			PaymentMethod:  dto.PaymentMethod(paymentMethod),
// 			DeliveryMethod: dto.DeliveryMethod(deliveryMethod),
// 		})
// 	}

// 	return purchases, nil
// }

// func (c *CartPurchaseClient) GetCartByID(ctx context.Context, cartID uuid.UUID) (*dto.Cart, error) {
// 	protoReq := &cartPurchaseProto.GetCartByIDRequest{
// 		CartId: cartID.String(),
// 	}

// 	resp, err := c.client.GetCartByID(ctx, protoReq)
// 	if err != nil {
// 		return nil, errors.Wrap(ErrCartNotFound, "cart not found")
// 	}

// 	cart := &dto.Cart{
// 		ID:      uuid.MustParse(resp.Cart.Id),
// 		UserID:  uuid.MustParse(resp.Cart.UserId),
// 		Status:  entity.CartStatus(resp.Cart.Status),
// 		Adverts: []dto.PreviewAdvertCard{},
// 	}

// 	for _, advert := range resp.Cart.Adverts {
// 		ad := dto.PreviewAdvertCard{
// 			Preview: dto.PreviewAdvert{
// 				ID:          uuid.MustParse(advert.Preview.AdvertId),
// 				Title:       advert.Preview.Title,
// 				Price:       uint(advert.Preview.Price),
// 				ImageId:     uuid.MustParse(advert.Preview.ImageId),
// 				Status:      dto.AdvertStatus(advert.Preview.Status),
// 				Location:    advert.Preview.Location,
// 				HasDelivery: advert.Preview.HasDelivery,
// 			},
// 			IsSaved:  advert.IsSaved,
// 			IsViewed: advert.IsViewed,
// 		}
// 		cart.Adverts = append(cart.Adverts, ad)
// 	}

// 	return cart, nil
// }

// func (c *CartPurchaseClient) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*dto.Cart, error) {
// 	protoReq := &cartPurchaseProto.GetCartByUserIDRequest{
// 		UserId: userID.String(),
// 	}

// 	resp, err := c.client.GetCartByUserID(ctx, protoReq)
// 	if err != nil {
// 		return nil, errors.Wrap(ErrCartNotFound, "cart not found")
// 	}

// 	cart := &dto.Cart{
// 		ID:      uuid.MustParse(resp.Cart.Id),
// 		UserID:  uuid.MustParse(resp.Cart.UserId),
// 		Status:  entity.CartStatus(resp.Cart.Status),
// 		Adverts: []dto.PreviewAdvertCard{},
// 	}

// 	for _, advert := range resp.Cart.Adverts {
// 		ad := dto.PreviewAdvertCard{
// 			Preview: dto.PreviewAdvert{
// 				ID:          uuid.MustParse(advert.Preview.AdvertId),
// 				Title:       advert.Preview.Title,
// 				Price:       uint(advert.Preview.Price),
// 				ImageId:     uuid.MustParse(advert.Preview.ImageId),
// 				Status:      dto.AdvertStatus(advert.Preview.Status),
// 				Location:    advert.Preview.Location,
// 				HasDelivery: advert.Preview.HasDelivery,
// 			},
// 			IsSaved:  advert.IsSaved,
// 			IsViewed: advert.IsViewed,
// 		}
// 		cart.Adverts = append(cart.Adverts, ad)
// 	}

// 	return cart, nil
// }

// func (c *CartPurchaseClient) AddAdvertToCart(ctx context.Context, userID uuid.UUID, advertID uuid.UUID) (string, error) {
// 	protoReq := &cartPurchaseProto.AddAdvertToCartRequest{
// 		UserId:   userID.String(),
// 		AdvertId: advertID.String(),
// 	}

// 	resp, err := c.client.AddAdvertToCart(ctx, protoReq)
// 	if err != nil {
// 		return "", errors.Wrap(ErrCartNotFound, "cart not found")
// 	}

// 	return resp.Message, nil
// }

// func (c *CartPurchaseClient) DeleteAdvertFromCart(ctx context.Context, cartID uuid.UUID, advertID uuid.UUID) (string, error) {
// 	protoReq := &cartPurchaseProto.DeleteAdvertFromCartRequest{
// 		CartId:   cartID.String(),
// 		AdvertId: advertID.String(),
// 	}

// 	resp, err := c.client.DeleteAdvertFromCart(ctx, protoReq)
// 	if err != nil {
// 		return "", err
// 	}

// 	return resp.Message, nil
// }

// func (c *CartPurchaseClient) CheckCartExists(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
// 	protoReq := &cartPurchaseProto.CheckCartExistsRequest{
// 		UserId: userID.String(),
// 	}

// 	resp, err := c.client.CheckCartExists(ctx, protoReq)
// 	if err != nil {
// 		return uuid.Nil, errors.Wrap(ErrCartNotFound, "cart not found")
// 	}

// 	return uuid.MustParse(resp.CartId), nil
// }

// func (c *CartPurchaseClient) Ping(ctx context.Context) error {
// 	_, err := c.client.Ping(ctx, &cartPurchaseProto.NoContent{})
// 	return err
// }
