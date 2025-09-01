package service

import (
	"context"
	"fmt"
	"time"

	"github.com/nolafw/projecttemplate/service/adapter/order"
	pbPost "github.com/nolafw/projecttemplate/service/adapter/post"
)

// サービスを別のモジュールから使う場合は、
// 直接このサービスを呼び出すのではなく、
// 一度ServiceAdapterを通して呼び出すこと
// serviceの返す値は必ずDTOにすること
// entityを返さないように実装すること
// entityはserviceの中で処理でのみ使う。
type UserService interface {
	Something() string
	GetPostContent(postId string) (string, error)
}

// gRPCクライアントが必要な場合は、クライアントの型を指定する
func NewUserService(postClient pbPost.PostClient, orderService order.OrderService) UserService {
	return &UserServiceImpl{
		postClient:   postClient,
		orderService: orderService,
	}
}

type UserServiceImpl struct {
	postClient   pbPost.PostClient
	orderService order.OrderService
}

func (s *UserServiceImpl) Something() string {
	fmt.Printf("injected orderService: %T\n", s.orderService) // DEBUG: injectされてない
	s.orderService.GetOrder()                                 // DEBUG:
	return "hoge"
}

func (s *UserServiceImpl) GetPostContent(postId string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := s.postClient.GetPost(ctx, &pbPost.GetPostRequest{
		PostId: postId,
	})
	if err != nil {
		return "", err
	}
	return p.Content, nil

}
