package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"
)

type CourceDataStruct struct {
	UUID             string
	Program_Name     string
	Cource_Name      string
	Cource_Code      string
	Cource_Aseesment string
	Video_List       string
	Module           string
	Book             string
	Exam             bool
}

type ProgramDataOut struct {
	Present      bool
	Program_Name string
	ProgramData  []CourceDataStruct
	Admin        AdminInfo
	Admin_Name string
}

type CourceDataUpdate struct {
	Update bool
	Data   CourceDataStruct
}

func UpdateCourceData(w http.ResponseWriter, r *http.Request) {

	uuid := r.URL.Query().Get("cource_uuid")

	fmt.Println("Hit")

	dbconn := dbcode.SqlRead().DB
	var cource_data CourceDataStruct

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	stmt, err := dbconn.Prepare("select uuid, program_name, cource_name, cource_code,cource_assesment, video_list, module,recomended_book from cource_table where uuid = ?")


	if err != nil {
		err_out := fmt.Errorf("Failed to read from DB, error out ONE: %w", err)
		fmt.Println(err_out)
	}

	defer stmt.Close()

	err = stmt.QueryRow(uuid).Scan(&cource_data.UUID, &cource_data.Program_Name, &cource_data.Cource_Name,  &cource_data.Cource_Code,&cource_data.Cource_Aseesment, &cource_data.Video_List, &cource_data.Module, &cource_data.Book)

	if err != nil {
		err_out := fmt.Errorf("Failed to read from DB, error out TWO: %w", err)

		fmt.Println(err_out)
	}


	fmt.Println("Cource UUID:::::: ", uuid)
	fmt.Println("Cource Data:::::: ", cource_data)


	err_out := tpl.ExecuteTemplate(w, "updatecoursedatanew", cource_data)

	if err_out != nil {
		fmt.Println("Failed to get data for update: ", err_out)
	}

}

