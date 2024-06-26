package repository

//
//type UserRepository struct {
//	db *database.DB
//}
//
//func NewUserRepository(db *database.DB) *UserRepository {
//	return &UserRepository{
//		db,
//	}
//}
//
//func (u *UserRepository) CreateUser(user *domain.User) error {
//	res, err := u.db.Exec("INSERT INTO users (login, password, salt, role) VALUES (?, ?, ?, ?) ON CONFLICT DO NOTHING",
//		user.Login, user.Password, user.Salt, user.Role)
//	if err != nil {
//		return fmt.Errorf("error creating user: %w", err)
//	}
//	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
//		return domain.ErrUserAlreadyExist
//	}
//	return nil
//}
//
//func (u *UserRepository) GetUserByLogin(login string) (*domain.User, error) {
//	row := u.db.QueryRow("SELECT id, login, password, salt, role FROM users WHERE login = ?", login)
//	var user domain.User
//	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Salt, &user.Role)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return nil, domain.ErrUserNotFound
//		}
//		return nil, fmt.Errorf("error getting user by login: %w", err)
//	}
//	return &user, nil
//}
//
//func (u *UserRepository) UpdateUser(user *domain.User) error {
//	_, err := u.db.Exec("UPDATE users SET login = ?, password = ?, salt = ?, role = ? WHERE id = ?",
//		user.Login, user.Password, user.Salt, user.Role, user.ID)
//	if err != nil {
//		return fmt.Errorf("error updating user: %w", err)
//	}
//	return nil
//}
