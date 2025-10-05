package routes

import (
	"fmt"
	"log"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"
)

type StudentProgramData struct {
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
}

type CourceData struct {
	UUID string
}

// CORRECT THE ACAMS STRUCT ERROR

func ApplyForCource(uuid, cource_name string) bool {
	cource_applied := true

	fmt.Println("Cource name: ", uuid, cource_name)
	dbread := dbcode.SqlRead().DB

	dbupdate_statement := fmt.Sprintf(`UPDATE %s SET applied = ? WHERE uuid = ? `, cource_name)

	statement, err := dbread.Prepare(dbupdate_statement)

	if err != nil {
		error_text := fmt.Sprintf("line 44 error from update prepare:: %s", err)
		ErrorPrintOut("studentportal", "ApplyForCource", error_text)

		cource_applied = false
	}

	defer statement.Close()

	_, errup := statement.Exec(true, uuid)

	if errup != nil {
		error_text := fmt.Sprintf("line 50 error from update prepare:: %s", errup)
		ErrorPrintOut("studentportal", "ApplyForCource", error_text)
		cource_applied = false
	}
	return cource_applied
}

func AddToProgramCources(student_uuid, date_in, payment_type, program_code string) bool {
	dbcreate := dbcode.SqlRead().DB

	course_tables, err := GetProgramCourses(program_code)

	if err != nil {
		fmt.Println("Failed to get program courses: ", err)
	}

	added_to_cource_table := true

	for _, item := range course_tables {

		uuid := encription.Generateuudi()

		cource_create, err := dbcreate.Begin()

		if err != nil {
			error_out := fmt.Sprintf("AddToProgramCources: %s", err)
			ErrorPrintOut("programcoursecreate", "AddToProgramCources", error_out)
		}

		insert_String := fmt.Sprintf(`insert into %s(
			uuid,
			student_uuid,
			cource_name,
			course_code,
			book,
			module,
			video,
			applied,
			approved,
			examined,
			continuorse_assesment,
			completed,
			date) values(?,?,?,?,?,?,?,?,?,?,?,?,?)`, item.Code)

		statment, err := cource_create.Prepare(insert_String)

		if err != nil {
			error_out := fmt.Sprintf("failed to insert student in: %s,error: %s", item, err)

			ErrorPrintOut("programcoursecreate", "AddToProgramCources", error_out)
		}

		defer statment.Close()

		var applied_value bool

		if payment_type == "lamp_sum" {
			applied_value = true
		} else {

			applied_value = false
		}
		var cource_name = item
		var book string = fmt.Sprintf("Book for: %s", item)
		var video string = fmt.Sprintf("video for: %s", item)
		var module string = fmt.Sprintf("video for: %s", item)
		var applied bool = applied_value
		var approved bool = false
		var examinde bool = false
		var continuorse_assesment string = "0"
		var completed bool = false

		_, err = statment.Exec(
			uuid,
			student_uuid,
			cource_name,
			book,
			video,
			module,
			applied,
			approved,
			examinde,
			continuorse_assesment,
			completed,
			date_in,
		)

		if err != nil {
			error_out := fmt.Sprintf("failde to aad to cource table: %s, error: %s", item, err)
			ErrorPrintOut("programcoursecreate", "AddToProgramCources", error_out)

		}

		err = cource_create.Commit()

		if err != nil {
			error_out := fmt.Sprintf("failde to commit to cource table: %s, error: %s", item, err)
			ErrorPrintOut("programcoursecreate", "AddToProgramCources", error_out)

		}

	}

	return added_to_cource_table
}

