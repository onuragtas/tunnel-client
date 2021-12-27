package models

import "time"

type Domain struct {
	Data struct {
		Domains []DomainItem `json:"domains"`
	} `json:"data"`
	Success bool `json:"success"`
}

type DomainItem struct {
	CreatedAt time.Time   `json:"CreatedAt"`
	UpdatedAt time.Time   `json:"UpdatedAt"`
	DeletedAt interface{} `json:"DeletedAt"`
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	Domain    string      `json:"domain"`
	Port      int         `json:"port"`
	KeepAlive int         `json:"keep_alive"`
}
