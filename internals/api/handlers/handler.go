package handlers

import (
	"user-profile-app/internals/app"
	"user-profile-app/internals/service"

	"golang.org/x/oauth2"
)

type Handler struct {
	App            *app.App
	Oauth          *oauth2.Config
	AuthService    service.AuthService
	ProfileService service.ProfileService
}