func GetFromProgramCourcesAdmin(student_uuid, program_code string) []CourceStruct{
	var cource_data_out CourceStruct

	var cource_data_out_list []CourceStruct

	course_tables, err := GetProgramCourses(program_code)

	if err != nil {
		fmt.Println("Failed program courses: ", err)
	}

	for _, item := range course_tables {
		dbread := dbcode.SqlRead().DB

		data_query_string := fmt.Sprintf("select uuid, student_uuid, cource_name,course_code, applied, approved, examined, completed, date from %s  where student_uuid = ?", item.Code)

		stmt, err := dbread.Prepare(data_query_string)

		if err != nil {
			fmt.Println(err)
		
		}

		defer stmt.Close()

		err = stmt.QueryRow(student_uuid).Scan(&cource_data_out.UUID,
			&cource_data_out.Student_UUID,
			&cource_data_out.Cource_Name,
			&cource_data_out.Course_Code,
			&cource_data_out.Applied,
			&cource_data_out.Approved,
			&cource_data_out.Examined,
			&cource_data_out.Completed,
			&cource_data_out.Date,
		)

		if err != nil {
			error_out := fmt.Sprintf("falied to query row: %s", err)
			ErrorPrintOut("acamsfile", "GetFromACAMSCource", error_out)
			log.Fatal(err)
		}

		cource_data_out_list = append(cource_data_out_list, cource_data_out)

	}

	return cource_data_out_list
}
//
//


///IF NOT IN USE TO BE DELETED

func GetCourseMaterialOne(courseCode string)(string, string, string){

	dbread := dbcode.SqlRead().DB

	var assesment_out string
	var video_list string
	var module string

	defer dbread.Close()

	stmt, err := dbread.Prepare("select  cource_assesment, video_list, module from cource_table where cource_code = ?")

	if err != nil {
		fmt.Println("Failed to get cource material")
	}

	err = stmt.QueryRow(courseCode).Scan(&assesment_out, &video_list, &module)


	if err != nil{
		fmt.Println("Failed to get courses")
	}


	return assesment_out, video_list, module



}

func GetFromProgramCources(student_uuid_in, program_code string) []CourceStruct {
	var cource_data_out CourceStruct

	var cource_data_out_list []CourceStruct

	course_tables, err := GetProgramCourses(program_code)

	if err != nil {
		fmt.Println("Failed program courses: ", err)
	}

	dbread := dbcode.SqlRead().DB

	defer dbread.Close()

	var uuid string
	var student_uuid string
	var cource_name string
	var cource_code string
	var applied bool
	var approved bool
	var examind bool
	var completed bool
	var date string 



	for _, item := range course_tables {
		

		data_query_string := fmt.Sprintf("select uuid, student_uuid, cource_name,course_code, applied, approved, examined,completed, date from %s  where student_uuid = ?", item.Code)

		stmt, err := dbread.Prepare(data_query_string)

		if err != nil {
			fmt.Println(err)
		
		}
		defer stmt.Close()

		err = stmt.QueryRow(student_uuid_in).Scan(
			&uuid,
			&student_uuid,
			&cource_name,
			&cource_code,
			&applied,
			&approved,
			&examind,
			&completed,
			&date,
		)



		//SET THE COURCE MATERIAL INTO STRUCT

		ass, video, module := GetCourseMaterialOne(item.Code)

		cource_data_out = CourceStruct{
			UUID: uuid,
			Student_UUID: student_uuid,
			Cource_Name: cource_name,
			Course_Code: cource_code,
			Applied: applied,
			Approved: approved,
			Examined: examind,
			Continuorse_Assesment: ass,
			Completed: completed,
			Date: date,
			Module: module,
			Video: video,
		}

		if err != nil {
			error_out := fmt.Sprintf("falied to query row: %s", err)
			ErrorPrintOut("programcoursecreate", "GetProgramCources", error_out)
			log.Fatal(err)
		}

		cource_data_out_list = append(cource_data_out_list, cource_data_out)

	}

	return cource_data_out_list
}


//TO BE DELETED


//
//1. add student data to the table of the program applied for 
//2. get list of cources that are associated with that program
//3. add the student data to the tables of said cources

