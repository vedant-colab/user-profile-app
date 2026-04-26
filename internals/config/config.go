package config

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
)

type Oauth struct {
	GOOGLE_CLIENT_ID     string
	AUTH_REDIRECT_URL    string
	GOOGLE_CLIENT_SECRET string
	GOOGLE_USER_INFO     string
}

type Config struct {
	Oauth   Oauth
	DSN     string
	BaseURL string
	PORT    string
}

var Cfg Config

func GetWorkDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func GetTemplates() (*template.Template, error) {
	workDir, err := GetWorkDir()
	if err != nil {
		log.Fatalf("couldn't fetch working dir: %v", err)
		return nil, err
	}
	funcMap := template.FuncMap{
		"deref": func(s *string) string {
			if s == nil {
				return ""
			}
			return *s
		},
	}
	Tpl, err := template.New("").Funcs(funcMap).ParseGlob(strings.Join([]string{workDir, "internals", "templates", "*.html"}, "/"))
	return Tpl, err
}

func InitializeConfig() Config {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	GOOGLE_CLIENT_ID := os.Getenv("GOOGLE_CLIENT_ID")
	GOOGLE_CLIENT_SECRET := os.Getenv("GOOGLE_CLIENT_SECRET")
	AUTH_REDIRECT_URL := os.Getenv("AUTH_REDIRECT_URL")
	GOOGLE_USER_INFO := os.Getenv("GOOGLE_USER_INFO")

	baseUrl := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	Cfg = Config{
		DSN: dsn,
		Oauth: Oauth{
			GOOGLE_CLIENT_ID:     GOOGLE_CLIENT_ID,
			GOOGLE_CLIENT_SECRET: GOOGLE_CLIENT_SECRET,
			AUTH_REDIRECT_URL:    AUTH_REDIRECT_URL,
			GOOGLE_USER_INFO:     GOOGLE_USER_INFO,
		},
		BaseURL: baseUrl,
		PORT:    port,
	}
	return Cfg
}
