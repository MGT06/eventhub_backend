package dto

import "time"

type DashboardOrganizer struct {
	TotalEvent    int `json:"totalEvent"`
	TotalAttendes int `json:"totalAttendes"`
	AVGFillRate   int `json:"AVGFillRate"`
}

type DashboardAdmin struct {
	TotalUsers       int `json:"totalUsers"`
	TotalEvents      int `json:"totalEvents"`
	TotalCommunities int `json:"totalCommunities"`
}

type AllUser struct {
	Id     int       `json:"id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	Status string    `json:"status"`
	Joined time.Time `json:"joined"`
}
