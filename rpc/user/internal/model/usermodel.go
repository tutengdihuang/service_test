package model

import (
	"sync"
	"time"
)

type User struct {
	UserId    int64
	Username  string
	Email     string
	CreatedAt int64
}

type UserModel struct {
	mu    sync.RWMutex
	users map[int64]*User
}

func NewUserModel() *UserModel {
	um := &UserModel{
		users: make(map[int64]*User),
	}
	// 初始化一些测试数据
	um.users[1] = &User{
		UserId:    1,
		Username:  "alice",
		Email:     "alice@example.com",
		CreatedAt: time.Now().Unix(),
	}
	um.users[2] = &User{
		UserId:    2,
		Username:  "bob",
		Email:     "bob@example.com",
		CreatedAt: time.Now().Unix(),
	}
	return um
}

func (m *UserModel) GetUser(userId int64) (*User, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	user, ok := m.users[userId]
	return user, ok
}

