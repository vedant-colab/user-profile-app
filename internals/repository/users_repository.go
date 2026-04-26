package repository

import (
	"database/sql"
	"fmt"

	"github.com/vedant-colab/github.com/vedant-colab/user-profile-app/internals/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (u *UserRepository) GetUserByGoogleID(id string) (*models.User, error) {
	var user models.User
	err := u.DB.QueryRow("SELECT id, email, password_hash, google_id FROM users WHERE google_id = ?", id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.GoogleID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := u.DB.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) UpdateGoogleID(id int64, googleId string) error {
	_, err := u.DB.Exec("UPDATE users (id, google_id) VALUES (?, ?)", id, googleId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) CreateUser(user *models.User) (int64, error) {
	res, err := u.DB.Exec("INSERT INTO users (email, password_hash, google_id) VALUES (?, ?, ?)", user.Email, user.PasswordHash, user.GoogleID)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (u *UserRepository) GetUserByID(userid int64) (*models.User, error) {
	var user models.User
	err := u.DB.QueryRow("SELECT id, email, password_hash,google_id FROM users WHERE id = ?", userid).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.GoogleID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) UpdateUser(userid int64, email string) error {
	res, err := u.DB.Exec("UPDATE users SET email = ? WHERE id = ? AND google_id IS NULL", email, userid)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("email update not allowed or user not found")
	}
	return nil
}
