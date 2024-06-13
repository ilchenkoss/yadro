package domain

type UserRole string

const (
	Admin     UserRole = "admin"
	SuperUser UserRole = "super_user"
	Ordinary  UserRole = "ordinary"
)
