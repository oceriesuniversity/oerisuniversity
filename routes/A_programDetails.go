package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kamukwamba/oerisuniversity/dbcode"

	_ "github.com/mattn/go-sqlite3"
)

type Course_Name struct {
	ID    int
	Code  string
	Name  string
	PCode string
}

type ProgramDataEntry struct {
	ID          int
	Name        string
	Code        string
	CourseNames []Course_Name
}
type ProgramsAvailabelSt struct {
	Present      bool
	EmailPresent bool
	ProgramList  []ProgramDataEntry
}


func CreateProgramDB() {
	dbread := dbcode.SqlRead().DB

	programData := `CREATE TABLE IF NOT EXISTS ProgramData(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		programName TEXT UNIQUE,
		programCode TEXT UNIQUE
	);`

	defer dbread.Close()

	_, err := dbread.Exec(programData)

	if err != nil {
		log.Printf("%q: %s\n", err, programData)

	}

}



func CreateNewProgramR(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	program_name := r.FormValue("program_name")
	program_code := r.FormValue("program_code")
	render := false

	dataCreated := ProgramDataEntry{
		Name: program_name,
		Code: program_code,
	}

	err := ConfirmProgramDataExists(program_name, program_code)

	if err != nil {
		err = CreateProgramEntry(program_name, program_code)

		if err != nil {
			fmt.Printf("Failed to Create Program Entry:: %s", err)
		} else {
			fmt.Println("Program Entry Created Sucesfully")
			err = CreateProgramTabel(program_code)
			if err != nil {
				fmt.Println("Failed to create program table")
			}
			render = true
		}
	} else {
		fmt.Println("Program Exist With that Program Code Or Program Name")

	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	if render {
		err := tpl.ExecuteTemplate(w, "programCardTempEmpty", dataCreated)

		if err != nil {
			log.Fatal(err)
		}
	}

}

func ConfirmProgramDataExists(programName, programCode string) error {
	db := dbcode.SqlRead().DB
	defer db.Close()

	var course_code string
	var course_name string

	err := db.QueryRow("SELECT courseName, courseCode FROM CourseNames WHERE  programCode = ? ", programCode).Scan(&course_name, &course_code)

	if err != nil {
		fmt.Printf("Program Name Already Exists Changed:: %s", err)
	}

	return err
}

func CreateProgramEntry(programName, programCode string) error {

	db := dbcode.SqlRead().DB

	defer db.Close()

	_, err := db.Exec("INSERT INTO ProgramData (programName, programCode) VALUES (?,?)",
		programName, programCode)

	if err != nil {
		fmt.Println("Failed to add to database")
		return err

	}

	return nil
}



func ProgramsAvailabel() ProgramsAvailabelSt {

	var present bool
	var program_list ProgramsAvailabelSt

	programs_available, err := GetAllProgramData()

	if err != nil {
		present = false

	} else {
		present = true
	}

	program_list = ProgramsAvailabelSt{
		Present:     present,
		ProgramList: programs_available,
	}

	return program_list
}

func GetAllProgramData() ([]ProgramDataEntry, error) {

	db := dbcode.SqlRead().DB
	var programData ProgramDataEntry

	rows, err := db.Query("SELECT programName, programCode FROM ProgramData")
	if err != nil {
		return nil, err

	}
	defer rows.Close()

	var resultList []ProgramDataEntry

	for rows.Next() {
		var program_name string
		var program_code string
		var course_names []Course_Name

		err := rows.Scan(&program_name, &program_code)

		if err != nil {
			fmt.Println("the getall error", err)

			return nil, err
		} else {
			course_names_check, errCourses := GetProgramCourses(program_code)

			if errCourses != nil {

				fmt.Printf("Failed to Get Program Courses: %s", err)
			} else {
				course_names = course_names_check
			}
		}

		programData = ProgramDataEntry{
			Name:        program_name,
			Code:        program_code,
			CourseNames: course_names,
		}

		resultList = append(resultList, programData)

	}

	return resultList, nil
}

///COURSES CRUD

func GetPorgamCourseR(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "programs.html", nil)

	if err != nil {
		log.Fatal(err)
	}

}

