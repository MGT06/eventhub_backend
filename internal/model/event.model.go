package model

import "time"

type Event struct {
	Id_event          int        `db:"id"`
	Organizer_id      int        `db:"organizer_id"`
	Community_id      *int       `db:"community_id"`
	Location_event_id int        `db:"location_event_id"`
	Title             string     `db:"title"`
	Description       string     `db:"description"`
	Image_event_url   string     `db:"image_event_url"`
	Start_at          time.Time  `db:"start_at"`
	End_at            time.Time  `db:"end_at"`
	Format            string     `db:"format"`
	Capacity          int        `db:"capacity"`
	Speakers          string     `db:"speakers"`
	Created_at        time.Time  `db:"created_at"`
	Updated_at        *time.Time `db:"updated_at"`
}

type EventDetail struct {
	Event
	Account
	category
	locationEvent
	Attendees int	`json:"attendees"`
}
