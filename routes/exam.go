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

type ExamDetails struct {
	UUID          string
	Program_Name  string
	Cource_Name   string
	Cource_Code   string
	Exam_Duration int
	Total_Marks   string
}

type ExamQuestons struct {
	Section   string
	Questions []string
}

type ExamStruct struct {
	UUID         string
	Program_Name string
	Cource_Name  string
	Questions    []ExamQuestons
}

type Displayed struct {
	Section        string
	Question       string
	Quesion_number int
}

type QuestionStruct struct {
	Section         string
	Question        string
	Question_Number string
	Answer          string
}

type QuestionsEntered struct {
	UUID            string
	Program_Name    string
	Question_Number string
	Cource_Name     string
	Question_Count  int
	Section_A       []QuestionStruct
	Section_B       []QuestionStruct
	Time            int
	Cource_Code     string
	Total_Marks     string
	Date            string
}

type Questions_Out struct {
	Section_A []QuestionStruct
	Section_B []QuestionStruct
}

type Answer_Out struct {
	UUID            string
	Cource_UUID     string
	Student_UUID    string
	Cource_Name     string
	Qustion_UUID    string
	Question_Number string
	Attemp_Number   string
	Question        string
	Answer          string
}

type Questions_Construct struct {
	UUID            string
	Section         string
	Cource_UUID     string
	Cource_Name     string
	Question        string
	Answer          string
	Question_Number int
}

type ExamOut struct {
		Present     bool
		Cource_Data CourceDataStruct
		ExamData    []Questions_Construct
		ExamDetails ExamDetails
		CourceNameFm string
		Admin_Name string
	}

type CreateExamResponse struct {
	Details_Messages string
	Questions        Questions_Construct
}

type DisplayExam struct {
	AlreadyTaken      bool
	Student_UUID      string
	Qusetions_Present bool
	Cource_UUIDOut    string
	Attempt_Number    int ``
	ExamData          ExamDetails
	Exam_Questions    []Question_Structure
	Cource_Name_Two   string
	Writen_UUID       string
	Attempt_Out       string
}

type ExamStructOut struct {
	Student_UUID   string
	Program_Name   string
	Cource_Name    string
	Cource_Code    string
	Time           string
	Attempt_Number string
	Question_List  []QuestionStruct
}

type ExamComplete struct {
	Student_UUID string
	Program_Name string
	Cource_Name  string
	Cource_Code  string
	Grade        string
}

type ExamTakenStruct struct {
	UUID            string
	Cource_Name     string
	Student_UUID    string
	Attemp_Number   string
	First_Attempted string
	Open_Period     string
	Answers         string
	Grade           string
	Date            string
	Completed       string
	Comment         string
	Passed          string
	Gradding    string
}

type Exam_Details struct {
	UUID         string
	Cource_UUID  string
	Program_Name string
	Cource_Name  string
	Cource_Code  string
	Duration     string
	Total_Marks  string
}

type Question_Structure struct {
	Section         string
	Question_UUID   string
	Question_Number string
	Question        string
}

type Complition struct {
	Logo    string
	Message string
}

type TakeExamStruct struct {
	Taken         bool
	ExamTaken     QuestionsEntered
	TakensResults ExamTakenStruct
}


type QuestionData struct {
	Question_UUID   string
	Question        string
	Question_Number string
}

func CounterFunc(str string) int {
	var counter int

	counter = counter + 1

	return counter

}

func ToUpperCase(str string) string {

	return strings.ToUpper(str)
}

func Clean(str string) string {

	str_out := strings.Split(str, "_")

	var join_string string
	
	for _, item :=  range str_out{
		join_string = join_string + " " + item
	}


	capitalised := ToUpperCase(join_string)

	return capitalised
}

func Listify(question_a, question_b string) ([]QuestionStruct, []QuestionStruct) {

	var question_list_a []QuestionStruct //section a questiions
	var question_list_b []QuestionStruct

	question_a = strings.Trim(question_a, "[ ]")
	question_b = strings.Trim(question_b, "[ ]")

	section_A := strings.Split(question_a, "}")
	section_B := strings.Split(question_b, "}")

	fmt.Println(question_a, question_b)

	if len(question_a) > 0 {
		var count_out int
		for _, item := range section_A {

			count_out++

			if count_out != len(section_A) {
				the_result := strings.Trim(item, "{")
				the_result_out := strings.Split(the_result, ":")
				questions_out := QuestionStruct{
					Section:  "A",
					Question: the_result_out[0],
					Answer:   the_result_out[1],
				}

				question_list_a = append(question_list_a, questions_out)

				fmt.Println(question_list_a)

			} else {
				fmt.Println("Out of bounds")
			}

		}
	}

	//Note This code is for presenting the code
	if len(question_b) > 0 {

		var count_out int
		for _, item := range section_B {

			if len(item) > 0 {

				the_result := strings.Trim(item, "{ }")
				fmt.Println(the_result)
				count_out++

				if count_out != len(section_B) {
					the_result := strings.Trim(item, "{")
					questions_out := QuestionStruct{
						Section:  "B",
						Question: the_result,
					}

					question_list_b = append(question_list_b, questions_out)

				} else {
					fmt.Println("Out of bounds")
				}
			}

		}

	}

	return question_list_a, question_list_b
}

