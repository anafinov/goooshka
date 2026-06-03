package domain

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // hashed
}

type Request struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	Status    string    `json:"status"` // "pending", "ready"
	CreatedAt time.Time `json:"created_at"`
}
