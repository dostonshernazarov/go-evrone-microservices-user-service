package entity

import "time"

type Users struct {
	GUID        string
	FullName      string
	Username       string
	Email        string
	Password        string
	Bio       string
	Website  string
	Role      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}