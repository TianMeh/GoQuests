package models

import "time"

type Quest struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Reward      int       `json:"reward"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type Session struct {
	ID      uint      `json:"id" gorm:"primary_key"`
	UserID  uint      `json:"user_id" gorm:"index"`
	Expires time.Time `json:"expires_at"`
	Token   string    `json:"token"`
}

func (s *Session) isExpired() bool {
	return s.Expires.Before(time.Now())
}