func ReadQuestions(cource_name string)([]Questions_Construct, bool){
	questions_present := true
	get_exam := dbcode.SqlRead().DB
	var queOut Questions_Construct
	var question_structure_list []Questions_Construct

	stmt, err := get_exam.Query("select uuid, section, cource_uuid,cource_name,question, answer from exam_questions where cource_name = ?", cource_name)

	if err != nil {
		log.Fatal("Failed to get exam questions: ", err)
		questions_present = false
	}

	

	defer stmt.Close()

	for stmt.Next() {
		err = stmt.Scan(&queOut.UUID, &queOut.Section,&queOut.Cource_UUID ,&queOut.Cource_Name, &queOut.Question,&queOut.Answer)

		if err != nil {
			fmt.Println("failed to query row: ", err)
			break
		} else {
				question_structure_list = append(question_structure_list, queOut)
		}
	}

	if err = stmt.Err(); err != nil {
		fmt.Println("stmt scan failed, error out: ", err)
		questions_present = false
	}

	if err != nil {
		log.Fatal(err)
	}

	return question_structure_list, questions_present

}

func Read_Exam(cource_name string) ([]Question_Structure, string, bool) {

	questions_present := true
	get_exam := dbcode.SqlRead().DB
	var question_structure Question_Structure
	var question_structure_list []Question_Structure

	stmt, err := get_exam.Query("select uuid, section, cource_uuid,cource_name,question, question_number from exam_questions where cource_name = ?", cource_name)

	if err != nil {
		log.Fatal("Failed to get exam questions: ", err)
		questions_present = false
	}

	var uuid string
	var section string
	var cource_uuid string
	var cource_name_out string
	var question string
	var question_number string

	defer stmt.Close()

	for stmt.Next() {
		err = stmt.Scan(&uuid, &section, &cource_uuid, &cource_name_out, &question, &question_number)

		if err != nil {
			fmt.Println("failed to query row: ", err)
			break
		} else {

			if section == "A" {
				question_structure = Question_Structure{
					Section:         section,
					Question_UUID:   uuid,
					Question_Number: question_number,
					Question:        question,
				}

				question_structure_list = append(question_structure_list, question_structure)

			} else if section == "B" {
				question_structure = Question_Structure{
					Section:         section,
					Question_UUID:   uuid,
					Question_Number: question_number,
					Question:        question,
				}
				question_structure_list = append(question_structure_list, question_structure)

			}
		}
	}

	if err = stmt.Err(); err != nil {
		fmt.Println("stmt scan failed, error out: ", err)
	}

	if err != nil {
		log.Fatal(err)
	}

	return question_structure_list, cource_uuid, questions_present

}

