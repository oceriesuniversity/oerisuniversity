package routes

import (
	"html/template"
	"log"
	"net/http"
	"fmt"
)

func Programs(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	programs, err := GetAllProgramData()

	if err != nil {
		fmt.Println("Failed to get programs available")
	}

	err = tpl.ExecuteTemplate(w, "programs.html", programs)

	if err != nil {
		log.Fatal(err)
	}
}

func Programcards(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))
	var programsAvailable bool
	user_name, err := GetUserName(r)

	

	programdata, errout := GetAllProgramData()

	if errout != nil {
		programsAvailable = false
	} else {
		programsAvailable = true
	}

	data_out := AdminLandingData{
		Admin_Name:         user_name,
		ProgramD:      programdata,
		DataAvailable: programsAvailable,
	}

	//debug failure to laod templates

	err = tpl.ExecuteTemplate(w, "A_programs.html", data_out)

	if err != nil {
		log.Fatal(err)
	}
}
