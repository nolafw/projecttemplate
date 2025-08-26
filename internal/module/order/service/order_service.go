package service

import (
	"fmt"

	"github.com/nolafw/projecttemplate/internal/module/order/dto"
)

type OrderService interface {
	// FIXME: 実際には引数や戻り値はオブジェクトになるので、そこをちゃんと実装して試してみる
	GetOrder() (*dto.Order, error)
}

func NewOrderService() OrderService {
	return &OrderServiceImpl{}
}

type OrderServiceImpl struct {
}

func (s *OrderServiceImpl) GetOrder() (*dto.Order, error) {
	fmt.Println("OrderServiceImpl GetOrder called")
	return &dto.Order{Id: 1, Amount: 100.0, Status: "completed"}, nil
}
