package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kamukwamba/oerisuniversity/dbcode"

	"github.com/kamukwamba/oerisuniversity/encription"
)

type MessageBody struct {
	UUID         string
	Message      string
	Sender_UUID  string
	Sender_Name  string
	Sender       bool
	Seen_Student bool
	Seen_Admin   bool
	Date         string
}

type MessageLading struct {
	StudentUUID string
	MsgBody     []MessageBody
	StInfo      StudentInfo
}

func GetMessages(uuid string) []MessageBody {
	var messages_out_list []MessageBody

	dbread := dbcode.SqlRead()
	stmt, err := dbread.DB.Query("select uuid, sender_uuid,sender_name, sender,message,seen_student,seen_admin,date from messages")

	uuidout := uuid
	if err != nil {

		fmt.Println("No messages")

	}
	defer stmt.Close()
	var message_out MessageBody

	for stmt.Next() {
		err = stmt.Scan(&message_out.UUID, &message_out.Sender_UUID, &message_out.Sender_Name, &message_out.Sender, &message_out.Message, &message_out.Seen_Student, &message_out.Seen_Admin, &message_out.Date)

		if err != nil {
			log.Fatal(err)
		}

		if message_out.Sender_UUID == uuidout {
			messages_out_list = append(messages_out_list, message_out)
		}

	}

	return messages_out_list
}

func arrayToString(arr []string) string {

	// seperating string elements with -
	return strings.Join([]string(arr), ",")
}

