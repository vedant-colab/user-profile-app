package pages

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/vedant-colab/user-profile-app/internals/utils"
)

func (h *PageHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		resp, err := http.PostForm(
			h.BaseURL+"/api/login",
			url.Values{
				"email":    {strings.TrimSpace(r.FormValue("email"))},
				"password": {strings.TrimSpace(r.FormValue("password"))},
			},
		)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"error": "internal error",
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "login failed", http.StatusUnauthorized)
			return
		}
		for _, c := range resp.Cookies() {
			http.SetCookie(w, c)
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	h.App.Tpl.ExecuteTemplate(w, "login.html", nil)
}

func (h *PageHandler) GoogleHandler(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateStateToken()

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})
	url := h.Oauth.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *PageHandler) Callback(w http.ResponseWriter, r *http.Request) {

	apiURL := h.BaseURL + "/api/auth/google/callback?" + r.URL.RawQuery

	req, _ := http.NewRequest("GET", apiURL, nil)

	if c, err := r.Cookie("oauth_state"); err == nil {
		req.AddCookie(c)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("callback: %v", err)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "oauth failed",
		})
		return
	}
	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		http.SetCookie(w, c)
	}

	var apiResp struct {
		Status string `json:"status"`
		Data   struct {
			UserID        int64 `json:"user_id"`
			ProfileExists bool  `json:"profileExists"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		http.Error(w, "invalid response", http.StatusInternalServerError)
		return
	}

	if !apiResp.Data.ProfileExists {
		http.Redirect(w, r, "/profile/create", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *PageHandler) Logout(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest("POST", h.BaseURL+"/api/logout", nil)
	cookie, err := r.Cookie("sid")
	if err == nil {
		req.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "logout failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	for _, c := range resp.Cookies() {
		http.SetCookie(w, c)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *PageHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		resp, err := http.PostForm(
			h.BaseURL+"/api/signup",
			url.Values{
				"email":    {strings.TrimSpace(r.FormValue("email"))},
				"password": {strings.TrimSpace(r.FormValue("password"))},
			},
		)

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		log.Printf("[signup - api status code]: %v", resp.StatusCode)
		if resp.StatusCode != http.StatusOK {
			http.Error(w, "signup failed", http.StatusBadRequest)
			return
		}
		for _, c := range resp.Cookies() {
			http.SetCookie(w, c)
		}

		http.Redirect(w, r, "/profile/create", http.StatusSeeOther)
		return
	}

	h.App.Tpl.ExecuteTemplate(w, "signup.html", nil)
}
