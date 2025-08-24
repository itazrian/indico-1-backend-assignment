package services

import (
    "backend-assignment/module/models"
    "backend-assignment/module/repositories"
)

type OrderService struct {
    orderRepo *repositories.OrderRepo
    prodRepo  *repositories.ProductRepo
}

func NewOrderService() *OrderService {
    return &OrderService{
        orderRepo: repositories.NewOrderRepo(),
        prodRepo:  repositories.NewProductRepo(),
    }
}

func (s *OrderService) CreateOrder(o *models.Order) error {
    if err := s.prodRepo.ReduceStock(o.ProductID, o.Quantity); err != nil {
        return err
    }
    return s.orderRepo.Create(o)
}

func (s *OrderService) GetOrder(id int64) (*models.Order, error) {
    return s.orderRepo.GetByID(id)
}
