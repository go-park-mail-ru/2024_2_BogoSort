package cart_purchase

import (
	"context"
	"errors"

	"github.com/google/uuid"

	proto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	proto.UnimplementedCartPurchaseServiceServer
	cartUC     usecase.Cart
	purchaseUC usecase.Purchase
}

func NewGrpcServer(cartUC usecase.Cart, purchaseUC usecase.Purchase) *GrpcServer {
	return &GrpcServer{
		cartUC:     cartUC,
		purchaseUC: purchaseUC,
	}
}

func (s *GrpcServer) AddPurchase(ctx context.Context, req *proto.AddPurchaseRequest) (*proto.AddPurchaseResponse, error) {
	paymentMethod := ConvertPaymentMethodToDB(req.PaymentMethod)
	deliveryMethod := ConvertDeliveryMethodToDB(req.DeliveryMethod)

	purchaseReq := dto.PurchaseRequest{
		CartID:         uuid.MustParse(req.CartId),
		Address:        req.Address,
		PaymentMethod:  dto.PaymentMethod(paymentMethod),
		DeliveryMethod: dto.DeliveryMethod(deliveryMethod),
		UserID:         uuid.MustParse(req.UserId),
	}

	purchase, err := s.purchaseUC.Add(purchaseReq, purchaseReq.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add purchase: %v", err)
	}

	protoAdverts := make([]*proto.PreviewAdvertCard, 0, len(purchase[0].Adverts))
	for _, a := range purchase[0].Adverts {
		protoAdvert := convertPreviewAdvertCardToProto(a)
		protoAdverts = append(protoAdverts, protoAdvert)
	}

	purchaseStatus, err := ConvertDBPurchaseStatusToEnum(string(purchase[0].Status))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert purchase status: %v", err)
	}
	purchasePaymentMethod, err := ConvertDBPaymentMethodToEnum(string(purchase[0].PaymentMethod))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert purchase payment method: %v", err)
	}
	purchaseDeliveryMethod, err := ConvertDBDeliveryMethodToEnum(string(purchase[0].DeliveryMethod))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert purchase delivery method: %v", err)
	}

	return &proto.AddPurchaseResponse{
		Id:             purchase[0].ID.String(),
		SellerId:       purchase[0].SellerID.String(),
		CustomerId:     purchase[0].CustomerID.String(),
		Address:        purchase[0].Address,
		Adverts:        protoAdverts,
		Status:         purchaseStatus,
		PaymentMethod:  purchasePaymentMethod,
		DeliveryMethod: purchaseDeliveryMethod,
	}, nil
}

func (s *GrpcServer) GetPurchasesByUserID(ctx context.Context, req *proto.GetPurchasesByUserIDRequest) (*proto.GetPurchasesByUserIDResponse, error) {
	userID := uuid.MustParse(req.UserId)
	purchases, err := s.purchaseUC.GetByUserId(userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get purchases: %v", err)
	}

	var protoPurchases []*proto.PurchaseResponse
	for _, p := range purchases {
		protoAdverts := make([]*proto.PreviewAdvertCard, 0, len(p.Adverts))
		for _, a := range p.Adverts {
			protoAdvert := &proto.PreviewAdvertCard{
				Preview: &proto.PreviewAdvert{
					Id:          a.Preview.ID.String(),
					Title:       a.Preview.Title,
					Price:       uint64(a.Preview.Price),
					ImageId:     a.Preview.ImageId.String(),
					CategoryId:  a.Preview.CategoryId.String(),
					SellerId:    a.Preview.SellerId.String(),
					Status:      proto.AdvertStatus(proto.AdvertStatus_value[string(a.Preview.Status)]),
					Location:    a.Preview.Location,
					HasDelivery: a.Preview.HasDelivery,
				},
				IsSaved:  a.IsSaved,
				IsViewed: a.IsViewed,
			}
			protoAdverts = append(protoAdverts, protoAdvert)
		}

		protoPurchases = append(protoPurchases, &proto.PurchaseResponse{
			Id:             p.ID.String(),
			SellerId:       p.SellerID.String(),
			CustomerId:     p.CustomerID.String(),
			Adverts:        protoAdverts,
			Address:        p.Address,
			Status:         proto.PurchaseStatus(proto.PurchaseStatus_value[string(p.Status)]),
			PaymentMethod:  proto.PaymentMethod(proto.PaymentMethod_value[string(p.PaymentMethod)]),
			DeliveryMethod: proto.DeliveryMethod(proto.DeliveryMethod_value[string(p.DeliveryMethod)]),
		})
	}

	return &proto.GetPurchasesByUserIDResponse{
		Purchases: protoPurchases,
	}, nil
}

