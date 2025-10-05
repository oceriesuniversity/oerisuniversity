package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kamukwamba/oerisuniversity/dbcode"
)

type ProgramAvailable struct {
	Available bool
}



type AllStudentCources struct {
	Student_dat StudentInfo
	All         []AllCourceData
}

type AllCourceData struct {
	ProgramStruct ProgramStruct
	Cource_Struct []CourceStruct
}




type ProgramStruct struct {
	UUID           string
	Student_UUID   string
	Program_Name   string
	First_Name     string
	Last_Name      string
	Email          string
	Payment_Method string
	Paid           string
	Approved       bool
	Applied        bool
	Completed      bool
	Date           string
	Admin_ID       string
}

type CourceStruct struct {
	UUID                  string
	Student_UUID          string
	Cource_Name           string
	Course_Code 		string
	Book                  string
	Module                string
	Video                 string
	Applied               bool
	Approved              bool
	Examined              bool
	Continuorse_Assesment string
	Completed             bool
	Date                  string
}

type StudentCourse struct {
	Available        bool
	StInfo           StudentInfo
	AllCourceDataOut []AllCourceData
	Admin            AdminInfo
	Admin_Name string
}

func ValidateSudent(email_in, password_in string) (bool, string) {
	isstudent := true
	dbread := dbcode.SqlRead().DB
	stmt, err := dbread.Prepare("select uuid, student_uuid, email, password from studentcridentials where email = ?")

	if err != nil {
		isstudent = false
		fmt.Println("First err", err)
	}

	defer stmt.Close()

	var uuid string
	var student_uuid string
	var email string
	var password string

	fmt.Println("Email: ", email_in)
	fmt.Println("Password: ", password_in)

	err = stmt.QueryRow(email_in).Scan(&uuid, &student_uuid, &email, &password)

	if err != nil {
		fmt.Println(":Second err this is the err: ", err)
		// log.Fatal(err)
		isstudent = false
	}

	fmt.Println(uuid, student_uuid, email, email, password)

	compareHashedKeys := CheckPassword(password, password_in)

	if compareHashedKeys != true {
		isstudent = false
	}

	return isstudent, student_uuid

}



func UpdateProgram(student_uuid, table_name string) bool {
	var programcompleted bool

	dbread := dbcode.SqlRead()

	// .dbread"UPDATE artist_t SET check_s = ? WHERE artist_n = ?", "2021-05-20", 42

	approval_string := fmt.Sprintf("UPDATE %s SET completed = ? WHERE student_uuid = ?", table_name)
	stmt, err := dbread.DB.Prepare(approval_string)

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	_, erre := stmt.Exec(true, student_uuid)

	if erre != nil {
		log.Fatal(err)
	}

	// _, err := dbread.DB.Exec("update  acams set accepted = ? where student_uuid = ?", "true", student_uuid)
	return programcompleted

}

func Update(student_uuid, table_name string) bool {

	dbread := dbcode.SqlRead()
	updated := true

	approval_string := fmt.Sprintf("UPDATE %s SET approved = ? WHERE student_uuid = ?", table_name)

	stmt, err := dbread.DB.Prepare(approval_string)

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, erre := stmt.Exec(true, student_uuid)

	if erre != nil {
		log.Fatal(err)
	}

	// _, err := dbread.DB.Exec("update  acams set accepted = ? where student_uuid = ?", "true", student_uuid)

	return updated

}



func GetStudentPrograms(student_uuid string) []string {
	dbread := dbcode.SqlRead()

	var listout []string

	stmt, err := dbread.DB.Prepare("select program_list from studentprogramlist where student_uuid = ?")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var program_list string

	err = stmt.QueryRow(student_uuid).Scan(&program_list)

	trimedlist := strings.Trim(program_list, "[]")
	list_out := strings.Split(trimedlist, ",")

	for _, item := range list_out {
		trimedlistone := strings.Trim(item, "\"")
		trimedlisttwo := strings.Trim(trimedlistone, " \"")

		if len(trimedlisttwo) > 1 {
			listout = append(listout, trimedlisttwo)

		} else {
			continue
		}
	}

	if err != nil {
		fmt.Println("FAILED TO GET STUDENT PROGRAM LIST")
	}



	return listout
}

func StudentPortal(w http.ResponseWriter, r *http.Request) {

	studentuuid := r.URL.Query().Get("student")

	tpl = template.Must(template.ParseGlob("templates/*.html"))
	var programdataout []AllCourceData
	var studentinfo StudentInfo

	studentprogramlist := GetStudentPrograms(studentuuid)//Get The list of programs applied for 

	programdataout, present := GetStudentProgramData(studentprogramlist, studentuuid)
	
	studentinfo = GetStudentAllDetails(studentuuid)

	students_data := StudentCourse{
		Available:        present,
		StInfo:           studentinfo,
		AllCourceDataOut: programdataout,
	}

	err := tpl.ExecuteTemplate(w, "studentportal.html", students_data)

	if err != nil {
		// http.Redirect(w, r, "/error", http.StatusSeeOther)
		// return

		fmt.Println("Failed to load the student material", err)
	}

}

