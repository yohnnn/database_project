package models

type User struct {
	ID       int64
	Username string
	Password string
	Email    string
	RoleID   int64
}
