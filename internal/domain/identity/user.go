package identity

import (
	"errors"
	"time"
)

var (
	ErrInvalidEmail    = errors.New("invalid user email")
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
)

type User struct {
	id           string
	name         string
	email        string
	passwordHash string
	roles        []Role
	isActive     bool
	createdAt    time.Time
}

func NewUser(id, name, email, passwordHash string, roles []Role) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}
	if passwordHash == "" {
		return nil, errors.New("password hash cannot be empty")
	}

	if len(roles) == 0 {
		roles = []Role{RoleUser}
	}

	return &User{
		id:           id,
		name:         name,
		email:        email,
		passwordHash: passwordHash,
		roles:        roles,
		isActive:     true,
		createdAt:    time.Now(),
	}, nil
}

// Getters
func (u *User) ID() string           { return u.id }
func (u *User) Name() string         { return u.name }
func (u *User) Email() string        { return u.email }
func (u *User) PasswordHash() string { return u.passwordHash }
func (u *User) Roles() []Role        { return u.roles }
func (u *User) IsActive() bool       { return u.isActive }
func (u *User) CreatedAt() time.Time { return u.createdAt }