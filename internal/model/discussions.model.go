package model

import "time"

type DiscussCommunity struct {
	DiscussId   int        `db:""`
	CommunityId int        `db:""`
	AccountId   int        `db:""`
	ParentId    *int       `db:""`
	Message     string     `db:""`
	CreatedAt   time.Time  `db:""`
	UpdatedAt   *time.Time `db:""`
}

type DiscussionsCommunityDetail struct {
	DiscussCommunity
	Account
}