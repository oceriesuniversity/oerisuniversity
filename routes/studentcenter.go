package routes

import (
	"html/template"
	"log"
	"net/http"
)


func StudentCenter(w http.ResponseWriter, r *http.Request){
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	//debug failure to laod templates

	err := tpl.ExecuteTemplate(w, "studentcenter.html", nil)

	if err != nil {
		log.Fatal(err)
	}
}

