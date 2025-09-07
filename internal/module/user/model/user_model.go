package model

import "time"

// Entityという名前だと、すべてのフィールドをValue Objectにするイメージがある
// そこまで厳密でないものにしたいため、modelにする

type User struct {
	Id        int
	Name      string
	CreatedAt time.Time // 適当な型
	UpdatedAt time.Time // 適当な型
}
