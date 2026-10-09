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
	Name   string    `json:"name"`
	Role   string    `json:"role"`
	Status string    `json:"status"`
	Joined time.Time `json:"joined"`
}
