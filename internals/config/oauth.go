package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func InitializeOauth() *oauth2.Config {
	GoogleOauthConfig := &oauth2.Config{
		ClientID:     Cfg.Oauth.GOOGLE_CLIENT_ID,
		ClientSecret: Cfg.Oauth.GOOGLE_CLIENT_SECRET,
		RedirectURL:  Cfg.Oauth.AUTH_REDIRECT_URL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return GoogleOauthConfig
}
