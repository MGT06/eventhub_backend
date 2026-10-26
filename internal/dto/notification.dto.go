package dto

import "time"

type Notification struct {
	Id       int        `json:"id"`
	UserName string     `json:"userName"`
	Title    string     `json:"title"`
	Message  string     `json:"message"`
	Type     string     `json:"type"`
	Read_at  *time.Time `json:"read_at"`
}
