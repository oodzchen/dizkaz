package service

import "github.com/oodzchen/dizkaz/store"

// Category 服务层结构
type Category struct {
	Store *store.Store
}

// NewCategory 创建新的 Category 服务
func NewCategory(store *store.Store) *Category {
	return &Category{
		Store: store,
	}
}