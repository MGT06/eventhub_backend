package dto

import "time"

type Notification struct {
	Id        int        `json:"id"`
	UserId    int        `json:"user_id"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	Type      string     `json:"type"`
	Read_at   *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type CreateNotification struct {
	UserId  int    `json:"user_id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Type    string `json:"type"`
}
