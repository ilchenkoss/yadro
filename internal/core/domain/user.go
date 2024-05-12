package domain

type UserRole string

const (
	Admin    UserRole = "admin"
	Ordinary UserRole = "ordinary"
)

type User struct {
	ID       uint64
	Login    string
	Password string
	Salt     string
	Role     UserRole
}
