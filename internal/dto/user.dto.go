package dto

import "mime/multipart"

type UserProfile struct {
	Bio           *string `json:"bio"`
	User_location *string `json:"user_location"`
	Position      *string `json:"position"`
}

type EditUserProfile struct {
	Name          string                `form:"name"`
	Bio           *string               `form:"bio"`
	User_location *string               `form:"user_location" binding:"omitempty,max=50"`
	Position      *string               `form:"position" binding:"omitempty,max=50"`
	Avatar        *multipart.FileHeader `form:"avatar"`
}

type ResponseEditUserProfile struct {
	Name          string  `json:"name"`
	Bio           *string `json:"bio"`
	User_location *string `json:"user_location"`
	Position      *string `json:"position"`
	Avatar        *string `json:"avatar_url"`
}

type ChangePassword struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

type UserHeaderInfo struct {
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Avatar_url      *string `json:"avatar_url"`
	Role            string  `json:"role"`
	EventJoin       []int   `json:"event_join"`
	CommunityMember []int   `json:"community_member"`
	EventSave       []int   `json:"event_save"`
}
