package models

import "time"

type User struct {
	Id				string
    Email			string
    FullName		string
    PasswordHash	string
    PasswordSalt	string
    CreatedTime		time.Time
    LastLoginTime	time.Time
    IsDisabled		bool
    IsActivated		bool
}

type UserSession struct {
	SessionKey		string
	UserId			string
	LoginTime 		time.Time
	LastSeenTime 	time.Time
}