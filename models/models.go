package models

import "time"

type Category struct {
	Id 			int 		`json:"id"`
	UserId		int			`json:"uid"`
	Name		string		`json:"name"`	
	Picture		string		`json:"picture"`
	CreatedTime	time.Time 	`json:"time"`
}

type Container struct {
	Id 			int 		`json:"id"`
	UserId		int			`json:"uid"`
	CategoryId	int			`json:"cid"`
	Name		string		`json:"name"`
	Description	string		`json:"desc"`
	Time		time.Time 	`json:"time"`
	IsPublic	bool		`json:"public"`
}

type Date struct {
	Id 			int 		`json:"id"`
	ContainerId	int			`json:"cid"`
	Name		string		`json:"name"`
	Type 	 	string		`json:"type"`
	Time		time.Time 	`json:"time"`
	Icon		string		`json:"icon"`
	Content		string		`json:"content"`
}

type Categories []Category
type Containers []Container
type Dates []Date