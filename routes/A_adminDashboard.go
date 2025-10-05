package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kamukwamba/oerisuniversity/dbcode"
)
type StudentDataPage struct {
	Admin_Name string
	ProgramData []ProgramStruct
}



type AdminLogData struct {
	Email    string
	Password string
}
type AdminInfo struct {
	ID       string
	Name     string
	Email    string
	Password string
}

// Marked For Removal
type AdminPage struct {
	Admin      AdminInfo
	AcamsData  []ProgramStruct
	Admin_Name string
}

// Marked For Removal
type DashData struct {
	ACAMSTotal int
	Admin      AdminInfo
}

type StudentProgramList struct {
	Program_Name string
}

func DeleteStudentInfo(uuid string) bool {

	var examtable = fmt.Sprintf("st%s", uuid)
	var deleted = true
	var tables = []string{"studentdata", "studentcridentials"}
	dbconn := dbcode.SqlRead().DB

	for item := range tables {

		var preparestmt = fmt.Sprintf("DELETE FROM %s WHERE uuid = ?", item)

		stmt, err := dbconn.Prepare(preparestmt)

		if err != nil {
			fmt.Println("failed to prepare delete: ", err)
			deleted = false

		}

		defer stmt.Close()

		_, err = stmt.Exec(uuid)

		if err != nil {
			fmt.Println("FAILED TO DELETE: ", err)
			deleted = false
		}

	}

	query := fmt.Sprintf("DROP TABLE IF EXISTS %s;", examtable)

	// Execute the SQL query
	_, err := dbconn.Exec(query)
	if err != nil {
		fmt.Errorf("error dropping table %q: %w", examtable, err)
		deleted = false
	}

	return deleted

}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {

	uuid := r.URL.Query().Get(".uuid")

	success := DeleteStudentInfo(uuid)

	if success {
		w.WriteHeader(http.StatusOK)

	} else {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return

	}

}

func AdminData(uuid string) AdminInfo {
	var admin_info AdminInfo
	dbconn := dbcode.SqlRead().DB

	fmt.Println(uuid)

	stmt, err := dbconn.Prepare("select uuid, first_name, last_name, email from admin where uuid = ?")

	if err != nil {
		fmt.Println("Failed to load admin prepare", err)
	}

	defer stmt.Close()

	var first_name string
	var last_name string
	var id string
	var email string

	err = stmt.QueryRow(uuid).Scan(&id, &first_name, &last_name, &email)

	name := fmt.Sprintf("%s %s", first_name, last_name)

	admin_info = AdminInfo{
		ID:    id,
		Name:  name,
		Email: email,
	}

	if err != nil {
		fmt.Println("Failed to QueryRow: ", err)
	}

	return admin_info
}

func AdminAuth(data AdminLogData, dataList []dbcode.AdminInfo) (bool, AdminInfo) {

	var result bool
	var admin_data AdminInfo

	for _, admin_info := range dataList {
		name_out := fmt.Sprintf("%s %s", admin_info.First_Name, admin_info.Last_Name)

		id := admin_info.ID
		name := name_out
		email := admin_info.Email
		password := admin_info.Password

		matchPassword := CheckPassword(password, data.Password)

		if matchPassword == true && data.Email == email {

			admin_data = AdminInfo{
				ID:       id,
				Name:     name,
				Email:    email,
				Password: password,
			}

			result = true

			break
		}
	}
	return result, admin_data
}

func StudentACAMSData(student_uuid string) {
	dbread := dbcode.SqlRead()

	stmt, err := dbread.DB.Prepare("select program_list from studentprogramlist where student_uuid = ?")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var program_list string

	err = stmt.QueryRow(student_uuid).Scan(&program_list)

	if err != nil {
		fmt.Println("FAILED TO GET STUDENT PROGRAM LIST")
	}

}







