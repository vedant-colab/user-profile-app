package service

import (
	"fmt"
	"time"
	"user-profile-app/internals/models"
	"user-profile-app/internals/repository"
	"user-profile-app/internals/utils"
)

type AuthService struct {
	UserRepo    *repository.UserRepository
	SessionRepo *repository.SessionRepository
	ProfileRepo *repository.ProfileRepository
}

func (a *AuthService) GoogleLoginService(googleUser *models.GoogleUser) (*models.User, error) {
	user, err := a.UserRepo.GetUserByGoogleID(googleUser.ID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	user, err = a.UserRepo.GetUserByEmail(googleUser.Email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		err = a.UserRepo.UpdateGoogleID(user.ID, googleUser.ID)
		if err != nil {
			return nil, err
		}
		user.GoogleID = &googleUser.ID
		return user, nil
	}
	newUser := &models.User{
		Email:    googleUser.Email,
		GoogleID: &googleUser.ID,
	}

	userID, err := a.UserRepo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	newUser.ID = userID
	return newUser, nil
}

func (a *AuthService) CreateSession(id int64) (string, error) {
	now := time.Now()
	expiry := now.Add(24 * time.Hour)
	sid, err := a.SessionRepo.CreateSession(id, expiry)
	if err != nil {
		return "", err
	}
	return sid, nil
}

func (a *AuthService) CreateManualUser(email, password string) (int64, error) {
	password_hash := utils.HashPassword(password)
	user := &models.User{
		PasswordHash: &password_hash,
		Email:        email,
	}
	id, err := a.UserRepo.CreateUser(user)
	if err != nil {
		return 0, err
	}
	return id, nil

}

func (a *AuthService) FetchProfileById(userid int64) (*models.Profile, error) {
	profile, err := a.ProfileRepo.GetProfileById(userid)
	if err != nil {
		return nil, err
	}
	return profile, nil

}

func (a *AuthService) GetUserIDFromSession(sid string) (int64, error) {
	userid, err := a.SessionRepo.GetUserIDFromSession(sid)
	if err != nil {
		return 0, err
	}
	return userid, nil

}

func (a *AuthService) GetUserByID(userid int64) (*models.User, error) {
	user, err := a.UserRepo.GetUserByID(userid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (a *AuthService) GetUserByEmail(email string) (*models.User, error) {
	user, err := a.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *AuthService) UpateUserByID(userid int64, email string) error {
	user, err := a.UserRepo.GetUserByID(userid)
	if err != nil {
		return err
	}
	if email == user.Email {
		return nil
	}
	existingUser, err := a.UserRepo.GetUserByEmail(email)
	if err == nil && existingUser.ID != userid {
		return fmt.Errorf("email already in use")
	}

	err = a.UserRepo.UpdateUser(userid, email)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) DeleteSession(sid string) error {
	err := a.SessionRepo.DeleteSession(sid)
	if err != nil {
		return err
	}
	return nil
}
