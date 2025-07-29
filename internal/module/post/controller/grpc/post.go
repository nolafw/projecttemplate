package grpc

import (
	"context"

	"github.com/nolafw/projecttemplate/internal/module/post/service"
	pb "github.com/nolafw/projecttemplate/service_adapter/post"
)

type PostAPI struct {
	pb.UnimplementedPostServer
	service service.PostService
}

func NewPostAPI(service service.PostService) *PostAPI {
	return &PostAPI{
		service: service,
	}
}

func (s *PostAPI) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	return &pb.GetPostResponse{
		PostId:  "1",
		Content: "Sample Post Content",
	}, nil
}
