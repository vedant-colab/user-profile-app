package service_test

import (
	"database/sql"
	"testing"

	"github.com/vedant-colab/user-profile-app/internals/models"
	"github.com/vedant-colab/user-profile-app/internals/repository"
	"github.com/vedant-colab/user-profile-app/internals/service"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/user_profile_test_db")
	if err != nil {
		t.Fatalf("failed to connect test db: %v", err)
	}
	_, _ = db.Exec("DELETE FROM users")
	_, _ = db.Exec("DELETE FROM sessions")

	return db
}

func setupAuthService(t *testing.T) *service.AuthService {
	db := setupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	return &service.AuthService{
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
	}
}

func TestGoogleLogin_NewUser(t *testing.T) {
	authService := setupAuthService(t)

	googleUser := &models.GoogleUser{
		Email: "test@gmail.com",
		Name:  "Test User",
	}

	user, err := authService.GoogleLoginService(googleUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatalf("expected user, got nil")
	}

	if user.Email != googleUser.Email {
		t.Fatalf("expected email %s, got %s", googleUser.Email, user.Email)
	}
}

func TestGoogleLogin_ExistingUser(t *testing.T) {
	authService := setupAuthService(t)

	googleUser := &models.GoogleUser{
		Email: "existing@gmail.com",
		Name:  "Existing User",
	}
	user1, err := authService.GoogleLoginService(googleUser)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	user2, err := authService.GoogleLoginService(googleUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user1.ID != user2.ID {
		t.Fatalf("expected same user, got different IDs")
	}
}
