package app

import (
	"database/sql"
	"html/template"
)

type App struct {
	Tpl *template.Template
	DB  *sql.DB
}
