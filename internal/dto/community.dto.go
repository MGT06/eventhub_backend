package dto

type Community struct {
	Id             int      `json:"id"`
	Community_name string   `json:"community_name"`
	Description    string   `json:"description"`
	Image_url      string   `json:"image_url"`
	TotalMembers   int      `json:"total_members"`
	UpcomingEvent  int      `json:"upcoming_events"`
	Category       []string `json:"category"`
}

type CommunityMembers struct {
	MemberName string `json:"member_name"`
}
