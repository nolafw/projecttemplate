package service

type PostService interface {
	Anything() string
}

func NewPostService() PostService {
	return &PostServiceImpl{}
}

type PostServiceImpl struct {
}

func (s *PostServiceImpl) Anything() string {
	return "post service response"
}
