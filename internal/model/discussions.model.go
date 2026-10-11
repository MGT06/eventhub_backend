package model

import "time"

type DiscussCommunity struct {
	DiscussId   int        `db:"id"`
	CommunityId int        `db:"community_id"`
	AccountId   int        `db:"account_id"`
	ParentId    *int       `db:"parent_id"`
	Message     string     `db:"message"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

type DiscussionsCommunityDetail struct {
	DiscussCommunity
	Account
}

type DiscussEvent struct {
	DiscussId int        `db:"id"`
	EventId   int        `db:"event_id"`
	AccountId int        `db:"account_id"`
	ParentId  *int       `db:"parent_id"`
	Message   string     `db:"message"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type DiscussionsEventDetail struct {
	DiscussEvent
	Account
}
