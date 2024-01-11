package models

import "time"

type Domain struct {
	Data struct {
		Domains []DomainItem `json:"domains"`
	} `json:"data"`
	Success bool `json:"success"`
}

type UserInfo struct {
	Data struct {
		User struct {
			ID        int         `json:"ID"`
			CreatedAt time.Time   `json:"CreatedAt"`
			UpdatedAt time.Time   `json:"UpdatedAt"`
			DeletedAt interface{} `json:"DeletedAt"`
			Id        int         `json:"id"`
			Username  string      `json:"username"`
			Email     string      `json:"email"`
			Password  string      `json:"password"`
			Subscribe int         `json:"subscribe"`
		} `json:"user"`
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
