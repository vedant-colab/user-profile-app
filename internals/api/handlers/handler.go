package handlers

import (
	"github.com/vedant-colab/github.com/vedant-colab/user-profile-app/internals/app"
	"github.com/vedant-colab/github.com/vedant-colab/user-profile-app/internals/service"

	"golang.org/x/oauth2"
)

type Handler struct {
	App            *app.App
	Oauth          *oauth2.Config
	AuthService    service.AuthService
	ProfileService service.ProfileService
}
