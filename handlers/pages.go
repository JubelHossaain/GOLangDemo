package handlers

import (
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

// HomeHandler serves the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

// UsersPageHandler serves the user management page
func UsersPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "users.html", nil)
}

// MessagesPageHandler serves the message management page
func MessagesPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "messages.html", nil)
}
