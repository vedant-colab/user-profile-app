package main

import (
	"log"
	"net/http"

	"github.com/vedant-colab/user-profile-app/internals/api/handlers"
	"github.com/vedant-colab/user-profile-app/internals/api/pages"
	"github.com/vedant-colab/user-profile-app/internals/app"
	"github.com/vedant-colab/user-profile-app/internals/config"
	middleware "github.com/vedant-colab/user-profile-app/internals/middlewares"
	"github.com/vedant-colab/user-profile-app/internals/repository"
	"github.com/vedant-colab/user-profile-app/internals/service"
)

func main() {

	config.InitializeConfig()
	googleOauth := config.InitializeOauth()

	db, err := repository.InitDB()

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	profileRepo := repository.NewProfileRepository(db)
	authService := service.AuthService{
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
		ProfileRepo: profileRepo,
	}

	profileService := service.ProfileService{
		ProfileRepo: profileRepo,
	}

	if err != nil {
		log.Fatalf("Error connecting database: %v", err)
	}

	template, err := config.GetTemplates()
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}
	app := &app.App{
		Tpl: template,
		DB:  db,
	}
	h := handlers.Handler{
		App:            app,
		Oauth:          googleOauth,
		AuthService:    authService,
		ProfileService: profileService,
	}

	pHand := pages.PageHandler{
		App:       app,
		APIClient: &http.Client{},
		Oauth:     googleOauth,
		BaseURL:   config.Cfg.BaseURL,
	}
	authMiddleware := middleware.AuthMiddleware(&h.AuthService)

	http.HandleFunc("/api/login", h.LoginAPI)
	http.HandleFunc("/api/signup", h.SignupAPI)
	http.HandleFunc("/api/auth/google/callback", h.GoogleCallbackAPI)
	http.Handle("/api/profile", authMiddleware(http.HandlerFunc(h.GetProfileAPI)))
	http.Handle("/api/profile/create", authMiddleware(http.HandlerFunc(h.CreateProfileAPI)))
	http.Handle("/api/profile/update", authMiddleware(http.HandlerFunc(h.UpdateProfileAPI)))

	http.HandleFunc("/", pHand.LoginHandler)
	http.HandleFunc("/signup", pHand.SignupHandler)
	http.HandleFunc("/auth/google", pHand.GoogleHandler)
	http.HandleFunc("/auth/google/callback", pHand.Callback)
	http.HandleFunc("/profile", pHand.ProfileHandler)
	http.HandleFunc("/profile/create", pHand.CreateProfileHandler)
	http.HandleFunc("/profile/update", pHand.EditProfileHandler)

	log.Printf("Server starting on port: %s", config.Cfg.PORT)
	http.ListenAndServe(":"+config.Cfg.PORT, nil)

}
