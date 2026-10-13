package model

import "time"

type Account struct {
	UserId        int        `db:"id"`
	Name          string     `db:"name"`
	Email         string     `db:"email"`
	Bio           *string    `db:"bio"`
	User_location *string    `db:"user_location"`
	Position      *string    `db:"position"`
	Password      string     `db:"password"`
	Role          string     `db:"role"`
	Avatar_url    *string    `db:"avatar_url"`
	Status        string     `db:"status"`
	Terms         bool       `db:"terms"`
	Created_at    time.Time  `db:"created_at"`
	Updated_at    *time.Time `db:"updated_at"`
}

type AccountDetail struct {
	Account
	EventJoin       []int `db:"event_join"`
	CommunityMember []int `db:"community_member"`
	EventSave       []int `db:"event_save"`
}
