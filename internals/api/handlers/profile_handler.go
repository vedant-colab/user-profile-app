package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/vedant-colab/github.com/vedant-colab/user-profile-app/internals/api/response"
	middleware "github.com/vedant-colab/github.com/vedant-colab/user-profile-app/internals/middlewares"
	"github.com/vedant-colab/github.com/vedant-colab/user-profile-app/internals/models"
)

func (h *Handler) GetProfileAPI(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	profile, err := h.AuthService.FetchProfileById(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch profile")
		return
	}

	user, err := h.AuthService.GetUserByID(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	var profileData models.Profile
	if profile != nil {
		profileData = *profile
	}

	resp := models.ProfileResponse{
		User:       *user,
		Profile:    profileData,
		ManualUser: user.GoogleID == nil,
	}

	log.Println(user.Email)
	response.Success(w, resp)
}

func (h *Handler) CreateProfileAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.Error(w, http.StatusMethodNotAllowed, "invalid method")
		return
	}

	if err := r.ParseForm(); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid form")
		return
	}

	fullName := strings.TrimSpace(r.FormValue("full_name"))
	phone := strings.TrimSpace(r.FormValue("phone"))

	if fullName == "" {
		response.Error(w, http.StatusBadRequest, "full name is required")
		return
	}

	if len(phone) != 10 {
		response.Error(w, http.StatusBadRequest, "invalid phone")
		return
	}
	userID := middleware.GetUserID(r)
	profile, err := h.AuthService.FetchProfileById(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch profile")
		return
	}
	if profile != nil {
		response.Error(w, http.StatusConflict, "profile already exists")
		return
	}
	if err := h.ProfileService.CreateProfile(userID, fullName, phone); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create profile")
		return
	}
	response.Success(w, map[string]string{
		"status": "created",
	})
}

func (h *Handler) UpdateProfileAPI(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	user, err := h.AuthService.GetUserByID(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}
	if r.Method == "GET" {
		profile, err := h.AuthService.FetchProfileById(userID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "failed to fetch profile")
			return
		}

		var profileData models.Profile
		if profile != nil {
			profileData = *profile
		}

		resp := models.ProfileResponse{
			Profile:    profileData,
			User:       *user,
			ManualUser: user.GoogleID == nil,
		}

		response.Success(w, resp)
		return
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid form")
			return
		}
		fullName := strings.TrimSpace(r.FormValue("full_name"))
		phone := strings.TrimSpace(r.FormValue("phone"))
		email := strings.TrimSpace(r.FormValue("email"))

		if fullName == "" {
			response.Error(w, http.StatusBadRequest, "full name is required")
			return
		}

		if len(phone) != 10 {
			response.Error(w, http.StatusBadRequest, "invalid phone")
			return
		}
		if user.GoogleID != nil {
			email = user.Email
		}
		if user.GoogleID == nil {
			err = h.AuthService.UpateUserByID(userID, email)
			if err != nil {
				response.Error(w, http.StatusInternalServerError, "failed to update user")
				return
			}
		}

		err = h.ProfileService.UpdateProfile(userID, fullName, phone)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "failed to update profile")
			return
		}
		response.Success(w, map[string]string{
			"status": "updated",
		})
	}
}
