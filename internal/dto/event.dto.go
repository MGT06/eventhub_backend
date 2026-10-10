package dto

import (
	"mime/multipart"
	"time"
)

type Event struct {
	Id            int       `json:"id"`
	OrganizerName string    `json:"organizer"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Image_url     string    `json:"image_url"`
	Start_at      time.Time `json:"start_at"`
	End_at        time.Time `json:"end_at"`
	Format        string    `json:"format"`
	Location      string    `json:"location"`
	Attendees 	  int 		`json:"attendees"`
	Capacity      int    `json:"capacity"`
	Speakers      string    `json:"speakers"`
	Categories    []string  `json:"categories"`
}

type AddEvent struct {
	OrganizerId int                  `form:"organizer_id"`
	CommunityId *int                 `form:"community_id"`
	Title       string               `form:"title"`
	Description string               `form:"description"`
	Image       multipart.FileHeader `form:"image"`
	Start_at    time.Time            `form:"start_at" example:"2026-10-05T10:30:00Z"`
	End_at      time.Time            `form:"end_at" example:"2026-10-05T10:30:00Z"`
	Format      string               `form:"format"`
	LocationId  int                  `form:"location_event_id"`
	Capacity    int                  `form:"capacity"`
	Speakers    string               `form:"speakers"`
	Categories  []int                `form:"categories"`
}

type EditEvent struct {
	CommunityId *int                  `form:"community_id"`
	Title       *string               `form:"title"`
	Description *string               `form:"description"`
	Image       *multipart.FileHeader `form:"image"`
	Start_at    *time.Time            `form:"start_at"`
	End_at      *time.Time            `form:"end_at"`
	Format      *string               `form:"format"`
	LocationId  *int                  `form:"location_event_id"`
	Capacity    *int                  `form:"capacity"`
	Speakers    *string               `form:"speakers"`
	Categories  []int                 `form:"categories"`
}

type Speaker struct {
	Name   string `json:"name"`
	Title  string `json:"title"`
	Avatar string `json:"avatar_url"`
}
