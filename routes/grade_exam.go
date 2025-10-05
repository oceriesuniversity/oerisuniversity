package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kamukwamba/oerisuniversity/dbcode"
)

type Grade_Data struct {
	Exam_Data   ExamTakenStruct
	Exam_Detail ExamDetails
	Answers_Out []Answer_Out
	Assesment   []AssesmentGrade
}

func GetParticularExam(w http.ResponseWriter, r *http.Request) {

	student_uuid := r.URL.Query().Get("rttse")
	cource_name := r.URL.Query().Get("ttre")
	
	
	

	attempt_number := r.FormValue("attempt")

	cleaned := CleanStudentUUID(student_uuid)

	var grade_answer Answer_Out
	var grade_answer_list []Answer_Out

	dbconn := dbcode.SqlRead().DB

	_, cource_uuid, _ := Read_Exam(cource_name)
	exam_details := GetExamDetails(cource_uuid)
	_, attemp_out := Read_Exam_Taken(student_uuid, cource_name)
	
	attemp_out.Attemp_Number = attempt_number

	query_stmt := fmt.Sprintf("select * from  %s where cource_name = ? AND attempt_number = ?", cleaned)

	stmt, err := dbconn.Query(query_stmt, cource_name, attempt_number)

	if err != nil {
		fmt.Println("Failed to create prepare statement")
	}

	defer stmt.Close()

	for stmt.Next() {
		err = stmt.Scan(&grade_answer.UUID, &grade_answer.Cource_UUID, &grade_answer.Cource_Name, &grade_answer.Student_UUID, &grade_answer.Question_Number, &grade_answer.Question, &grade_answer.Attemp_Number, &grade_answer.Answer)

		if err != nil {
			err_text := fmt.Sprintf("failed to scan ")
			ErrorPrintOut("grade_exam.go", "GetParticularExam", err_text)
		}

		grade_answer_list = append(grade_answer_list, grade_answer)
	}
	

	grade_data := Grade_Data{
		Exam_Data:   attemp_out,
		Exam_Detail: exam_details,
		Answers_Out: grade_answer_list,
	}
	
	fmt.Println("Attempt Number: ", attempt_number, "\nData Out::::: ", grade_data)
	
	
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, "grade_left", grade_data)

	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
}

func SaveGrades(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	dbconn := dbcode.SqlRead().DB
	
	cource_name := r.URL.Query().Get("rttcn")
	student_uuid := r.URL.Query().Get("rttsu")



	grade := r.FormValue("total")
	passed := r.FormValue("passed")
	comment := r.FormValue("comment")
	grading := "false"

	fmt.Println("Route Has Been Hit: ", cource_name, student_uuid)

	stmt, err := dbconn.Prepare("UPDATE write_exam SET grade = ?, comment = ?, passed = ?, grading = ? WHERE student_uuid = ? AND cource_name = ?")

	if err != nil {
		fmt.Println("Prepare statement error: ", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(grade, comment, passed, grading,student_uuid, cource_name)

	if err != nil {
		fmt.Println("failed to update", err)
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, "grade_saved", nil)

	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

}

func CurrentAttemp(cource_name, studentuuid string) string{

	dbread := dbcode.SqlRead().DB
	var numberStr string 
	
	stmt, err := dbread.Prepare("SELECT attemp_number FROM write_exam WHERE cource_name = ? AND student_uuid = ?")

	if err != nil {
		fmt.Println("PREPARE STATEMENT FAILED", err)
	}

	defer stmt.Close()
	err = stmt.QueryRow(cource_name,studentuuid).Scan(&numberStr)

	if err != nil {
		fmt.Println("EXECUTION FAILED", err)
	}

	

	return numberStr

}

func GradeExam(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Cource grade Triggerd")
	tpl = template.Must(template.ParseGlob("templates/*.html"))



	student_uuid := r.URL.Query().Get("student_uuid")
	cource_name := r.URL.Query().Get("cource_name")

	var asse_data []AssesmentGrade
	present, assesment_data := GetAssesmentData(student_uuid, cource_name)

	if present {
		asse_data = assesment_data
	}

	if_present, attemp_out := Read_Exam_Taken(student_uuid, cource_name)
	dbconn := dbcode.SqlRead().DB

	if if_present {
		attemp_number := attemp_out.Attemp_Number
		fmt.Println(attemp_number)

		var grade_answer Answer_Out
		var grade_answer_list []Answer_Out

		var grade_data Grade_Data

		cleaned := CleanStudentUUID(student_uuid)



		_, cource_uuid_out, _ := Read_Exam(cource_name)

		exam_details := GetExamDetails(cource_uuid_out)

		fmt.Println("Cleaned: ", cleaned)

		query_string := fmt.Sprintf("select * from %s", cleaned)

		stmt, err := dbconn.Query(query_string)

		attemptnumber := CurrentAttemp(cource_name, student_uuid)

		if err != nil {
			fmt.Println("Failed to initialize prepare statement: ", err)
		}

		defer stmt.Close()

		for stmt.Next() {
			err = stmt.Scan(&grade_answer.UUID, &grade_answer.Cource_UUID, &grade_answer.Cource_Name, &grade_answer.Student_UUID, &grade_answer.Question_Number, &grade_answer.Question, &grade_answer.Attemp_Number, &grade_answer.Answer)

			if err != nil {
				fmt.Println("Failed to obtain scan: ", err)
			}
			attmpstr := grade_answer.Attemp_Number

			
			if(attemptnumber == attmpstr){
			grade_answer_list = append(grade_answer_list, grade_answer)}
			


		}

		if err != nil {
			fmt.Println("Failed to create prapare statement: ", err)
		}

		grade_data = Grade_Data{
			Exam_Data:   attemp_out,
			Exam_Detail: exam_details,
			Answers_Out: grade_answer_list,
			Assesment:   asse_data,
		}

		tpl = template.Must(template.ParseGlob("templates/*.html"))

		err = tpl.ExecuteTemplate(w, "grade_exam.html", grade_data)

		if err != nil {
			log.Fatal(err)
		}

	} else {
		tpl = template.Must(template.ParseGlob("templates/*.html"))

		err := tpl.ExecuteTemplate(w, "exam_nottaken.html", nil)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(student_uuid, cource_name)

}
