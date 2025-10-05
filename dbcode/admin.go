package dbcode

import (
	"database/sql"
	"os"
	"fmt"
	"log"
	"io"
	"path/filepath"
	"net/http"
)

func BackUpData(w http.ResponseWriter, r *http.Request){

	
	dbFilePath := "./data/ucms.db"

	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		http.Error(w, "Database file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(dbFilePath)))

	
	file, err := os.Open(dbFilePath)
	if err != nil {
		http.Error(w, "Unable to open the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error writing file to response", http.StatusInternalServerError)
		return
	}
	
}


type AdminInfo struct {
	ID         string
	First_Name string
	Last_Name  string
	Email      string
	Password   string
	Date       string
}

var inforOutLsit []AdminInfo

func AdminGet() []AdminInfo {

	dbread := SqlRead()
	var infor_out AdminInfo
	rows, err := dbread.DB.Query("select uuid, first_name, last_name,email, password, auth, date from admin")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var first_name string
		var last_name string
		var email string
		var auth string
		var date string
		var password sql.NullString
		err = rows.Scan(&id, &first_name, &last_name, &email, &password, &auth, &date)

		infor_out = AdminInfo{
			ID:         id,
			First_Name: first_name,
			Last_Name:  last_name,
			Email:      email,
			Password:   password.String,
		}
		inforOutLsit = append(inforOutLsit, infor_out)
		if err != nil {
			log.Fatal(err)
		}

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return inforOutLsit
}

func CreateAdmin(info AdminInfo) {
	dbread := SqlRead()
	tx, err := dbread.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into admin(uuid, admin_name, admin_email, admin_password) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちは世界%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func LoadDB() {

	dbread := SqlRead()

	defer dbread.DB.Close()

	course_data := `
		create table if not exists coursedata(
		uuid blob not null,
		program_name text,
		course_name text,
		course_recomended_book text,
		module text)

	`

	course_data_videos := `
		create  table if not exists coursedatavideos( 
		uuid blob not null,
		cource_uuid blob,
		video_link text)
		`
	studentprogramlist :=
		`create table if not exists studentprogramlist(
		uuid blob not null,
		student_uuid blob,
		program_list blob
	)`
	studentcridentials :=
		`create table if not exists studentcridentials(
		uuid blob not null,
		student_uuid blob,
		email text,
		password text
	)`

	studentsdata := `create table if not exists studentdata( 
												uuid blob not null, 
												first_name text, 
												last_name text,
												phone text,
												email text, 
												date_of_birth text, 
												gender text,
												marital_status text, 
												country text, 
												eduction_background text, 
												program text, 
												high_scholl_confirmation text,
												grammer_comprihention text, 
												waiver text, 
												number_of_children text,
												school_atteneded text, 
												major_studied text, 
												degree_obtained text, 
												current_occupetion text,
												field_interested_in text, 
												mps_techqnique_Practiced text, 
												previouse_experince text, 
												purpose_of_enrollment text, 
												use_of_degree text, 
												reason_for_choice text, 
												method_of_incounter text,
												student_id text);
		`

	programvideosname := `create table if not exists programvideos ( 
					uuid blob,
			 		Cource_name text, 
					video_link_dic text);`

	messages := `create table if not exists messages(
		uuid blob,
		sender_uuid blob,
		sender_name text,
		sender bool,
		message text,
		seen_student bool,
		seen_admin bool,
		date string
	)`

	news_story := `
		create table if not exists news(
		uuid blob,
		title text,
		auther text, 
		image text,
		story text,
		date text)

	`

	lectureras := `
		create table if not exists lectureras(
		uuid blob not null,
		date text,
		first_name text,
		last_name text,
		emial text,
		password,
		permision text)

	`

	_, lectureraserror := dbread.DB.Exec(lectureras)
	if lectureraserror != nil {
		log.Printf("%q: %s\n", lectureraserror, lectureras)
	}

	_, newsstoryerror := dbread.DB.Exec(news_story)
	if newsstoryerror != nil {
		log.Printf("%q: %s\n", newsstoryerror, news_story)
	}

	_, msgerror := dbread.DB.Exec(messages)
	if msgerror != nil {
		log.Printf("%q: %s\n", msgerror, course_data)

	}
	_, errcourse_data := dbread.DB.Exec(course_data)

	if errcourse_data != nil {
		log.Printf("%q: %s\n", errcourse_data, course_data)
	}

	_, errcourse_data_videos := dbread.DB.Exec(course_data_videos)

	if errcourse_data_videos != nil {
		log.Printf("%q: %s\n", errcourse_data_videos, course_data_videos)
	}
	_, errvideos := dbread.DB.Exec(programvideosname)

	if errvideos != nil {
		log.Printf("%q: %s\n", errvideos, programvideosname)
	}

	_, errstp := dbread.DB.Exec(studentprogramlist)
	if errstp != nil {
		log.Printf("%q: %s\n", errstp, studentprogramlist)
		return
	}

	_, errstc := dbread.DB.Exec(studentcridentials)
	if errstc != nil {
		log.Printf("%q: %s\n", errstc, studentcridentials)
		return
	}

	_, errstd := dbread.DB.Exec(studentsdata)

	if errstd != nil {
		log.Printf("%q: %s\n", errstd)
		return
	}

	defer dbread.DB.Close()
}