func SendMsg(w http.ResponseWriter, r *http.Request) {

	uuid := encription.Generateuudi()
	write_msg := r.PostFormValue("message_content")

	student_uuid := r.URL.Query().Get("student_uuid")
	from := r.URL.Query().Get("from")

	student_name_get := GetStudentAllDetails(student_uuid)
	sender_name := fmt.Sprintf("%s %s", student_name_get.First_Name, student_name_get.Last_Name)

	fmt.Println(sender_name)

	msgyear := strconv.Itoa(time.Now().Year())
	msgmonth := time.Now().Month().String()
	msgday := strconv.Itoa(time.Now().Day())
	var seen_student bool
	var seen_admin bool
	var sender bool

	if from == "student" {
		sender = true
		seen_student = true
		seen_admin = false
	} else if from == "admin" {
		sender = false
		seen_student = false
		seen_admin = true

	}

	datein := []string{msgyear, msgmonth, msgday}
	dateout := arrayToString(datein)

	dbcode := dbcode.SqlRead()

	sendmsg, err := dbcode.DB.Begin()

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := sendmsg.Prepare("insert into messages (uuid, sender_uuid,sender_name, sender,message,seen_student,seen_admin,date) values(?,?,?,?,?,?,?,?)")

	if err != nil {
		ErrorPrintOut("studentpersonal", "sendmessage", err.Error())
	}

	defer stmt.Close()
	_, err = stmt.Exec(uuid, student_uuid, sender_name, sender, write_msg, seen_student, seen_admin, dateout)

	if err != nil {
		log.Fatal(err)
	}

	err = sendmsg.Commit()
	if err != nil {
		log.Fatal(err)
	}

	data := MessageBody{
		UUID:        uuid,
		Sender_UUID: student_uuid,
		Sender:      sender,
		Message:     write_msg,
		Date:        dateout,
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	errtmpl := tpl.ExecuteTemplate(w, "right", data)

	if errtmpl != nil {
		log.Fatal(errtmpl)
	}

}

func ProgramCompleted(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	user_uuid := r.URL.Query().Get("user_uuid")
	program := r.URL.Query().Get("program")

	switch program {
	case "acams":
		if UpdateProgram(user_uuid, program) {
			fmt.Println("Update was succesfull")

		}

	case "acms":
		if UpdateProgram(user_uuid, program) {
			fmt.Println("Update was succesfull")
		}
	case "adms":
		if UpdateProgram(user_uuid, program) {
			fmt.Println("Update was succesfull")
		}
	case "abdms":
		if UpdateProgram(user_uuid, program) {
			fmt.Println("Update was succesfull")
		}

	}

	// datastring := fmt.Sprintf("The querris are %s ", dataout)
	// fmt.Fprint(w, datastring)

	// keys, ok := r.URL.Query()["id"]

	// if ok {
	// 	fmt.Println(keys)
	// }

	err := tpl.ExecuteTemplate(w, "completed", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func Messages(w http.ResponseWriter, r *http.Request) {

	student_uuid := r.PathValue("id")

	message_out := GetMessages(student_uuid)
	studentdata := GetStudentAllDetails(student_uuid)

	data := MessageLading{
		StudentUUID: student_uuid,
		MsgBody:     message_out,
		StInfo:      studentdata,
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "studentmessages.html", data)

	if err != nil {
		log.Fatal("ERROR=== ", err, " ===END")
	}

}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	var setroute string
	var settext string

	r.ParseForm()
	uuid := r.URL.Query().Get("uuid")

	fmt.Println(uuid)

	dbread := dbcode.SqlRead()
	stmt, err := dbread.DB.Prepare("select uuid, student_uuid, email, password from studentcridentials where student_uuid = ?")

	if err != nil {

		fmt.Println("First err")
		log.Fatal(err)
	}

	defer stmt.Close()

	var uuid_out string
	var student_uuid string
	var email_out string
	var password string

	err = stmt.QueryRow(uuid).Scan(&uuid_out, &student_uuid, &email_out, &password)

	if err != nil {
		fmt.Println("Second err", err)
		// log.Fatal(err)
	}

	old_password := r.FormValue("old_password")
	new_password := r.FormValue("new_password")
	confirm_password := r.FormValue("confirm_password")

	if old_password != password {
		setroute = "updateresponce"
		settext = "old password is in correct"
	} else {
		if new_password != confirm_password {
			setroute = "checkpassword"
			settext = "Passwords do noot match!!!"
		} else {
			dbread := dbcode.SqlRead().DB
			statement, err := dbread.Prepare("update studentcridentials SET password = ? WHERE uuid = ?")

			if err != nil {
				error_text := fmt.Sprintf("line 44 error from update prepare:: %s", err)
				ErrorPrintOut("studentportal", "ApplyForCource", error_text)
			}
			defer statement.Close()

			_, errup := statement.Exec(confirm_password, uuid_out)

			if errup != nil {
				error_text := fmt.Sprintf("line 50 error from update prepare:: %s", errup)
				ErrorPrintOut("studentportal", "ChangeStudentPassword", error_text)
			}
			setroute = "updateresponce"
			settext = "password updated"
		}
	}

	err = tpl.ExecuteTemplate(w, setroute, settext)

	if err != nil {
		log.Fatal(err)

	}

}

func CloseUpdateData(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	fmt.Println("Close Update Passwsord")

	err := tpl.ExecuteTemplate(w, "closeupdate", nil)

	if err != nil {
		log.Fatal(err)

	}
}

func ChangeStudentPassword(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	student_uuid := r.URL.Query().Get("student")

	fmt.Println("Student UUID", student_uuid)

	err := tpl.ExecuteTemplate(w, "changepassword", student_uuid)

	if err != nil {
		log.Fatal(err)

	}
}

func StudentSettings(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	uuid := r.URL.Query().Get("uuid")

	data_out := GetStudentAllDetails(uuid)

	students_data := StudentCourse{
		StInfo: data_out,
	}

	err := tpl.ExecuteTemplate(w, "studentprofile.html", students_data)

	if err != nil {
		log.Fatal(err)

	}

}

func StudentLogOut(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	uuid := r.URL.Query().Get("uuid")

	data_out := GetStudentAllDetails(uuid)

	err := tpl.ExecuteTemplate(w, "studentdata", data_out)

	if err != nil {
		log.Fatal(err)

	}
}

func ContactInstitution(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	videopath := r.PathValue("id")
	fmt.Println(videopath)

	err := tpl.ExecuteTemplate(w, "messagesstudent", nil)

	if err != nil {
		log.Fatal(err)

	}

}

func StudentProfilePortal(w http.ResponseWriter, r *http.Request) {
	studentuuid := r.PathValue("id")
	var programdataout []AllCourceData
	var studentinfo StudentInfo

	var present bool
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	studentprogramlist := GetStudentPrograms(studentuuid)
	fmt.Println("student programs", studentprogramlist)
	programdataout, present = GetStudentProgramData(studentprogramlist, studentuuid)
	studentinfo = GetStudentAllDetails(studentuuid)

	students_data := StudentCourse{
		Available:        present,
		StInfo:           studentinfo,
		AllCourceDataOut: programdataout,
	}

	fmt.Println("WORKING")

	tpl.ExecuteTemplate(w, "studentportal.html", students_data)
}

type VideoStruct struct {
	Video string
}

type VideoDisplay struct {
	Present     bool
	Cource_Name string
	MainVideo   VideoStruct
	VideoList   []VideoStruct
}

func DeleteStudentExam(cource_uuid string) {

}

func CloseAssesmentDiv(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	fmt.Println("close button pressed")

	err := tpl.ExecuteTemplate(w, "empty_div", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func RecordStudentMarks(student_answers Answer_Out) bool {
	dbconn := dbcode.SqlRead().DB

	saved := true

	table_name := CleanStudentUUID(student_answers.Student_UUID)

	uuid := encription.Generateuudi()
	prepare_statment := fmt.Sprintf("insert into %s (uuid, cource_uuid, cource_name,student_uuid, question_number, question, attempt_number,answer) values(?,?,?,?,?,?,?,?)", table_name)
	stmt, err := dbconn.Prepare(prepare_statment)

	if err != nil {
		fmt.Println("Failed to prepare create statment:  ", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(uuid, student_answers.Cource_UUID, student_answers.Cource_Name, student_answers.Student_UUID, student_answers.Question_Number, student_answers.Question, student_answers.Attemp_Number, student_answers.Answer)

	if err != nil {
		fmt.Println("Failed to save exam answers: ", err)
	}

	return saved
}

func CleanStudentUUID(uuid string) string {

	new_uuid := strings.Split(uuid, "-")

	new_string := "st"

	for _, item := range new_uuid {
		new_string = new_string + item
	}

	return new_string
}

func MakeStudentExamTable(student_uuid string) {

	dbconn := dbcode.SqlRead().DB

	new_string := CleanStudentUUID(student_uuid)

	fmt.Println("The string name:::", new_string, "::")

	create_table := fmt.Sprintf(`create table if not exists %s(
		uuid blob not null,
		cource_uuid text,
		cource_name text,
		student_uuid text,
		question_number text,
		question text,
		attempt_number text,
		answer text)`, new_string)

	stmt, err := dbconn.Prepare(create_table)

	if err != nil {
		fmt.Println("Failed to create table error_one::: ", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec()

	if err != nil {
		fmt.Println("Failed to load table error_two: ", err)
	}

	if err != nil {
		fmt.Println("Failed to create student exam answer table error_three: ", err)
	}

}

func WatcVideo(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))
	var videos_present bool

	cource_name := r.URL.Query().Get("cource_name")
	fmt.Println(cource_name)

	video_list := GetCourceMaterial(cource_name, "video")

	video_list_out := strings.Split(video_list, ",")

	if len(video_list) >= 1 {
		videos_present = true

	} else {
		videos_present = false
	}

	var video_link VideoStruct

	var all_videos []VideoStruct
	var AllVidoes VideoDisplay

	if videos_present {
		for _, item := range video_list_out {
			video_link = VideoStruct{
				Video: item,
			}

			all_videos = append(all_videos, video_link)

		}

		AllVidoes = VideoDisplay{
			Present:     true,
			Cource_Name: cource_name,
			MainVideo:   all_videos[0],
			VideoList:   all_videos,
		}

	} else {
		AllVidoes = VideoDisplay{
			Present:     false,
			Cource_Name: cource_name,
		}
	}

	err := tpl.ExecuteTemplate(w, "videos.html", AllVidoes)

	if err != nil {
		log.Fatal(err)

	}

}