func CreateProgramCourseR(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	dbread := dbcode.SqlRead().DB
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	program_code := r.URL.Query().Get("programcode")
	course_name := r.FormValue("course_name")
	course_code := r.FormValue("course_code")

	var setTemplatName string

	err := CheckCourseInDataBase(program_code, course_name, course_code)

	if err != nil {
		fmt.Printf("Failed To Create New Course Because Course With Either The Same Code Or Name Already Exists In The Data Base")
		setTemplatName = "failedToCreateCourse"
	}

	stmt, err := dbread.Prepare("INSERT INTO CourseNames(courseCode,courseCode, programCode) values(?,?,?)")

	if err != nil {
		fmt.Println("Failed to get the email and password", err)
		setTemplatName = "failedToCreateCourse"

	}else{

		err = CreateCourseTable(course_code) 
		if err != nil{
			fmt.Println("Failed to create course table: ", err)
		}
	}

	defer stmt.Close()

	_, err = stmt.Exec(course_name, course_code, program_code)

	if err != nil {
		fmt.Printf("Failed to execute db command create HOW:: %s", err)
		setTemplatName = "failedToCreateCourse"

	}

	err = tpl.ExecuteTemplate(w, setTemplatName, nil)

	if err != nil {
		log.Fatal(err)
	}

}

func CreateCourseMaterial(program_code, course_name, course_code string) error {

	dbread := dbcode.SqlRead().DB

	stmt, err := dbread.Prepare("INSERT INTO CourseNames(courseName,courseCode, programCode) values(?,?,?)")

	if err != nil {
		fmt.Println("Failed to get the email and password", err)

	}

	defer stmt.Close()

	_, err = stmt.Exec(course_name, course_code, program_code)

	if err != nil {
		fmt.Printf("Failed to execute db command create NOW:: %s", err)

	}

	return nil

}

func CheckCourseInDataBase(program_code, course_name, course_code string) error {

	dbread := dbcode.SqlRead().DB
	var course_data Course_Name

	statement, err := dbread.Prepare("SELECT courseName, courseCode,programCode FROM CourseNames = ? WHERE courseName = ? OR courseCode = ?")

	if err != nil {
		return err
	}

	defer statement.Close()

	err = statement.QueryRow(course_name, course_code).Scan(
		&course_data.Name,
		&course_data.Code,
		&course_data.PCode,
	)

	if err != nil {
		return err
	}

	return nil
}



func GetProgramCourses(program_code string) ([]Course_Name, error) {

	db := dbcode.SqlRead().DB
	var courseData Course_Name
	var courseDataListOut []Course_Name

	rows, err := db.Query("SELECT courseName, courseCode FROM CourseNames WHERE programCode = ?", program_code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var course_name string
		var course_code string

		err := rows.Scan(&course_name, &course_code)
		if err != nil {
			return nil, err
		} else {
		}

		courseData = Course_Name{
			Name: course_name,
			Code: course_code,
		}

		courseDataListOut = append(courseDataListOut, courseData)

	}

	return courseDataListOut, nil
}

func CreateCourseAssignmentFolder(cource_namme string) error{


	return nil
	
}

func CreateCourseDB() {
	db := dbcode.SqlRead().DB

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS CourseNames (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		courseName TEXT UNIQUE,
		courseCode TEXT UNIQUE,
		programCode TEXT
	)`)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

}


func CreateNewCourseTable(course_code string) error{

	dbread := dbcode.SqlRead()

	defer dbread.DB.Close()

	stmt_str := fmt.Sprintf(`
		create table if not exists %s (uuid blob not null, 
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
	fmt.Println("the course code",course_code)

	_, err := dbread.DB.Exec(stmt_str)

	if err != nil {
		log.Printf("%q: %s\n", err, stmt_str)
		return err

	}


	return nil

}