func ConfirmStudentLogin(w http.ResponseWriter, r *http.Request) {

	// var students_data_acams ACAMS

	//WORK ON STUDENT VALIDAION AND SECURITY CHECK

	tpl = template.Must(template.ParseGlob("templates/*.html"))
	var programdataout []AllCourceData
	var studentinfo StudentInfo

	var present bool
	var setroute string
	r.ParseForm()

	if r.Method == "POST" {

		studentemail := r.FormValue("studentemail")
		password := r.FormValue("studentpassword")

		confirm, studentuuid := ValidateSudent(studentemail, password)
		fmt.Println("Confirmed")

		if confirm {

			// Redirect with query parameters

			http.Redirect(w, r, fmt.Sprintf("/studentportal?student=%s", studentuuid), http.StatusSeeOther)

			return

		} else {
			setroute = "loginerror.html"
		}

	}

	students_data := StudentCourse{
		Available:        present,
		StInfo:           studentinfo,
		AllCourceDataOut: programdataout,
	}

	tpl.ExecuteTemplate(w, setroute, students_data)

}


func checkUserProgramPrepared(student_uuid, program_code string) (bool, error) {
	dbread := dbcode.SqlRead().DB

	defer dbread.Close()

	var exists bool
	query_stmt  := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE student_uuid = ?)", program_code)
	stmt, err := dbread.Prepare(query_stmt)
	if err != nil {
		log.Fatal(err)
	}
	
	err = stmt.QueryRow(student_uuid).Scan(&exists)

	if err != nil {
		exists = false
	}

	return exists, err
}

func checlIfApplied(student_uuid string) []ProgramDataEntry{

	programs, err := GetAllProgramData()
	var programs_not_studing []ProgramDataEntry 
	var present bool
	

	if err != nil {
		fmt.Println("No Programs Exist")
	}

	

	
	for _, program := range programs{
		present, err = checkUserProgramPrepared(student_uuid, program.Code)

		if err != nil{
			fmt.Println("Failed to get user from checkUserExistsPrepared")
		}else{
			fmt.Println("Is present: ",present)
			if !present{
				programs_not_studing = append(programs_not_studing, program)
			}
		}
		
	}

	return programs_not_studing

}


type ProceedStruct struct{
	Programs_Av []ProgramDataEntry
	StInfo  StudentInfo
}

func StudentProcced(w http.ResponseWriter, r *http.Request) {

	
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	student_uuid := r.URL.Query().Get("student_uuid")

	programs_available := checlIfApplied(student_uuid)
	studentdata := GetStudentAllDetails(student_uuid)

	dataout := ProceedStruct{
		Programs_Av: programs_available,
		StInfo: studentdata,
	}


	


	err := tpl.ExecuteTemplate(w, "st_apply_for_more.html", dataout)

	if err != nil {
		// http.Redirect(w, r, "/error", http.StatusSeeOther)
		// return

		fmt.Println("Failed to load the student material", err)
	}
}


func ApplyProceed(w http.ResponseWriter, r *http.Request){

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	student_uuid := r.URL.Query().Get("student_uuid")
	program_code := r.URL.Query().Get("program")

	var setteltemplate string 

	programs := GetStudentPrograms(student_uuid)


	complete := AppendProgramList(program_code, student_uuid,programs)

	getSD := GetStudentAllDetails(student_uuid)
	payment_type := ""
	date_applied :=fmt.Sprintf("%s",time.Now())
	
	program_data := StudentProgramData{
		Student_UUID:   student_uuid,
		First_Name:     getSD.First_Name,
		Last_Name:      getSD.Last_Name,
		Email:          getSD.Email,
		Applied:        true,
		Approved:       false,
		Payment_Method: payment_type,
		Paid:           "pending",
		Completed:      false,
		Date:           date_applied,
	}
	if complete{
		err := CreateProgamData(program_data, payment_type, program_code)
		if err != nil {
			fmt.Println("Failed to apply to new program")
		}else{
			setteltemplate = "progragramapplied"

		}
	}else{
		setteltemplate = "programapplicationfailed"
	}	

	err := tpl.ExecuteTemplate(w, setteltemplate, nil)
	

	if err != nil {
		log.Fatal(err)
	}
			

}


	
	
	