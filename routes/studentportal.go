package routes

import (
	"fmt"
	"html/template"
	"net/http"
)

func ApproveCource(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))
	uuid := r.URL.Query().Get("uuid")
	cource_name := r.URL.Query().Get("cource_name")
	var setresult_out string



	result := ApplyForCource(uuid, cource_name)

	if result {
		setresult_out = "cource_applied_btn"
	}

	err := tpl.ExecuteTemplate(w, setresult_out, nil)

	if err != nil {
		fmt.Fprintln(w, "Something Went Wrong!!!!")
	}

}
