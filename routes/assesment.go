package routes

import (
	"log"
	"fmt"
	"time"
	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"

)

type FileDirectory struct{
	UUID string
	Cource_Name string
	Student_UUID string
	File_Name string
	Date string
}

func LoadExam() {
	dbread := dbcode.SqlRead().DB

	defer dbread.Close()

	//CREATE ADMS
	create_exam := `
		create table if not exists exam_table(
			uuid blod,
			student_uuid blob,
			program_name text,
			cource_name text,
			cource_code text,
			grade text,
			remark text,
			comment text);`

	_, err := dbread.Exec(create_exam)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadAssesments() {
	dbread := dbcode.SqlRead().DB

	defer dbread.Close()

	//CREATE ADMS
	create_table := `
		create table if not exists assesment_table(
			uuid blod,
			student_uuid blob,
			program_name text,
			cource_name text,
			cource_code text,
			grade text,
			remark text,
			comment text);`

	_, err := dbread.Exec(create_table)
	if err != nil {
		log.Fatal(err)
	}

	fileDirectory := `create table if not exists assesmentdirectory(
			uuid blod,
			student_uuid blob,
			cource_name text,
			cource_code text,
			file text,
			date text);`

	_, err = dbread.Exec(fileDirectory)
	if err != nil{
		log.Fatal(err)
	}
}



func ListFileDirectories(student_uuid, cource_name string) (bool, []FileDirectory){


	dbread := dbcode.SqlRead().DB
	stCleaned := CleanStudentUUID(student_uuid)
	present := true
	var fileDir FileDirectory
	var fileDirList []FileDirectory

	stmt, err := dbread.Query("SELECT uuid, cource_name, student_uuid, file, date FROM assesmentdirectory WHERE student_uuid = ? AND cource_name = ?",stCleaned, cource_name )

	if err != nil {
		fmt.Println("QUER STATEMENT FIALED", err)
		present = false


	}

	defer stmt.Close()

	for stmt.Next(){
		err = stmt.Scan(&fileDir.UUID, &fileDir.Cource_Name, &fileDir.Student_UUID, &fileDir.File_Name, &fileDir.Date)
		if err != nil {
			fmt.Println("FAILED TO SCAN FILE")
			present = false
		}
		fileDirList = append(fileDirList, fileDir)
	}


	return present, fileDirList
}



func CreateFileDirectory(cource_name, cource_code,student_uuid, file_name string){
	dbread := dbcode.SqlRead().DB
	uuid := encription.Generateuudi()
	date := time.Now()

	cource_code = SanitizeCookieValue(cource_code)
	

	stCleaned := CleanStudentUUID(student_uuid)
	stmt, err := dbread.Prepare("INSERT INTO assesmentdirectory(uuid, student_uuid, cource_name,course_code,file, date) values(?,?,?,?,?,?)")

	if err != nil {

		fmt.Println("PREPARE STATEMENT FAILED", err)
	}

	defer stmt.Close()

	_,err = stmt.Exec(uuid, stCleaned, cource_name, cource_code,file_name, date)
	if err != nil {
		fmt.Println("FAILED TO CREATE FILE DIRECTORY")
	}
}

func LoadCource() {
	dbread := dbcode.SqlRead().DB

	defer dbread.Close()

	//CREATE ADMS
	create_cource := `
		create table if not exists cource_table(
			uuid blod,
			program_name text,
			cource_name text,
			cource_code,
			cource_assesment text,
			video_list text,
			module text,
			recomended_book text,
			exam_file text);`

	_, err := dbread.Exec(create_cource)
	if err != nil {
		log.Fatal(err)
	}
}