func (s *GrpcServer) GetCartByID(ctx context.Context, req *proto.GetCartByIDRequest) (*proto.GetCartByIDResponse, error) {
	cart, err := s.cartUC.GetById(uuid.MustParse(req.CartId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cart: %v", err)
	}

	protoCart := &proto.Cart{
		Id:            cart.ID.String(),
		UserId:        cart.UserID.String(),
		Status:        proto.CartStatus(proto.CartStatus_value[string(cart.Status)]),
		CartPurchases: []*proto.CartPurchase{},
	}

	for _, purchase := range cart.CartPurchases {
		protoPurchase := &proto.CartPurchase{
			SellerId: purchase.SellerID.String(),
			Adverts:  []*proto.PreviewAdvertCard{},
		}

		for _, advert := range purchase.Adverts {
			protoPurchase.Adverts = append(protoPurchase.Adverts, convertPreviewAdvertCardToProto(advert))
		}
		protoCart.CartPurchases = append(protoCart.CartPurchases, protoPurchase)
	}

	return &proto.GetCartByIDResponse{
		Cart: protoCart,
	}, nil
}

func (s *GrpcServer) GetCartByUserID(ctx context.Context, req *proto.GetCartByUserIDRequest) (*proto.GetCartByUserIDResponse, error) {
	cart, err := s.cartUC.GetByUserId(uuid.MustParse(req.UserId))
	if err != nil {
		if errors.Is(err, repository.ErrCartNotFound) {
			return nil, status.Errorf(codes.NotFound, "cart not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get cart: %v", err)
	}

	protoCart := &proto.Cart{
		Id:            cart.ID.String(),
		UserId:        cart.UserID.String(),
		Status:        proto.CartStatus(proto.CartStatus_value[string(cart.Status)]),
		CartPurchases: []*proto.CartPurchase{},
	}

	for _, purchase := range cart.CartPurchases {
		protoPurchase := &proto.CartPurchase{
			SellerId: purchase.SellerID.String(),
			Adverts:  []*proto.PreviewAdvertCard{},
		}

		for _, advert := range purchase.Adverts {
			protoPurchase.Adverts = append(protoPurchase.Adverts, convertPreviewAdvertCardToProto(advert))
		}
		protoCart.CartPurchases = append(protoCart.CartPurchases, protoPurchase)
	}

	return &proto.GetCartByUserIDResponse{
		Cart: protoCart,
	}, nil
}

func (s *GrpcServer) AddAdvertToCart(ctx context.Context, req *proto.AddAdvertToCartRequest) (*proto.AddAdvertToCartResponse, error) {
	err := s.cartUC.AddAdvert(uuid.MustParse(req.UserId), uuid.MustParse(req.AdvertId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add advert to cart: %v", err)
	}

	return &proto.AddAdvertToCartResponse{
		Message: "advert added to user cart",
	}, nil
}

func (s *GrpcServer) DeleteAdvertFromCart(ctx context.Context, req *proto.DeleteAdvertFromCartRequest) (*proto.DeleteAdvertFromCartResponse, error) {
	err := s.cartUC.DeleteAdvert(uuid.MustParse(req.CartId), uuid.MustParse(req.AdvertId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete advert from cart: %v", err)
	}

	return &proto.DeleteAdvertFromCartResponse{
		Message: "advert deleted from user cart",
	}, nil
}

func (s *GrpcServer) CheckCartExists(ctx context.Context, req *proto.CheckCartExistsRequest) (*proto.CheckCartExistsResponse, error) {
	cart, err := s.cartUC.CheckExists(uuid.MustParse(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check cart existence: %v", err)
	}

	return &proto.CheckCartExistsResponse{
		CartId: cart.String(),
	}, nil
}

func (s *GrpcServer) Ping(ctx context.Context, req *proto.NoContent) (*proto.NoContent, error) {
	return &proto.NoContent{}, nil
}
