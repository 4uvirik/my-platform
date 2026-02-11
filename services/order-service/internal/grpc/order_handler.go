package grpc

import (
	"context"
	"google.golang.org/grpc"

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
	id, err := h.svc.Create(ctx, req.UserId, req.Amount)
	if err != nil {
		return nil, err
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