func CreateProgamData(data_in StudentProgramData, payment_type, program_code string) error {
	

	dbcreate := dbcode.SqlRead().DB

	uuid := encription.Generateuudi()

	student_create, err := dbcreate.Begin()

	if err != nil {
		error_out := fmt.Sprintf("%s", err)
		ErrorPrintOut("programcreate", "Create Program data", error_out)
		return err
	}

	prepare_str := fmt.Sprintf(`insert into %s(
		uuid,
		student_uuid,
		program_name,
		first_name,
		last_name,
		email,
		applied,
		approved,
		payment_method,
		paid,
		completed,
		date) values(?,?,?,?,?,?,?,?,?,?,?,?)`, program_code)

	statment, err := student_create.Prepare(prepare_str)

	if err != nil {
		error_out := fmt.Sprintf("the prepare statment: %s", err)
		ErrorPrintOut("programcoursecreate", "CreateProgram", error_out)
		return err

	}

	defer statment.Close()

	_, err = statment.Exec(
		uuid,
		data_in.Student_UUID,
		program_code,
		data_in.First_Name,
		data_in.Last_Name,
		data_in.Email,
		data_in.Applied,
		data_in.Approved,
		data_in.Payment_Method,
		data_in.Paid,
		data_in.Completed,
		data_in.Date,
	)

	if err != nil {
		error_out := fmt.Sprintf("execusion statement: %s", err)

		ErrorPrintOut("acams", "CreateACAMS: ", error_out)
		return err

	}

	err = student_create.Commit()

	if err != nil {
		error_out := fmt.Sprintf("commit statement: %s", err)

		ErrorPrintOut("programcoursecreate", "Create program data: ", error_out)
		return err

	}
	add_to_cources := AddToProgramCources(data_in.Student_UUID, data_in.Date, payment_type, program_code)

	fmt.Println("adding to cources was succesful", add_to_cources)

	fmt.Printf("Creating new  student in  database complete\n")
	return nil

}

func GetProgramAdmin(students_uuid_in, promt, program_code string) (bool, StudentProgramData, []StudentProgramData) {
	var confirmacms bool
	var acams_data_out StudentProgramData
	var acams_data_out_list []StudentProgramData

	var prepare_str string
	promtout := promt
	dbread := dbcode.SqlRead().DB

	switch promtout {

	case "one":

		prepare_str = fmt.Sprintf("select uuid, student_uuid, program_name,first_name, last_name, email, applied, approved, payment_method, paid, completed, date from %s where student_uuid = ?", program_code)
		statement, err := dbread.Prepare(prepare_str)

		if err != nil {
			error_out := fmt.Sprintf("%s prepare", err)
			ErrorPrintOut("programcoursecreate", "GetPrograms", error_out)
			confirmacms = false
		}

		defer statement.Close()

		err = statement.QueryRow(students_uuid_in).Scan(&acams_data_out.UUID, &acams_data_out.Student_UUID, &acams_data_out.Program_Name, &acams_data_out.First_Name, &acams_data_out.Last_Name, &acams_data_out.Email, &acams_data_out.Applied, &acams_data_out.Approved, &acams_data_out.Payment_Method, &acams_data_out.Paid, &acams_data_out.Completed, &acams_data_out.Date)

		if err != nil {
			error_out := fmt.Sprintf("%s assigning", err)
			ErrorPrintOut("programcoursecreate", "GetPogramAdmin", error_out)
		}

		fmt.Println("the approved tag", acams_data_out.Approved)

		if acams_data_out.Applied {
			confirmacms = true
		} else {
			confirmacms = false
		}
		return confirmacms, acams_data_out, acams_data_out_list

	case "multiple":
		prepare_str = fmt.Sprintf("select * from %s", program_code)
		rows, err := dbread.Query(prepare_str)

		if err != nil {
			error_out := fmt.Sprintf("getting multiple acams data: %s", err)
			ErrorPrintOut("programcoursecreate", "GetPrograms", error_out)
		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(
				&acams_data_out.UUID,
				&acams_data_out.Student_UUID,
				&acams_data_out.Program_Name,
				&acams_data_out.First_Name,
				&acams_data_out.Last_Name,
				&acams_data_out.Email,
				&acams_data_out.Applied,
				&acams_data_out.Approved,
				&acams_data_out.Payment_Method,
				&acams_data_out.Paid,
				&acams_data_out.Completed,
				&acams_data_out.Date,
			)

			if err != nil {
				error_out := fmt.Sprintf("getting multiple acams data for loop: %s", err)
				ErrorPrintOut("programcoursecreate", "GetPrograms", error_out)
			}

			acams_data_out_list = append(acams_data_out_list, acams_data_out)
		}

	}

	return confirmacms, acams_data_out, acams_data_out_list

}

