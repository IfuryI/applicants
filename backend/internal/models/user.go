package models

// User структура юзера
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	IsAdmin  int    `json:"isAdmin"`
}

// UserNoPassword структура юзера без пароля
type UserNoPassword struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	IsAdmin  string `json:"isAdmin"`
}

// FromUser из User в UserNoPassword
func FromUser(user User) UserNoPassword {
	return UserNoPassword{
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
	}
}
