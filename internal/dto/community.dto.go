package dto

type Community struct {
	Id             int      `json:"id"`
	Community_name string   `json:"community_name"`
	Description    string   `json:"description"`
	Image_url      string   `json:"image_url"`
	IsActive       bool     `json:"isActive"`
	TotalMembers   int      `json:"total_members"`
	UpcomingEvent  int      `json:"upcoming_events"`
	Category       []string `json:"category"`
}

type CommunityMembers struct {
	Name string `json:"name"`
	Avatar_url *string `json:"avatar_url"`
	Position *string `json:"position"`
}
