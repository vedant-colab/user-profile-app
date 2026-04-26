package repository

import (
	"database/sql"
	"fmt"
	"user-profile-app/internals/models"
)

type ProfileRepository struct {
	DB *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{
		DB: db,
	}
}

func (p *ProfileRepository) GetProfileById(userid int64) (*models.Profile, error) {
	var profile models.Profile
	err := p.DB.QueryRow("SELECT id, user_id, full_name, phone FROM profiles WHERE user_id = ?", userid).Scan(&profile.ID, &profile.UserID, &profile.FullName, &profile.Phone)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (p *ProfileRepository) CreateProfile(userid int64, fullName string, phone string) error {
	_, err := p.DB.Exec("INSERT INTO profiles (user_id, full_name, phone) VALUES (?, ?, ?)", userid, fullName, phone)
	return err
}

func (p *ProfileRepository) UpdateProfile(userid int64, fullName string, phone string) error {
	res, err := p.DB.Exec("UPDATE profiles SET full_name = ?, phone = ?, updated_at = NOW() WHERE user_id = ?", fullName, phone, userid)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no profile found for user_id %d", userid)
	}

	return nil
}
