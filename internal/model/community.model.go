package model

import "time"

type Community struct {
	Id_community        int       `db:"id"`
	Community_name      string    `db:"community_name"`
	Description         string    `db:"description"`
	Image_community_url string    `db:"image_community_url"`
	Created_at          time.Time `db:"created_at"`
	Updated_at          time.Time `db:"updated_at"`
}

type CommunityDetail struct {
	Community
	category
	Account
	TotalMembers int	`json:"total_members"`
	UpcomingEvent int 	`json:"upcoming_events"`
}