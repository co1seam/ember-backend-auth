package entity

import "time"

type User struct {
	ID       int       `json:"-"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

type SendOtpRequest struct {
	Email string `json:"email"`
}
