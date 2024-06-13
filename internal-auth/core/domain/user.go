package domain

type UserRole string

const (
	Admin     UserRole = "admin"
	SuperUser UserRole = "super_user"
	Ordinary  UserRole = "ordinary"
)

type User struct {
	ID       int64
	Login    string
	Password string
	Salt     string
	Role     UserRole
}