func GetAllStudentsData() []ProgramStruct {
	dbread := dbcode.SqlRead()
	var conuter int
	var datalist []ProgramStruct

	getallprograms, err := GetAllProgramData()


	if err != nil {
		fmt.Println("Failed to get student Program Data: ", err)
	}


	for _, item := range getallprograms{
		fmt.Println(item.Code)
		var cource_data ProgramStruct

		query_str := fmt.Sprintf("select * from %s", item.Code)
		rows, err := dbread.DB.Query(query_str)
		if err != nil {
			fmt.Printf("Failed to get student data from %s ", query_str)
		}
		defer rows.Close()

		for rows.Next() {
			conuter += 1

			var uuid string
			var student_uuid string
			var program_name string
			var first_name string
			var last_name string
			var email string
			var applied bool
			var approved bool
			var payment_method string
			var paid string
			var completed bool
			var date string
			err = rows.Scan(&uuid, &student_uuid, &program_name, &first_name, &last_name, &email, &applied, &approved, &payment_method, &paid, &completed, &date)

			cource_data = ProgramStruct{
				UUID:           uuid,
				Student_UUID:   student_uuid,
				Program_Name:   program_name,
				First_Name:     first_name,
				Last_Name:      last_name,
				Email:          email,
				Applied:        applied,
				Approved:       approved,
				Payment_Method: payment_method,
				Paid:           paid,
				Completed:      completed,
				Date:           date,
			}

			if err != nil {
				fmt.Println("Check the scan for student data admindashboard")
				log.Fatal(err)
			}

			datalist = append(datalist, cource_data)

		}
	}

	

	return datalist
}

func ACAMSCount() int {
	dbread := dbcode.SqlRead()
	var counter int

	rows, err := dbread.DB.Query("select * from  acams")
	if err != nil {
		fmt.Println("Failed to get acams student data")
	}
	defer rows.Close()

	for rows.Next() {
		counter += 1

	}

	return counter
}

func GetStudentAllDetails(uuid string) StudentInfo {
	dbread := dbcode.SqlRead()

	stmt, err := dbread.DB.Prepare("select  uuid, first_name, last_name, phone, email, date_of_birth,gender,marital_status,country,eduction_background,program,high_scholl_confirmation,grammer_comprihention,waiver,number_of_children,school_atteneded,major_studied,degree_obtained,current_occupetion,field_interested_in,mps_techqnique_Practiced,previouse_experince,purpose_of_enrollment,use_of_degree,reason_for_choice,method_of_incounter from studentdata where uuid = ?")

	if err != nil {
		fmt.Println("the student uuid", uuid)

	}
	defer stmt.Close()

	var data StudentInfo

	err = stmt.QueryRow(uuid).Scan(&data.UUID, &data.First_Name, &data.Last_Name, &data.Phone, &data.Email, &data.Date_Of_Birth, &data.Gender, &data.Marital_Status, &data.Country, &data.Education_Background, &data.Program, &data.High_School, &data.Grammer_Confirmation, &data.Waiver, &data.Children, &data.School_Attended, &data.Major_In, &data.Degree_Obtained, &data.Current_Occupation, &data.Field_Interested, &data.Prio_Techniques, &data.Previouse_Experience, &data.Purpose_Of_Enrollment, &data.Use_Of_Knowledge, &data.Reason_For_Choice, &data.Method_Of_Encounter)

	if err != nil {
		log.Fatal(err)

	}

	// fmt.Println("student info: ", data) //PRINT STUDENT DATA

	return data
}

//ROUTER CODE

//Approve Cources

