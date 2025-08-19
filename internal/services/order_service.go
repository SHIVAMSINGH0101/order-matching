
package services

import (
	orderModel "github.com/SHIVAMSINGH0101/go-demo/internal/models"
	 "github.com/SHIVAMSINGH0101/go-demo/internal/repository"
)
// This is OrderService layer
// All the business logic related to Order are performed
type OrderService interface {
	CreateLocation(loc *orderModel.Location) (int64, error)
	GetLocationByID(id int64) (*orderModel.Location, error)
	GetLocationsByIDs(ids []int64) ([]orderModel.Location, error)

	CreateOrder(order *orderModel.Order) (int64, error)
	GetOrderByID(orderId int64) (*orderModel.Order, error)
	GetOrdersByIDs(ids []int64) ([]orderModel.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(r repository.OrderRepository) OrderService {
	return &orderService{
		repo: r,
	}
}

func (s *orderService) CreateLocation(loc *orderModel.Location) (int64, error) {
	return s.repo.InsertLocation(loc)
}

func (s *orderService) GetLocationByID(id int64) (*orderModel.Location, error) {
	return s.repo.GetLocationByID(id)
}

func (s *orderService) GetLocationsByIDs(ids []int64) ([]orderModel.Location, error) {
	return s.repo.GetLocationsByIDs(ids)
}

func (s *orderService) CreateOrder(order *orderModel.Order) (int64, error) {
	return s.repo.InsertOrder(order)
}

func (s *orderService) GetOrderByID(orderId int64) (*orderModel.Order, error) {
	return s.repo.GetOrderByID(orderId)
}

func (s *orderService) GetOrdersByIDs(ids []int64) ([]orderModel.Order, error) {
	return s.repo.GetOrdersByIDs(ids)
}

