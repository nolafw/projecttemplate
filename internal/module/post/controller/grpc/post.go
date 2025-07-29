package grpc

import (
	"context"

	"github.com/nolafw/projecttemplate/internal/module/post/service"
	pb "github.com/nolafw/projecttemplate/service_adapter/post"
)

type PostGRPCService struct {
	pb.UnimplementedPostServer
	service service.PostService
}

func NewPostGRPCService(service service.PostService) *PostGRPCService {
	return &PostGRPCService{
		service: service,
	}
}

func (s *PostGRPCService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	return &pb.GetPostResponse{
		PostId:  "1",
		Content: "Sample Post Content",
	}, nil
}