func UpdateCource(uuid, cource_name string) bool {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	cource_approved := true

	dbread := dbcode.SqlRead().DB

	dbupdate_statement := fmt.Sprintf(`UPDATE %s SET approved = ? WHERE student_uuid = ? `, cource_name)

	statement, err := dbread.Prepare(dbupdate_statement)

	if err != nil {
		error_text := fmt.Sprintf("line 44 error from update prepare:: %s", err)
		ErrorPrintOut("studentportal", "ApplyForCource", error_text)

		cource_approved = false
	}
	defer statement.Close()

	_, errup := statement.Exec(true, uuid)

	if errup != nil {
		error_text := fmt.Sprintf("line 50 error from update prepare:: %s", errup)
		ErrorPrintOut("studentportal", "ApplyForCource", error_text)
		cource_approved = false
	}
	return cource_approved

}

func DropTable(cource_code string) error{
	dbread := dbcode.SqlRead().DB

    defer dbread.Close()


    fmt.Println("THE TABLE TO DROP: ", cource_code)
    
    dropTableSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", cource_code)
    
    _, err := dbread.Exec(dropTableSQL)
    if err != nil {
        log.Fatal(err)
        return err
    }

    return nil
    

}



func DeleteCourceNames(cource_code string) error{
	dbconn := dbcode.SqlRead().DB

	stmt, err := dbconn.Prepare("Delete from CourseNames where courseCode = ?")

	if err != nil {
		fmt.Println("Failed to delete from cource names error: ", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(cource_code)
	if err != nil {
		fmt.Println("Failed to delete from cource names error: ", err)
		return err
	}

	return nil
}

func DeleteCource(w http.ResponseWriter, r *http.Request) {

	cource_code := r.URL.Query().Get("cource_name")
	cource_uuid := r.URL.Query().Get("cource_uuid")

	dbconn := dbcode.SqlRead().DB

	stmt, err := dbconn.Prepare("Delete from cource_table where uuid = ?")

	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	defer stmt.Close()

	_, err = stmt.Exec(cource_uuid)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)

	}

	if !DeleteAllQuestions(cource_uuid) {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		fmt.Println("Failed to delete questions")
	}

	if DropTable(cource_code) != nil {
		fmt.Println("Failed to drop cource tables")
	}

	if DeleteCourceNames(cource_code) != nil{
		fmt.Println("Failed to delete cource names")
	}


	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, "empty_tr", nil)

	if err != nil {
		fmt.Println("Failed to delete")
	}

}

func ApproveCourceUpdate(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	uuid := r.URL.Query().Get("user_uuid")
	cource_name := r.URL.Query().Get("cource_name")

	var setroute string
	result := UpdateCource(uuid, cource_name)

	if result {

		setroute = "cource_approved_admin"
	} else {
		setroute = "approval_error"
	}
	// dbread := dbcode.SqlRead().DB

	err := tpl.ExecuteTemplate(w, setroute, nil)

	if err != nil {
		log.Fatal(err)
	}

}


type AllStudentProgramData struct{
	Program_Name string
	Course_List string
}

func GetStudentProgramDataAdmin(programlist []string, students_uuid string) ([]AllCourceData, bool){

	//Get information 


	//GET ALL THE COURCES NAME ASSOCIATED WITH  THE PROGRAM
	//RUN THE LOOP TO GET THE DATA FROM THE PROGRAM TABLE
	//WITHIN THE PROGRAM LOOP RUN A LOOP TO WITH THE COURCE ASSOCIATED WITH THE PROGRAM


	var programdata ProgramStruct
	var courcedata []CourceStruct
	var allcourcedataout AllCourceData
	var allcourcedataoutlist []AllCourceData
	

	var programsavailable ProgramAvailable
	var available []bool

	

	for _, program := range programlist {
		is_present, dataout, _ := GetProgramAdmin(students_uuid, "one", program)

		fmt.Println("The Program: ", program)
		fmt.Println("The dataout: ", dataout)

		fmt.Println("IS Present: ", is_present)


		if is_present {

			programdata = ProgramStruct{
				UUID: dataout.UUID,
				Student_UUID: dataout.Student_UUID,
				Program_Name: dataout.Program_Name,
				First_Name: dataout.First_Name,
				Last_Name: dataout.Last_Name,
				Email: dataout.Email,
				Payment_Method: dataout.Payment_Method,
				Paid: dataout.Paid,
				Approved: dataout.Approved,
				Applied: dataout.Applied,
				Completed: dataout.Completed,
				Date: dataout.Date,
			}

			

			courcedata = GetFromProgramCourcesAdmin(students_uuid, program)

			allcourcedataout = AllCourceData{
				ProgramStruct: programdata,
				Cource_Struct: courcedata,
			}


			allcourcedataoutlist = append(allcourcedataoutlist, allcourcedataout)

			available = append(available, true)

			


		
		} else {
			available = append(available, false)
		}

	}


	programsavailable.Available = available[0]

	
	return allcourcedataoutlist, programsavailable.Available

}

