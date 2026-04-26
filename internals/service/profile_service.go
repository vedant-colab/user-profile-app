package service

import (
	"user-profile-app/internals/repository"
)

type ProfileService struct {
	UserRepo    *repository.UserRepository
	SessionRepo *repository.SessionRepository
	ProfileRepo *repository.ProfileRepository
}

func (p *ProfileService) CreateProfile(userid int64, fullName, phone string) error {
	err := p.ProfileRepo.CreateProfile(userid, fullName, phone)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProfileService) UpdateProfile(userid int64, fullName, phone string) error {
	profile, err := p.ProfileRepo.GetProfileById(userid)
	if err != nil {
		return err
	}

	if profile == nil {
		err = p.ProfileRepo.CreateProfile(userid, fullName, phone)
		if err != nil {
			return err
		}
		return nil
	}

	err = p.ProfileRepo.UpdateProfile(userid, fullName, phone)
	if err != nil {
		return err
	}
	return nil
}
