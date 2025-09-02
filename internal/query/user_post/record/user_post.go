package record

// CQRSのQuery(composite_repository)の戻り値は、`***Record`という名前にして
// module内の、modelとは区別したものとする。
// Queryの戻り値は、ドメインロジックを含まない、データ専用のオブジェクトとして扱う。
type UserPostRecord struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	PostID    int64  `json:"post_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
