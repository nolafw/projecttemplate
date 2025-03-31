package userpost

// CQRSのQuery(composite_repository)の戻り値は、`***Model`という名前にして
// module内の、Entityとは区別したものとする。
type UserPostModel struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	PostID    int64  `json:"post_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
