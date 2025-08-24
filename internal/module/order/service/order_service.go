package service

type OrderService interface {
	GetOrder() (string, error)
}

func NewOrderService() OrderService {
	return &OrderServiceImpl{}
}

type OrderServiceImpl struct {
}

func (s *OrderServiceImpl) GetOrder() (string, error) {
	return "order", nil
}