func Update_Exam_Details(details Exam_Details) error{
	
	
	
	
	dbupdate := dbcode.SqlRead().DB
	stmt, err := dbupdate.Prepare("UPDATE exam_details SET duration = ?, total_marks = ? where cource_uuid = ?")

	if err != nil {
		fmt.Println("Failed to Execute Update Query: ", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(details.Duration, details.Total_Marks,details.Cource_UUID)

	if err != nil {
		fmt.Println("Failed To Update Exam Details: ", err)
		return err
	}



	return nil

}

func Create_Exam_Details(details Exam_Details) error {
	
	uuid := encription.Generateuudi()


	dbcreate := dbcode.SqlRead().DB
	stmt, err := dbcreate.Prepare("insert into exam_details(uuid,cource_uuid,program_name,cource_name,cource_code,duration, total_marks) values(?,?,?,?,?,?,?)")

	if err != nil {
		fmt.Println("PREPARE STATEMENT FAILED: ", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(uuid, details.Cource_UUID, details.Program_Name, details.Cource_Name, details.Cource_Code, details.Duration, details.Total_Marks)


	if err != nil {
		return err
	}

	return nil
}


func UpdateDurationTotalMarks(){

}

func Create_Exam(question_in Questions_Construct) (bool, string) {
	result := true
	create_exam := dbcode.SqlRead().DB

	stmt, err := create_exam.Prepare("insert into exam_questions(uuid, section,cource_uuid,cource_name,question,answer, question_number) values(?,?,?,?,?,?,?)")

	if err != nil {
		log.Fatal("Failed to entere exam", err)
	}

	defer stmt.Close()

	uuid := encription.Generateuudi()
	var section string
	var cource_uuid string
	var cource_name string
	var answer string
	var question string
	var question_number int

	section_enter := question_in.Section

	switch section_enter {
	case "A":
		cource_name = question_in.Cource_Name
		cource_uuid = question_in.Cource_UUID
		section = question_in.Section
		answer = question_in.Answer
		question = question_in.Question
		question_number = question_in.Question_Number

	case "B":

		section = question_in.Section
		cource_name = question_in.Cource_Name
		cource_uuid = question_in.Cource_UUID
		answer = question_in.Answer
		question = question_in.Question
		question_number = question_in.Question_Number

	}

	_, err = stmt.Exec(uuid, section, cource_uuid, cource_name, question, answer, question_number)

	if err != nil {
		log.Fatal(err)
	}

	return result, uuid
}

func DeleteAllQuestions(cource_name string) bool {

	deleted := true
	dbconn := dbcode.SqlRead().DB
	stmt, err := dbconn.Prepare("DELETE FROM exam_questions WHERE  cource_uuid = ?")

	if err != nil {
		deleted = false
		fmt.Println("deleleting problem")
		log.Fatal(err)

	}

	defer stmt.Close()

	_, err = stmt.Exec(cource_name)

	if err != nil {

		fmt.Println("failed to delete questions")
		deleted = false
	}

	return deleted
}



func Delete_Exam(w http.ResponseWriter, r *http.Request) {

	uuid := r.URL.Query().Get("uuid")
	delete := dbcode.SqlRead().DB

	stmt, err := delete.Prepare("DELETE FROM exam_questions WHERE  uuid = ?")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(uuid)

	if err != nil {
		log.Fatal(err)
	}

}

func Read_Exam_Taken(uuid, cource_name string) (bool, ExamTakenStruct) {

	exam_taken := true

	var exdata ExamTakenStruct

	dbread := dbcode.SqlRead().DB

	stmt, err := dbread.Prepare("SELECT uuid, cource_name, student_uuid, attemp_number, first_attempted, open_period, grade, comment, passed,grading, date FROM write_exam WHERE student_uuid = ? AND cource_name = ?")

	if err != nil {
		fmt.Println("EXAM TAKEN PREPARE STATEMENT FAILED: ", err )
		exam_taken = false
	}


	err = stmt.QueryRow(uuid, cource_name).Scan(&exdata.UUID, &exdata.Cource_Name, &exdata.Student_UUID, &exdata.Attemp_Number, &exdata.First_Attempted, &exdata.Open_Period, &exdata.Grade, &exdata.Comment, &exdata.Passed, &exdata.Gradding,&exdata.Date )

	if err != nil {
		fmt.Println("EXAM TAKEN PREPARE STATEMENT FAILED: ", err )
		exam_taken = false
	}

	return exam_taken, exdata

}


func CourceCompleted(w http.ResponseWriter, r *http.Request){
	cource_name := r.URL.Query().Get("cource_name")
	uuid := r.URL.Query().Get("student")
	
	fmt.Println(cource_name, uuid)
	
	id_out := CleanStudentUUID(uuid)
	var template_name string
	
	
	is_taken, _ := Read_Exam_Taken(uuid, cource_name)
	
	if is_taken {
		fmt.Println(id_out)
		template_name = "cource_complete_tag"
	}else{
		template_name = "cource_notcomplete_tag"
		
	}
	
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	
	
	err := tpl.ExecuteTemplate(w, template_name, nil)
	
	if err != nil {
		log.Fatal(err)
	}
}


func SaveQuestionUpdates(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	

	uuid := r.URL.Query().Get("uuid")
	section := r.URL.Query().Get("section")

	dbread  := dbcode.SqlRead().DB
	var question string
	var answer string

	type UpdateOut struct{
		UUID string
		Question string
		Answer string

	}
	var data_out UpdateOut
	var templateOut string 

	switch section {
		case "A":
			question = r.FormValue("question_a")
			answer = r.FormValue("answer")

			stmt, err := dbread.Prepare("UPDATE exam_questions SET question = ?, answer = ? WHERE UUID = ?")

			if err != nil {
				fmt.Println("FAILED PREPARE UPDATE FOR SECTION A:: ", err)
			}

			defer stmt.Close()


			_, err = stmt.Exec(question, answer, uuid)

			if err != nil {
				fmt.Println("FAILED TO UPDATE: ", err)

			}
			data_out = UpdateOut{
				UUID: uuid,
				Question: question,
				Answer: answer,
			}

			templateOut = "question_section_a"


		case "B":
			question = r.FormValue("question_b")
			

			stmt, err := dbread.Prepare("UPDATE exam_questions SET question = ?  WHERE UUID = ?")

			if err != nil {
				fmt.Println("FAILED PREPARE UPDATE FOR SECTION A:: ", err)
			}

			defer stmt.Close()


			_, err = stmt.Exec(question, uuid)

			if err != nil{
				fmt.Println("FAILED TO UPDATE:: ", err)
			}
			data_out = UpdateOut{
				UUID: uuid,
				Question: question,
				
			}

			templateOut = "question_section_b"

	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))
	err := tpl.ExecuteTemplate(w, templateOut, data_out)

	if err != nil {
		log.Fatal(err)
	}
}

func CreatePage(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	user_name, err := GetUserName(r)

	program_data := r.URL.Query().Get("uuid")

	fmt.Println("Cource UUID: ", program_data)
	var exam_value string
	
	get_exam_details := GetExamDetails(program_data)

	get_program_data := GetProgramDetailsSingle(program_data)
	
	
	dbread := dbcode.SqlRead().DB

	stmt, err := dbread.Prepare("SELECT exam_file FROM cource_table WHERE uuid = ?")
	if err != nil {
		fmt.Println("FAILED TO GET EAM VALUE")
	}
	defer stmt.Close()

	err = stmt.QueryRow(program_data).Scan(&exam_value)

	if err != nil {
		fmt.Println("FAILED TO GET EXAM VALUE")
	}
	

	var to_show ExamOut

	cource_name_out := get_program_data.Cource_Name

	formatCourceName := Clean(cource_name_out)
	
	

	if exam_value == "true"{

		result_out, _ := ReadQuestions(cource_name_out)
		
		
		to_show = ExamOut{
			Present:     true,
			Cource_Data: get_program_data,
			ExamData:    result_out,
			ExamDetails: get_exam_details,
			CourceNameFm: formatCourceName,
			Admin_Name: user_name,
		}

	} else {
		to_show = ExamOut{
			Present:     false,
			Cource_Data: get_program_data,
			CourceNameFm: formatCourceName,
			ExamDetails: get_exam_details,
			Admin_Name: user_name,

		}

	}

	err = tpl.ExecuteTemplate(w, "create_exam.html", to_show)

	if err != nil {
		log.Fatal(err)
	}

}

func ExtraChance(w http.ResponseWriter, r *http.Request){

}

func TakeExam(w http.ResponseWriter, r *http.Request) {

	

	funcMap := template.FuncMap{
		"ToCapital": ToUpperCase,
		"Clean":     Clean,
	}

	var display_number int

	t := template.New("").Funcs(funcMap)

	tpl = template.Must(t.ParseGlob("templates/*.html"))

	var template_name string

	// GET VARIABLES FROM ROUTE
	cource_name := r.URL.Query().Get("cource_name")
	uuid := r.URL.Query().Get("uuid")

	// VERIFY IF EXAM HAS BEEN ATTEMPTED
	if_taken, result_out := Read_Exam_Taken(uuid, cource_name)

	var attemped_number int
	var write_exam_uuid string
	var open_period int
	var passed string
	var attempt_out string
	var gradding string
	fmt.Println("POINT ONE: ")

	
	if if_taken {
		
		attemped_number, _ = strconv.Atoi(result_out.Attemp_Number)
		


		write_exam_uuid = result_out.UUID
		passed = result_out.Passed
		gradding = result_out.Gradding

		open_period, _ = strconv.Atoi(result_out.Open_Period)

		if gradding == "true"{
			
			fmt.Println("POINT THREE")
			err := tpl.ExecuteTemplate(w, "examgradding.html", nil)

			if err != nil {
				log.Fatal(err)
			}
			return

		}else{
			

			if passed == "true" {
				
				template_name = "exam_passed.html"
			

			} else{
				

				
				if attemped_number > 3 || open_period > 7 {
					display_number = attemped_number;
				

				} else {
					
					template_name = "exam_code.html"
					fmt.Println("Hitting the passed failed less than 3")
					
					_, attempt_out = UpdateExamTaken(true, uuid, cource_name, attemped_number)
				
					attemped_number, _ = strconv.Atoi(attempt_out)
					display_number = attemped_number


				}

			}
		}

	} else {

		open_period = 7
		template_name = "exam_code.html"

		display_number = 1
		_, attempt_out = UpdateExamTaken(false, uuid, cource_name,attemped_number)
		attemped_number, _ = strconv.Atoi(attempt_out)

	}

	
	get_exam_questions, cource_uuid_out, questiions_present := Read_Exam(cource_name)


	get_exam_details := GetExamDetails(cource_uuid_out)

	cource_name_out := cource_name

	if display_number <= 3 {
		display_data := DisplayExam{
			AlreadyTaken:      false,
			Qusetions_Present: questiions_present,
			Student_UUID:      uuid,
			Cource_UUIDOut:    cource_uuid_out,
			Attempt_Number:    attemped_number,
			ExamData:          get_exam_details,
			Exam_Questions:    get_exam_questions,
			Cource_Name_Two:   cource_name_out,
			Writen_UUID:       write_exam_uuid,
			Attempt_Out:       attempt_out,
		}
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") 
		w.Header().Set("Pragma", "no-cache") 
		w.Header().Set("Expires", "0") 
		err := tpl.ExecuteTemplate(w, template_name, display_data)

		if err != nil {
			log.Fatal(err)
		}
		return

	} else if ( display_number > 3 ){
		
		err := tpl.ExecuteTemplate(w, "exam_failed.html", nil)

		if err != nil {
			log.Fatal(err)
		}
		return
	}

}

func CheckForCource(uuid, question string) {

}

func Question_Count(uuid string) (bool, string) {

	dbcounter := dbcode.SqlRead().DB
	present := true
	counter_present := false
	var number string
	var number_list []string
	var last_number string

	

	stmt, err := dbcounter.Query("SELECT question_number  FROM exam_questions WHERE cource_uuid = ?", uuid)

	if err != nil {
		fmt.Println("failed to work properly", err)
		present = false
	}

	defer stmt.Close()

	for stmt.Next() {
		err = stmt.Scan(&number)
		if err != nil {
			fmt.Println("failed to query row: ", err)
			break
		} else {
			number_list = append(number_list, number)
			counter_present = true
		}

	}

	if err = stmt.Err(); err != nil {
		fmt.Println("stmt scan failed, error out: ", err)
	}

	if counter_present {
		last_number = number_list[len(number_list)-1]
		number = last_number
	}

	fmt.Println("Length of number list: ", len(number_list))

	return present, number

}

func GetExamDetails(courceuuid string) ExamDetails {

	var exam_details ExamDetails

	
	dbconn := dbcode.SqlRead().DB

	defer dbconn.Close()

	stmt, err := dbconn.Prepare("select program_name, cource_name, cource_code, duration, total_marks from exam_details where cource_uuid = ?")

	if err != nil {
		fmt.Println("Prepare Statement failed  in 'GetExamDetails' error out : ", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(courceuuid).Scan(&exam_details.Program_Name, &exam_details.Cource_Name, &exam_details.Cource_Code, &exam_details.Exam_Duration, &exam_details.Total_Marks)

	if err != nil {
		fmt.Println("Failed to get program details: ", err)


	}

	

	return exam_details

}

func QuestionUUID(cource_uuid string) []string {

	dbcounter := dbcode.SqlRead().DB

	var uuid_list []string
	var question_uuid string

	fmt.Println("The UUID ENTERED: ", cource_uuid)

	stmt, err := dbcounter.Query("SELECT uuid  FROM exam_questions WHERE cource_uuid = ?", cource_uuid)

	if err != nil {
		fmt.Println("failed to work properly", err)

	}

	defer stmt.Close()

	for stmt.Next() {
		err = stmt.Scan(&question_uuid)
		if err != nil {
			fmt.Println("failed to query row: ", err)
			break
		} else {
			uuid_list = append(uuid_list, question_uuid)

		}

	}

	if err = stmt.Err(); err != nil {
		fmt.Println("stmt scan failed, error out: ", err)
	}

	return uuid_list

}


func GetQuestionData(question_uuid string) QuestionData {

	var question_data QuestionData

	dbconn := dbcode.SqlRead().DB

	stmt, err := dbconn.Prepare("select uuid, question, question_number from exam_questions where uuid = ?")

	if err != nil {
		fmt.Println("Failed to get qustin data, error text: ", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(question_uuid).Scan(&question_data.Question_UUID, &question_data.Question, &question_data.Question_Number)

	if err != nil {
		fmt.Println("Failed to get question data", err)
	}

	return question_data

}

func SubTwoDay(date_two string) int64 {

	date_now := time.Now()
	date_t := strings.Split(date_two, ".")

	var days_out int64

	year_out := date_now.Year()
	month_out := date_now.Month()
	day_out := date_now.Day()

	year_t, _ := strconv.Atoi(date_t[0])
	month_t, _ := strconv.Atoi(date_t[1])
	day_t, _ := strconv.Atoi(date_t[2])

	form_day_one := time.Date(year_out, month_out, day_out, 0, 0, 0, 0, time.UTC)

	form_day_two := time.Date(year_t, time.Month(month_t), day_t, 0, 0, 0, 0, time.UTC)

	difference := form_day_one.Sub(form_day_two)

	days_out = int64(difference.Hours() / 24)
	return days_out
}


func IsGradding(cource_name, student_uuid string){
	dbread := dbcode.SqlRead().DB


	stmt, err := dbread.Prepare("UPDATE write_exam SET grading = ? WHERE cource_name = ?  AND student_uuid = ?")

	if err != nil {
		fmt.Println("FAILED TO GET PREPARE STATEMENT: ", err)

	}

	defer stmt.Close()

	_,err = stmt.Exec("true", cource_name, student_uuid)


	if err != nil {
		fmt.Println("FAILED TO UPDATE UPDATE: ", err)
	}
}


func UpdateExamTaken(is_taken bool, student_uuid string, cource_name string, attemped_number int) (bool, string) {

	updated := true
	dbconn := dbcode.SqlRead().DB
	uuid := encription.Generateuudi()
	var new_attmp string
	var attempt_string string



	if is_taken {

		
		uuid_out := student_uuid
		attemp_number_out := attemped_number
		if attemp_number_out < 3 {
			attempt_out := attemp_number_out + 1
			attempt_string = strconv.Itoa(attempt_out)

			new_attmp = attempt_string
			stmt, err := dbconn.Prepare("UPDATE write_exam SET attemp_number = ? WHERE student_uuid = ? AND cource_name = ?")

	
			if err != nil {
				fmt.Println("PREPARE STATEMENT FAILED ", err)
			}

			defer stmt.Close()

			_, err = stmt.Exec(new_attmp, uuid_out,cource_name)

			if err != nil {
				fmt.Println("FAILED TO UPDATE THE EXAM DETAILS", err)
			}
		}else{
			attempt_string = strconv.Itoa(4)

			new_attmp = attempt_string


		}

		
		

	} else {

		attemped_number := 1
		first_attempted := time.Now().Format("2017.09.07")
		open_period := 7
		grade := ""
		comment := ""
		passed := ""
		grading := "false"
		date := time.Now().Format("2017.09.07")

		new_attmp = strconv.Itoa(attemped_number)

		stmt, err := dbconn.Prepare("insert into write_exam(uuid, student_uuid, cource_name,attemp_number, first_attempted, open_period, grade,comment,passed,grading, date) values(?,?,?,?,?,?,?,?,?,?,?)")

		if err != nil {
			fmt.Println("Failed to insert into write_exam one: ", err)
		}

		defer stmt.Close()

		_, err = stmt.Exec(uuid, student_uuid, cource_name, attemped_number, first_attempted, open_period, grade, comment, passed,grading, date)

		if err != nil {
			fmt.Println("Failed to insert into write_exam two:", err)
		}

	}

	return updated, new_attmp
}

func SubmitExam(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	student_uuid := r.URL.Query().Get("student_uuid")
	cource_uuid := r.URL.Query().Get("cource_uuid")
	cource_name := r.URL.Query().Get("cource_name")

	_, data_out := Read_Exam_Taken(student_uuid, cource_name)

	attempt_number := data_out.Attemp_Number

	fmt.Println("Cource_UUID: ", cource_uuid)

	question_uuids := QuestionUUID(cource_uuid)

	
	IsGradding(cource_name, student_uuid)

	for _, item := range question_uuids {

		answer := r.FormValue(item)
		data_out := GetQuestionData(item)

		answer_out := Answer_Out{
			Cource_UUID:     cource_uuid,
			Student_UUID:    student_uuid,
			Qustion_UUID:    item,
			Question_Number: data_out.Question_Number,
			Question:        data_out.Question,
			Attemp_Number:   attempt_number,
			Answer:          answer,
			Cource_Name:     cource_name,
		}

		RecordStudentMarks(answer_out)

	}

	

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "exam_taken.html", student_uuid)

	if err != nil {
		log.Fatal(err)
	}

}

func UpdateExamEntered(uuid string){
	fmt.Println("WHAT  THE FUCL IS GOING")
	update_saved := true 
	dbread := dbcode.SqlRead().DB
	
	
	stmt, err := dbread.Prepare("UPDATE cource_table SET exam_file = ? WHERE uuid = ?")

	if err != nil{
		fmt.Println("FAILED TO LOAD COURCE TABLE:: ", err)
		update_saved = false
	
	}
	
	defer stmt.Close()
	
	_, err = stmt.Exec("true", uuid)
	if err != nil {
		fmt.Println("FIALED TO MAKE UPDATE::: ", err)
		update_saved = false
	}

	
	fmt.Println("WHAT  THE FUCL IS GOING", update_saved)
	
	
}

func AddExamDetails(w http.ResponseWriter, r *http.Request) {

	
	
	var exam_responce CreateExamResponse
	var template_name string
	var details_message string


	cource_uuid := r.URL.Query().Get("uuid")
	exam_time := r.FormValue("exam_time")
	total_marks := r.FormValue("total_marks")



	create_exam_detaile := Exam_Details{
		Cource_UUID: cource_uuid,
		Duration:     exam_time,
		Total_Marks:  total_marks,
	}


	
	err := Update_Exam_Details(create_exam_detaile)

	if err != nil{
		fmt.Println("Failed to update exam details")
	}
	exam_responce = CreateExamResponse{
		Details_Messages: details_message}

	template_name = "details_saved_temp false"

	

	fmt.Println("The Responce",exam_responce)
	
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	fmt.Println(template_name)

	err = tpl.ExecuteTemplate(w,  "durationTime", create_exam_detaile)

	if err != nil {
		log.Fatal(err)
	}
}


func UpdateQuestion(w http.ResponseWriter, r *http.Request){

	uuid := r.URL.Query().Get("qustion_uuid")
	section := r.URL.Query().Get("section")
	var updateTextArea string
	var question Questions_Construct
	dbread := dbcode.SqlRead().DB


	fmt.Println("THE UUID IS OUT:: ", uuid)
	stmt, err := dbread.Prepare("SELECT uuid, section, cource_uuid, cource_name,question, answer FROM exam_questions  WHERE uuid = ?")

	if err != nil{
		fmt.Println("FAILED TO PREPARE STATEMENT FOR SELECT:: ", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(uuid).Scan(&question.UUID, &question.Section,&question.Cource_UUID, &question.Cource_Name, &question.Question, &question.Answer)

	if err != nil {
		fmt.Println("FAILED TO GET QUESTION DATA", err)
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	fmt.Println("The Section: ", section)
	switch section{
		case "A":
			updateTextArea = "sectionAUpdate"
		case "B":
			updateTextArea = "sectionBUpate"
	}

	fmt.Println("THE QUESTION:::",question.Question)

	err = tpl.ExecuteTemplate(w, updateTextArea, question)

	if err != nil {
		log.Fatal(err)
	}



}

func DeleteQuestion(w http.ResponseWriter, r *http.Request){
	

	uuid := r.URL.Query().Get("qustion_uuid")
	fmt.Println("THE UUID: ", uuid)
	dbread := dbcode.SqlRead().DB
	stmt, err := dbread.Prepare("DELETE FROM exam_questions WHERE uuid = ?")

	if err != nil {
		fmt.Println("FAILED TO CREATE PREPARE DELETE STATEMENT: ", err)

	}

	defer stmt.Close()

	_, err = stmt.Exec(uuid)


	if err != nil{
		fmt.Println("FAILED TO DELETE: ", err)
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, "empty_div_none", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func AddExam(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var is_present bool
	var number string
	
	var exam_responce Questions_Construct

	var template_name string

	question_section := r.URL.Query().Get("section")
	cource_name := r.URL.Query().Get("cource_name")
	cource_uuid := r.URL.Query().Get("uuid")
	
	exam_time := r.FormValue("questionType")
	

	

	exam_time_out, _ := strconv.Atoi(exam_time)

	fmt.Println(exam_time_out)

	
	is_present, number = Question_Count(cource_uuid)



	switch question_section {

	case "A":

		question_a := r.FormValue("questionText")
		answers := r.FormValue("answer")

		if len(answers) < 1 {
			answers = "false"
		}

		var question_content Questions_Construct

		if is_present {

			number_out, _ := strconv.Atoi(number)
			

			question_content = Questions_Construct{
				Cource_UUID:     cource_uuid,
				Section:         "A",
				Cource_Name:     cource_name,
				Question:        question_a,
				Answer:          answers,
				Question_Number: number_out + 1,
			}

			_, quuid := Create_Exam(question_content)

			UpdateExamEntered(cource_uuid)
			
			question_content.UUID = quuid

			exam_responce = question_content

			fmt.Println(exam_responce)

			template_name = "question_section_a"
			

		} else {

			number_out := 1
			
			question_content = Questions_Construct{
				Cource_UUID:     cource_uuid,
				Cource_Name:     cource_name,
				Section:         "A",
				Question:        question_a,
				Answer:          answers,
				Question_Number: number_out,
			}

			_, quuid := Create_Exam(question_content)
			UpdateExamEntered(cource_uuid)
			

			question_content.UUID = quuid

			exam_responce = question_content
			
				
			template_name = "examquestion"
			
		}

	case "B":
		

		question_b := r.FormValue("questionText")

		var question_content Questions_Construct

		if is_present {

			number_out, _ := strconv.Atoi(number)

			question_content = Questions_Construct{
				Cource_UUID:     cource_uuid,
				Section:         "B",
				Cource_Name:     cource_name,
				Question:        question_b,
				Question_Number: number_out + 1,
			}

			_, quuid := Create_Exam(question_content)

			UpdateExamEntered(cource_uuid)
			

			question_content.UUID = quuid

			exam_responce = question_content

			template_name = "question_section_b"
			

		} else {

			number_out := 1
			question_content = Questions_Construct{
				Cource_UUID:     cource_uuid,
				Cource_Name:     cource_name,
				Section:         "B",
				Question:        question_b,
				Question_Number: number_out,
			}

			_, quuid := Create_Exam(question_content)
			UpdateExamEntered(cource_uuid)
		
			question_content.UUID = quuid

			exam_responce = question_content

			template_name = "question_section_b"

			
		}

	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	fmt.Println(template_name)

	err := tpl.ExecuteTemplate(w, template_name, exam_responce)

	if err != nil {
		log.Fatal(err)
	}
}

func OpenGradeExam(w http.ResponseWriter, r *http.Request) {

	student_uuid := r.URL.Query().Get("student_uuid")
	cource_name := r.URL.Query().Get("courcer_name")
	dbconn := dbcode.SqlRead().DB

	var answer_out Answer_Out
	var answers_out []Answer_Out

	clean_uuid := CleanStudentUUID(student_uuid)

	_, data_out := Read_Exam_Taken(student_uuid, cource_name)
	attempt_number := data_out.Attemp_Number

	fmt.Println(clean_uuid, data_out)

	query_statement := fmt.Sprintf("select * from %s ")

	stmt, err := dbconn.Query(query_statement)

	if err != nil {
		fmt.Print("Failed to To Load Prepare statement")
	}

	defer stmt.Close()

	for stmt.Next() {
		err = stmt.Scan(&answer_out.Qustion_UUID, &answer_out.Cource_UUID, &answer_out.Student_UUID, &answer_out.Question_Number, &answer_out.Question, &answer_out.Attemp_Number, &answer_out.Answer)
		attempt_number_out := &answer_out.Attemp_Number
		cource_name_out := &answer_out.Cource_Name

		if *cource_name_out == cource_name && *attempt_number_out == attempt_number {
			answers_out = append(answers_out, answer_out)
		}

		if err != nil {
			fmt.Println("Failed to get questions and answers out")

		}

	}

	type Display_To_Grade struct {
	}

	to_grade := Display_To_Grade{}

	err = tpl.ExecuteTemplate(w, "grade_exam.html", to_grade)

	if err != nil {
		log.Fatal(err)
	}
}

func LoadExamTable() {

	exam_code := dbcode.SqlRead().DB

	defer exam_code.Close()

	exam_details := `create table if not exists exam_details(
		uuid blob not null,
		cource_uuid text,
		program_name text,
		cource_name text,
		cource_code  text,
		duration int,
		total_marks string
	)`

	_, exam_details_error := exam_code.Exec(exam_details)
	if exam_details_error != nil {
		log.Printf("%q: %s\n", exam_details_error, exam_details)
	}

	create_exam := `
		create table if not exists exam_questions(
		uuid blob not null,
		section text,
		cource_uuid text,
		cource_name text,
		question text,
		answer text,
		question_number int)`

	_, create_exam_error := exam_code.Exec(create_exam)
	if create_exam_error != nil {
		log.Printf("%q: %s\n", create_exam_error, create_exam)
	}

	take_exam := `
		create table if not exists write_exam(
		uuid blob not null,
		cource_name text,
		student_uuid text,
		attemp_number int,
		first_attempted int,
		open_period int,
		grade text,
		comment text,
		passed bool,
		grading text,
		date text)`

	_, take_exam_error := exam_code.Exec(take_exam)
	if create_exam_error != nil {
		log.Printf("%q: %s\n", take_exam_error, take_exam)
	}

}
