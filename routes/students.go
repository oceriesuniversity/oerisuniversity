package routes

import (
	"html/template"
	"log"
	"net/http"
)

func Students(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "studentdataadmin.html", nil)

	if err != nil {
		log.Fatal(err)
	}
}
