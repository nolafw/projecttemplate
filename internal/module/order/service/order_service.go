package service

type OrderService interface {
	// FIXME: 実際には引数や戻り値はオブジェクトになるので、そこをちゃんと実装して試してみる
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
