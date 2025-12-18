package model

import (
	"sync"
)

type Product struct {
	ProductId   int64
	Name        string
	Description string
	Price       float64
	Stock       int64
}

type ProductModel struct {
	mu       sync.RWMutex
	products map[int64]*Product
}

func NewProductModel() *ProductModel {
	pm := &ProductModel{
		products: make(map[int64]*Product),
	}
	// 初始化一些测试数据
	pm.products[1] = &Product{
		ProductId:   1,
		Name:        "iPhone 15",
		Description: "Apple iPhone 15 Pro Max",
		Price:       9999.00,
		Stock:       100,
	}
	pm.products[2] = &Product{
		ProductId:   2,
		Name:        "MacBook Pro",
		Description: "Apple MacBook Pro 16 inch",
		Price:       19999.00,
		Stock:       50,
	}
	return pm
}

func (m *ProductModel) GetProduct(productId int64) (*Product, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	product, ok := m.products[productId]
	return product, ok
}

func (m *ProductModel) ReduceStock(productId int64, quantity int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	product, ok := m.products[productId]
	if !ok || product.Stock < quantity {
		return false
	}
	product.Stock -= quantity
	return true
}
