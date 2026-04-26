package repository

import (
	"database/sql"
	"time"

	"log"

	"github.com/vedant-colab/user-profile-app/internals/utils"
)

type SessionRepository struct {
	DB *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		DB: db,
	}
}

func (s *SessionRepository) CreateSession(userid int64, expiry time.Time) (string, error) {
	sid := utils.GenerateSessionID()
	_, err := s.DB.Exec("INSERT INTO sessions (id, user_id, created_at, expires_at) VALUES (?,?, ?, ?)", sid, userid, time.Now(), expiry)
	if err != nil {
		return "", err
	}
	return sid, nil
}

func (s *SessionRepository) GetUserIDFromSession(sid string) (int64, error) {
	var userid int64
	err := s.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ?", sid).Scan(&userid)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return userid, nil
}

func (s *SessionRepository) DeleteSession(sid string) error {
	_, err := s.DB.Exec("DELETE FROM sessions WHERE id = ?", sid)
	return err
}

func (s *SessionRepository) DeleteSessionByUserID(userid int64) error {
	_, err := s.DB.Exec("DELETE FROM sessions WHERE user_id = ?", userid)
	log.Println(err)
	return err
}
