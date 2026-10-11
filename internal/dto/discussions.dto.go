package dto

import "time"

type Author struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url"`
}

type DiscussionCommunity struct {
	Id          int                    `json:"id"`
	CommunityId int                    `json:"community_id"`
	ParentId    *int                   `json:"parent_id"`
	Message     string                 `json:"message"`
	Author      Author                 `json:"author"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   *time.Time             `json:"updated_at"`
	Replies     []*DiscussionCommunity `json:"replies"`
}

type DiscussionEvent struct {
	Id        int                `json:"id"`
	EventId   int                `json:"event_id"`
	ParentId  *int               `json:"parent_id"`
	Message   string             `json:"message"`
	Author    Author             `json:"author"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at"`
	Replies   []DiscussionEvent `json:"replies"`
}