func ReadMessageAdmin(w http.ResponseWriter, r *http.Request) {

	student_uuid := r.URL.Query().Get("student_uuid")
	message_seen := r.URL.Query().Get("message_seen")

	switch message_seen {
	case "not":
		UpdateMessages(student_uuid)

	}

	messages := ReadMessage(student_uuid)

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "readMessagesAdmin", messages)

	if err != nil {
		log.Fatal(err)
	}

}

func StudentProfileData(w http.ResponseWriter, r *http.Request) {

	studentuuid := r.URL.Query().Get("student_uuid")

	user_name, err := GetUserName(r)
	
	studentdataout := GetStudentAllDetails(studentuuid)
	listout := GetStudentPrograms(studentuuid)

	programdataout, present := GetStudentProgramDataAdmin(listout, studentuuid)

	studentdataoutadmin := StudentCourse{
		Available:        	present,
		StInfo:           	studentdataout,
		AllCourceDataOut: 	programdataout,
		Admin_Name: 		user_name,
	}


	
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, "studentdetailstemplate.html", studentdataoutadmin)

	if err != nil {
		log.Fatal(err)
	}

}


// func CloseAdmintDiv(w http.ResponseWriter, r *http.Request) {
// 	tpl = template.Must(template.ParseGlob("templates/*.html"))

// 	err := tpl.ExecuteTemplate(w, "student_cource_assesment", nil)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

func StudentData(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	student_data := GetAllStudentsData()

	user_name, err := GetUserName(r)

	if err != nil {
		fmt.Println("Failed to get cookie")
	}

	data_out := StudentDataPage{
		ProgramData:  student_data,
		Admin_Name: user_name,
	}

	err = tpl.ExecuteTemplate(w, "A_studentdataadmin.html", data_out)

	if err != nil {
		log.Fatal(err)
	}
}

type CreatNews struct {
	Admin      AdminInfo
	AllNews    []NewsStruct
	Admin_Name string
}

func AdminNews(w http.ResponseWriter, r *http.Request) {

	_, all := ReadNews("o", "many")

	user_name, err := GetUserName(r)

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	display_data := CreatNews{
		Admin_Name: user_name,
		AllNews:    all,
	}

	err = tpl.ExecuteTemplate(w, "NewsAdminCreate.html", display_data)

	if err != nil {
		log.Fatal(err)
	}
}

func DeleteMessageRouter(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	deleted := DeleteMessage(uuid)
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	if deleted {
		fmt.Println("Deleted")
	}
	err := tpl.ExecuteTemplate(w, "deleted_replacement", nil)

	if err != nil {
		log.Fatal(err)
	}

}

func ApproveProgram(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	user_uuid := r.URL.Query().Get("user_uuid")
	program := r.URL.Query().Get("program")

	if Update(user_uuid, program) {
		fmt.Println("Update succesfull")
	} else {
		fmt.Println("Failed to Update Program")
	}

	err := tpl.ExecuteTemplate(w, "approved", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func AdminMessagesGet(w http.ResponseWriter, r *http.Request) {

}

func AdminMessagesSend(w http.ResponseWriter, r *http.Request) {

}
