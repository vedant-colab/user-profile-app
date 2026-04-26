package pages

import (
	"net/http"

	"github.com/vedant-colab/github.com/vedant-colab/user-profile-app/internals/app"

	"golang.org/x/oauth2"
)

type PageHandler struct {
	App       *app.App
	APIClient *http.Client
	BaseURL   string
	Oauth     *oauth2.Config
}
