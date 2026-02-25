package entity

import "time"

type User struct {
	id        string
	createdAt time.Time
}

func NewUser(id string, createAt time.Time) *User {
	return &User{
		id:        id,
		createdAt: createAt,
	}
}

func (u *User) ID() string {
	return u.id
}
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}