func GetProgramsStudents(students_uuid_in, promt, program_name string) (bool, StudentProgramData, []StudentProgramData) {
	var confirmacms bool
	var acams_data_out StudentProgramData
	var acams_data_out_list []StudentProgramData

	promtout := promt
	dbread := dbcode.SqlRead().DB

	var prepare_str string

	switch promtout {

	case "one":
		prepare_str = fmt.Sprintf("select uuid, student_uuid, program_name,first_name, last_name, email, applied, approved, payment_method, paid, completed, date from %s where student_uuid = ?", program_name)
		statement, err := dbread.Prepare(prepare_str)

		if err != nil {
			error_out := fmt.Sprintf("%s prepare", err)
			ErrorPrintOut("programcoursecreate 395", "GetProgramsStudents", error_out)
			confirmacms = false
		}

		defer statement.Close()

		err = statement.QueryRow(students_uuid_in).Scan(&acams_data_out.UUID, &acams_data_out.Student_UUID, &acams_data_out.Program_Name, &acams_data_out.First_Name, &acams_data_out.Last_Name, &acams_data_out.Email, &acams_data_out.Applied, &acams_data_out.Approved, &acams_data_out.Payment_Method, &acams_data_out.Paid, &acams_data_out.Completed, &acams_data_out.Date)

		if err != nil {
			error_out := fmt.Sprintf("%s assigning", err)
			ErrorPrintOut("acams 405", "Get Student Program Data", error_out)
		}

		fmt.Println("the approved tag", acams_data_out.Approved)

		if acams_data_out.Approved {
			confirmacms = true
		} else {
			confirmacms = false
		}
		return confirmacms, acams_data_out, acams_data_out_list

	case "multiple":
		prepare_str = fmt.Sprintf("select * from %s", program_name)
		rows, err := dbread.Query(prepare_str)

		if err != nil {
			error_out := fmt.Sprintf("getting multiple acams data: %s", err)
			ErrorPrintOut("acams 422", "GetACAMS", error_out)
			log.Fatal(err)
		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(
				&acams_data_out.UUID,
				&acams_data_out.Student_UUID,
				&acams_data_out.Program_Name,
				&acams_data_out.First_Name,
				&acams_data_out.Last_Name,
				&acams_data_out.Email,
				&acams_data_out.Applied,
				&acams_data_out.Approved,
				&acams_data_out.Payment_Method,
				&acams_data_out.Paid,
				&acams_data_out.Completed,
				&acams_data_out.Date,
			)

			if err != nil {
				error_out := fmt.Sprintf("getting multiple acams data for loop: %s", err)
				ErrorPrintOut("acams 445", "GetACAMS", error_out)
			}

			acams_data_out_list = append(acams_data_out_list, acams_data_out)
		}

	}

	return confirmacms, acams_data_out, acams_data_out_list

}

func CreateCourseTable(course_code string) error {

	dbread := dbcode.SqlRead()

	defer dbread.DB.Close()

	create_course_table := fmt.Sprintf(`
	create table if not exists %s(uuid blob not null, 
		student_uuid text,
		cource_name text,
		course_code text,
		book text,
		module text,
		video text,
		applied bool,
		approved bool,
		continuorse_assesment text,
		examined bool,
		completed bool,
		date text);`, course_code)

	_, err := dbread.DB.Exec(create_course_table)
	if err != nil {
		log.Printf("%q: %s\n", err, create_course_table)
		return err
	}

	return nil

}

func CreateProgramTabel(table_name string) error {

	dbread := dbcode.SqlRead()

	defer dbread.DB.Close()

	create_program := fmt.Sprintf(`
		create table if not exists %s(
			uuid blob not null,
			student_uuid text,
			program_name text,
			first_name text,
			last_name text,
			email text,
			applied bool,
			approved bool,
			payment_method text,
			paid text,
			completed bool,
			date text
		);`, table_name)

	_, err := dbread.DB.Exec(create_program)
	if err != nil {
		log.Printf("%q: %s\n", err, create_program)
		return err
	}

	return nil
}
