package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"

	"gitlab.com/4uvirik/my-platform/services/order-service/internal/service"
	orderpb "gitlab.com/4uvirik/my-platform/services/order-service/proto"
)

type OrderHandler struct {
	orderpb.UnimplementedOrderServiceServer
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func Register(server *grpc.Server, h *OrderHandler) {
	orderpb.RegisterOrderServiceServer(server, h)
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	idempotencyKey := ""
	if v := md.Get("idempotency-key"); len(v) > 0 {
		idempotencyKey = v[0]
	}

	if idempotencyKey == "" {
		return nil, status.Error(codes.InvalidArgument, "missing idempotency-key")
	}

	id, err := h.svc.Create(ctx, req.UserId, req.Amount, idempotencyKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &orderpb.CreateOrderResponse{OrderId: id.String()}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	id, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, err
	}

	order, err := h.svc.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &orderpb.GetOrderResponse{
		Order: &orderpb.Order{
			Id:     order.ID.String(),
			UserId: order.UserID,
			Amount: order.Amount,
			Status: order.Status,
		},
	}, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	orders, err := h.svc.ListByUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var pbOrders []*orderpb.Order
	for _, o := range orders {
		pbOrders = append(pbOrders, &orderpb.Order{
			Id:     o.ID.String(),
			UserId: o.UserID,
			Amount: o.Amount,
			Status: o.Status,
		})
	}

	return &orderpb.ListOrdersResponse{Orders: pbOrders}, nil
}
