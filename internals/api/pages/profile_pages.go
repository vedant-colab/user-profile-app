package pages

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/vedant-colab/github.com/vedant-colab/user-profile-app/internals/models"
)

func (h *PageHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {

	req, _ := http.NewRequest("GET", h.BaseURL+"/api/profile", nil)

	if c, err := r.Cookie("sid"); err == nil {
		req.AddCookie(c)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "failed to fetch profile", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var data models.ProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid response", http.StatusInternalServerError)
		return
	}
	h.App.Tpl.ExecuteTemplate(w, "profile.html", data)
}

func (h *PageHandler) CreateProfileHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		req, _ := http.NewRequest("POST", h.BaseURL+"/api/profile/create",
			strings.NewReader(url.Values{
				"full_name": {strings.TrimSpace(r.FormValue("full_name"))},
				"phone":     {strings.TrimSpace(r.FormValue("phone"))},
			}.Encode()),
		)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		if c, err := r.Cookie("sid"); err == nil {
			req.AddCookie(c)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			http.Error(w, "failed to create profile", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	h.App.Tpl.ExecuteTemplate(w, "profile_create.html", nil)
}

func (h *PageHandler) EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		req, _ := http.NewRequest("GET", h.BaseURL+"/api/profile/update", nil)
		if c, err := r.Cookie("sid"); err == nil {
			req.AddCookie(c)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "failed to fetch profile", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var data models.ProfileResponse
		json.NewDecoder(resp.Body).Decode(&data)

		h.App.Tpl.ExecuteTemplate(w, "profile_edit.html", data)
		return
	}
	if r.Method == "POST" {

		form := url.Values{
			"full_name": {strings.TrimSpace(r.FormValue("full_name"))},
			"phone":     {strings.TrimSpace(r.FormValue("phone"))},
			"email":     {strings.TrimSpace(r.FormValue("email"))},
		}

		req, _ := http.NewRequest("POST", h.BaseURL+"/api/profile/update", strings.NewReader(form.Encode()))

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if c, err := r.Cookie("sid"); err == nil {
			req.AddCookie(c)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "update failed", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "update failed", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
}