func UpdateProgramDetails(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	uuid := r.URL.Query().Get("uuid")

	program_name := r.FormValue("program_name")
	cource_name := r.FormValue("cource_name")

	book_link := r.FormValue("book_link")
	module_link := r.FormValue("module_link")
	video_link := r.FormValue("video_link")
	assesment_link := r.FormValue("assesment_link")

	update := dbcode.SqlRead().DB

	stmt, err := update.Prepare("UPDATE cource_table SET program_name = ?, cource_name = ? , cource_assesment = ?, video_list = ?,module = ?,recomended_book = ? where uuid = ? ")
	if err != nil {
		fmt.Println(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(program_name, cource_name, assesment_link, video_link, module_link, book_link, uuid)

	if err != nil {
		log.Fatal(err)
	}

	cource_data := CourceDataStruct{
		UUID:             uuid,
		Program_Name:     program_name,
		Cource_Name:      cource_name,
		Cource_Aseesment: assesment_link,
		Video_List:       video_link,
		Module:           module_link,
		Book:             book_link,
	}

	err_out := tpl.ExecuteTemplate(w, "cource_data_tr", cource_data)

	if err_out != nil {
		log.Fatal(err)
	}

}

// UPATE THE ALLOW FOR CREATED EXAM

func ExamTrue(uuid string) {

	updateexam := dbcode.SqlRead().DB
	fmt.Println("Exam entered true", uuid)

	stmt, err := updateexam.Prepare("UPDATE cource_table SET exam_file = ?  where uuid = ? ")
	if err != nil {
		fmt.Println("The error is here error: ", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(true, uuid)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Exam exiting true")

}

func GetProgramDetailsSingle(uuid_out string) CourceDataStruct {
	var cource_data_out CourceDataStruct
	get_one := dbcode.SqlRead().DB

	stmt, err := get_one.Prepare("select uuid, program_name, cource_name,cource_code, cource_assesment, video_list,module,recomended_book, exam_file from cource_table where uuid = ?")

	if err != nil {
		fmt.Println("Error One", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(uuid_out).Scan(&cource_data_out.UUID, &cource_data_out.Program_Name, &cource_data_out.Cource_Name, &cource_data_out.Cource_Code,&cource_data_out.Cource_Aseesment, &cource_data_out.Video_List, &cource_data_out.Module, &cource_data_out.Book, &cource_data_out.Exam)

	if err != nil {
		fmt.Println("UUID: ", uuid_out)

		fmt.Println("Error two", err)
	}

	return cource_data_out
}

func GetProgramDetails(program_name string) ([]CourceDataStruct, bool) {
	var cuorce_data_out_list []CourceDataStruct
	var cource_data_out CourceDataStruct

	data_present := true

	get_cource_data := dbcode.SqlRead().DB

	statement, err := get_cource_data.Query("select * from cource_table")

	if err != nil {
		log.Fatal(err)
		data_present = false
	}
	defer statement.Close()

	

	for statement.Next() {
		err := statement.Scan(
			&cource_data_out.UUID, &cource_data_out.Program_Name, &cource_data_out.Cource_Name,&cource_data_out.Cource_Code, &cource_data_out.Cource_Aseesment, &cource_data_out.Video_List, &cource_data_out.Module, &cource_data_out.Book, &cource_data_out.Exam,
		)
		if err != nil {
			log.Fatal(err)
		}

		if cource_data_out.Program_Name == program_name {
			cuorce_data_out_list = append(cuorce_data_out_list, cource_data_out)
		} else {
			continue
		}

	}
	if len(cuorce_data_out_list) < 1 {
		data_present = false

	}

	return cuorce_data_out_list, data_present
}











func ProgramDetails(w http.ResponseWriter, r *http.Request) {

	

	admin_name, err := GetUserName(r)

	code := r.URL.Query().Get("programcode")

	fmt.Println("Program Code: ", code)
	
	// admin_infor := AdminData(admin_id)


	var program_data ProgramDataOut

	result, present := GetProgramDetails(code)

	tpl = template.Must(template.ParseGlob("templates/*.html"))



	if present {

		program_data = ProgramDataOut{
			Present:      true,
			Program_Name: code,
			ProgramData:  result,
			Admin_Name: admin_name,
		}

	} else {
		program_data = ProgramDataOut{
			Present:      false,
			Program_Name: code,
			Admin_Name:        admin_name,
		}
	}


	err = tpl.ExecuteTemplate(w, "A_programedetails.html", program_data)

	if err != nil {
		log.Fatal(err)
	}
}

func CreateCourseData(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	parameter_in := r.URL.Query().Get("parameter")
	program_name := r.URL.Query().Get("program_name")

	var data_out CourceDataUpdate
	fmt.Println("The Parameter in: ", program_name)
	fmt.Println("The Program in: ", program_name)

	var setrout string

	if parameter_in == "update" {
		uuid := r.URL.Query().Get("uuid")
		data_out.Update = true
		data_out.Data = GetProgramDetailsSingle(uuid)
		setrout = "form_update"

	} else if parameter_in == "create" {
		data_out.Update = false
		setrout = "create_cource_data"
		data_out.Data.Program_Name = program_name
	}

	err := tpl.ExecuteTemplate(w, setrout, data_out)

	if err != nil {
		log.Fatal(err)
	}

}

func CloseCreateCourseData(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	fmt.Println("The close div template called")

	err := tpl.ExecuteTemplate(w, "cource_data_close", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func GetCourceMaterial(cource_name, material_name string) string {
	material_out := dbcode.SqlRead().DB
	var link_out string

	if material_name == "video" {
		stmt, err := material_out.Prepare("select video_list from cource_table where cource_name = ?")

		if err != nil {
			fmt.Println("not working")
		}

		defer stmt.Close()

		err = stmt.QueryRow(cource_name).Scan(&link_out)
		if err != nil {
			fmt.Println("failed to get video")
		}
	}

	return link_out

}

func GetLink(cource_name, material_type string) string {
	dbconn := dbcode.SqlRead().DB

	fmt.Println("THE MATERIAL TYPE", material_type)

	var link_name string
	query_statement := fmt.Sprintf("select %s from cource_table where cource_name = ?", material_type)

	fmt.Println("The Query string: ", query_statement, "Cource Name: ", cource_name)
	stmt, err := dbconn.Prepare(query_statement)

	if err != nil {
		fmt.Println("Prepare statment failed ", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(cource_name).Scan(&link_name)

	if err != nil {
		fmt.Println("Query Failed: ", err)
	}

	return link_name
}

func GetStudyMaterial(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Route Has Been Hit")
	material_type := r.URL.Query().Get("material_type")
	cource_name := r.URL.Query().Get("cource_name")

	link_out := GetLink(cource_name, material_type)

	http.Redirect(w, r, link_out, http.StatusSeeOther)
}
func CleanVideoLinks(links string) string {
	var clean_list string

	list_out := strings.Split(links, "_")

	if len(list_out) <= 1 {

	}

	return clean_list
}

func AddCourceData(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	var course_code string
	var cource_name string
	var book_link string
	var module_link string 
	var video_link string 
	var assesment_link string 
	
	program_name := r.FormValue("program_name")
	cource_name = r.FormValue("cource_name")
	course_code = SanitizeCookieValue(r.FormValue("course_code"))
	book_link = r.FormValue("book_link")
	module_link = r.FormValue("module_link")
	video_link = r.FormValue("video_link")
	assesment_link = r.FormValue("assesment_link")

	create_uuid := encription.Generateuudi()
	exam := "false"



	create_cource := dbcode.SqlRead().DB
	statment, err := create_cource.Prepare("insert into cource_table(uuid, program_name, cource_name, cource_code,cource_assesment, video_list,module,recomended_book, exam_file) values(?,?,?,?,?,?,?,?,?)")

	if err != nil {
		log.Fatal(err)
	}

	defer statment.Close()

	_, err = statment.Exec(create_uuid, program_name, cource_name, course_code,assesment_link, video_link, module_link, book_link, exam)

	if err != nil {
		log.Fatal(err)
	}

	data_out := CourceDataStruct{
		UUID:             create_uuid,
		Program_Name:     program_name,
		Cource_Name:      cource_name,
		Cource_Code:      course_code,
		Book:             book_link,
		Module:           module_link,
		Video_List:       video_link,
		Cource_Aseesment: assesment_link,
	}

	err = CreateNewCourseTable(course_code)


	if err != nil {
		fmt.Println("Failed to create new course table")

	}

	err = CreateCourseMaterial(program_name, cource_name, course_code)

	if err != nil {
		fmt.Println("Failed to add course primary data to CourseNames: ", err)
	}

	exam_details := Exam_Details{
		Cource_UUID:  create_uuid,
		Program_Name: program_name,
		Cource_Name:  cource_name,
		Cource_Code:  course_code,
		Duration:     "0",
		Total_Marks: "0",				
	
	}

	err = Create_Exam_Details(exam_details)
	if err != nil{
		fmt.Println("Create Exam Details in \"program details\" returns error: ", err)
	}

	err = tpl.ExecuteTemplate(w, "cource_data_tr_two", data_out)

	if err != nil {
		log.Fatal(err)
	}

}


