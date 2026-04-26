package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/vedant-colab/user-profile-app/internals/api/response"
	"github.com/vedant-colab/user-profile-app/internals/models"
	"github.com/vedant-colab/user-profile-app/internals/utils"
)

func (h *Handler) LoginAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.Error(w, http.StatusMethodNotAllowed, "invalid method")
		return
	}

	r.ParseForm()

	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	if email == "" || password == "" {
		response.Error(w, http.StatusBadRequest, "email and password required")
		return
	}

	user, err := h.AuthService.GetUserByEmail(email)
	if err != nil || user == nil {
		response.Error(w, http.StatusUnauthorized, "invalid user")
		return
	}

	if !utils.CompareHashPasswords(*user.PasswordHash, password) {
		response.Error(w, http.StatusUnauthorized, "wrong password")
		return
	}

	sid, err := h.AuthService.CreateSession(user.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Domain:   "",
		SameSite: http.SameSiteNoneMode,
	})

	w.Header().Set("Content-Type", "application/json")
	response.Success(w, map[string]string{
		"status": "ok",
	})
}

func (h *Handler) SignupAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.Error(w, http.StatusMethodNotAllowed, "invalid method")
		return
	}

	r.ParseForm()

	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	if email == "" || password == "" {
		response.Error(w, http.StatusBadRequest, "email and password required")
		return
	}

	user, err := h.AuthService.GetUserByEmail(email)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "error fetching user")
		return
	}
	if user != nil {
		response.Error(w, http.StatusConflict, "error fetching user")
		return
	}

	userID, err := h.AuthService.CreateManualUser(email, password)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "error creating user")
		return
	}

	sid, err := h.AuthService.CreateSession(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create session")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	response.Success(w, map[string]string{
		"status": "created",
	})
}

func (h *Handler) LogoutAPI(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sid")
	if err == nil {
		_ = h.AuthService.DeleteSession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	response.Success(w, map[string]string{
		"status": "logged out",
	})
}
func (h *Handler) GoogleCallbackAPI(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("oauth_state")
	if err != nil || cookie.Value != r.FormValue("state") {
		response.Error(w, http.StatusUnauthorized, "invalid oauth state")
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		response.Error(w, http.StatusBadRequest, "missing authorization code")
		return
	}
	token, err := h.Oauth.Exchange(r.Context(), code)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "code exchange failed")
		return
	}
	client := h.Oauth.Client(r.Context(), token)
	resp, err := client.Get(os.Getenv("GOOGLE_USER_INFO"))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch user info")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		response.Error(w, http.StatusBadGateway, "google user info fetch failed")
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to read user info")
		return
	}

	var googleUser models.GoogleUser
	if err := json.Unmarshal(data, &googleUser); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to parse user info")
		return
	}
	user, err := h.AuthService.GoogleLoginService(&googleUser)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	sid, err := h.AuthService.CreateSession(user.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
	profile, err := h.AuthService.FetchProfileById(user.ID)
	log.Println("profile userid", profile.UserID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch profile")
		return
	}
	response.Success(w, map[string]interface{}{
		"user_id":       user.ID,
		"profileExists": profile != nil,
	})
}
