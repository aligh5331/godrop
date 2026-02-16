package entity

import (
	"auth/internal/domain"
	"strings"
	"time"
)

type HashedPassword string
type User struct {
	id        string
	name      string
	email     string
	password  HashedPassword
	active    bool
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(id, name, email string, password HashedPassword, createdAt time.Time, updatedAt time.Time) (*User, error) {

	var user = &User{
		id:        strings.TrimSpace(id),
		name:      strings.TrimSpace(name),
		email:     strings.TrimSpace(email),
		password:  password,
		createdAt: createdAt.UTC(),
		updatedAt: updatedAt.UTC(),
	}

	if err := user.validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) validate() error {

	if u.id == "" {
		return domain.ErrEmptyId
	}

	if u.name == "" {
		return domain.ErrEmptyName
	}

	if u.email == "" {
		return domain.ErrEmptyEmail
	}

	if u.password == "" {
		return domain.ErrEmptyPassword
	}

	return nil
}

func (u *User) SetHashedPassword(p HashedPassword, now time.Time) {
	u.password = p
	u.updatedAt = now.UTC()
}

func (u *User) UpdateProfile(name string, now time.Time) error {
	newUser := *u
	newUser.name = strings.TrimSpace(name)

	if err := newUser.validate(); err != nil {
		return err
	}

	*u = newUser
	u.updatedAt = now.UTC()
	return nil
}

func (u *User) Id() string {
	return u.id
}
func (u *User) Name() string {
	return u.name
}
func (u *User) Email() string {
	return u.email
}

func (u *User) Password() HashedPassword {
	return u.password
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) ChangePassword(hashedPassword HashedPassword) error {
	if strings.TrimSpace(string(hashedPassword)) == "" {
		return domain.ErrEmptyPassword
	}

	u.password = hashedPassword
	u.updatedAt = time.Now()
	return nil
}

func (u *User) ChangeName(newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return domain.ErrEmptyName
	}
	u.name = newName
	u.updatedAt = time.Now()
	return nil
}

func (u *User) ChangeEmail(newEmail string) error {
	newEmail = strings.TrimSpace(newEmail)
	if newEmail == "" {
		return domain.ErrEmptyEmail
	}
	u.email = newEmail
	u.updatedAt = time.Now()
	return nil
}
