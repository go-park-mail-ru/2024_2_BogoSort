package cart_purchase

import (
	"context"

	"github.com/google/uuid"

	proto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/cart_purchase/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
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
	purchaseReq := dto.PurchaseRequest{
		CartID:         uuid.MustParse(req.CartId),
		Address:        req.Address,
		PaymentMethod:  dto.PaymentMethod(req.PaymentMethod.String()),
		DeliveryMethod: dto.DeliveryMethod(req.DeliveryMethod.String()),
	}

	userID := uuid.MustParse(req.UserId)
	purchaseResp, err := s.purchaseUC.Add(purchaseReq, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add purchase: %v", err)
	}

	return &proto.AddPurchaseResponse{
		Id:             purchaseResp.ID.String(),
		CartId:         purchaseResp.CartID.String(),
		Address:        purchaseResp.Address,
		Status:         proto.PurchaseStatus(proto.PurchaseStatus_value[string(purchaseResp.Status)]),
		PaymentMethod:  proto.PaymentMethod(proto.PaymentMethod_value[string(purchaseResp.PaymentMethod)]),
		DeliveryMethod: proto.DeliveryMethod(proto.DeliveryMethod_value[string(purchaseResp.DeliveryMethod)]),
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
		protoPurchases = append(protoPurchases, &proto.PurchaseResponse{
			Id:             p.ID.String(),
			CartId:         p.CartID.String(),
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
	exists, err := s.cartUC.CheckExists(uuid.MustParse(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check cart existence: %v", err)
	}

	return &proto.CheckCartExistsResponse{
		Exists: exists,
	}, nil
}

func (s *GrpcServer) GetCartByID(ctx context.Context, req *proto.GetCartByIDRequest) (*proto.GetCartByIDResponse, error) {
	cart, err := s.cartUC.GetById(uuid.MustParse(req.CartId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cart: %v", err)
	}

	protoCart := &proto.Cart{
		Id:      cart.ID.String(),
		UserId:  cart.UserID.String(),
		Status:  proto.CartStatus(proto.CartStatus_value[string(cart.Status)]),
		Adverts: []*proto.PreviewAdvertCard{},
	}

	for _, advert := range cart.Adverts {
		protoCart.Adverts = append(protoCart.Adverts, &proto.PreviewAdvertCard{
			Preview: &proto.PreviewAdvert{
				AdvertId:    advert.Preview.ID.String(),
				Title:       advert.Preview.Title,
				Price:       uint64(advert.Preview.Price),
				ImageUrl:    advert.Preview.ImageURL,
				Status:      proto.AdvertStatus(proto.AdvertStatus_value[string(advert.Preview.Status)]),
				Location:    advert.Preview.Location,
				HasDelivery: advert.Preview.HasDelivery,
			},
			IsSaved: advert.IsSaved,
			IsViewed: advert.IsViewed,
		})
	}

	return &proto.GetCartByIDResponse{
		Cart: protoCart,
	}, nil
}

func (s *GrpcServer) GetByUserID(ctx context.Context, req *proto.GetCartByUserIDRequest) (*proto.GetCartByUserIDResponse, error) {
	cart, err := s.cartUC.GetByUserId(uuid.MustParse(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cart: %v", err)
	}

	protoCart := &proto.Cart{
		Id:      cart.ID.String(),
		UserId:  cart.UserID.String(),
		Status:  proto.CartStatus(proto.CartStatus_value[string(cart.Status)]),
		Adverts: []*proto.PreviewAdvertCard{},
	}

	for _, advert := range cart.Adverts {
		protoCart.Adverts = append(protoCart.Adverts, &proto.PreviewAdvertCard{
			Preview: &proto.PreviewAdvert{
				AdvertId:    advert.Preview.ID.String(),
				Title:       advert.Preview.Title,
				Price:       uint64(advert.Preview.Price),
				ImageUrl:    advert.Preview.ImageURL,
				Status:      proto.AdvertStatus(proto.AdvertStatus_value[string(advert.Preview.Status)]),
				Location:    advert.Preview.Location,
				HasDelivery: advert.Preview.HasDelivery,
			},
			IsSaved: advert.IsSaved,
			IsViewed: advert.IsViewed,
		})
	}

	return &proto.GetCartByUserIDResponse{
		Cart: protoCart,
	}, nil
}

func (s *GrpcServer) Ping(ctx context.Context, req *proto.NoContent) (*proto.NoContent, error) {
	return &proto.NoContent{}, nil
}