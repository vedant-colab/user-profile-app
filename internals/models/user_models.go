package models

type User struct {
	ID           int64
	Email        string
	PasswordHash *string
	GoogleID     *string
}

type GoogleUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ManualUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
