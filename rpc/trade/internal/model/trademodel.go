package model

import (
	"sync"
	"time"
)

type Trade struct {
	TradeId    int64
	UserId     int64
	ProductId  int64
	Quantity   int64
	TotalAmount float64
	Status     string
	CreatedAt  int64
}

type TradeModel struct {
	mu     sync.RWMutex
	trades map[int64]*Trade
	nextId int64
}

func NewTradeModel() *TradeModel {
	return &TradeModel{
		trades: make(map[int64]*Trade),
		nextId: 1,
	}
}

func (m *TradeModel) CreateTrade(userId, productId, quantity int64, totalAmount float64) *Trade {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	trade := &Trade{
		TradeId:     m.nextId,
		UserId:      userId,
		ProductId:   productId,
		Quantity:    quantity,
		TotalAmount: totalAmount,
		Status:      "success",
		CreatedAt:   time.Now().Unix(),
	}
	m.trades[m.nextId] = trade
	m.nextId++
	return trade
}

