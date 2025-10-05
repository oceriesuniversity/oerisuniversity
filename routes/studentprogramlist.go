package routes

import (
	"fmt"
	"log"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"
)

func AddStudentPrograms(studentuuid, program_code string) {
	var programlistname StringSlice

	uuid := encription.Generateuudi()

	programlistname = append(programlistname, program_code)

	dbread := dbcode.SqlRead()
	program_name_list, err := dbread.DB.Begin()
	if err != nil {
		log.Fatal()
	}

	stmt, err := program_name_list.Prepare("insert into studentprogramlist(uuid, student_uuid, program_list) values(?,?,?)")

	if err != nil {
		log.Fatal("Failed to add student to studentprogramlist :AddStudentPrograms",err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(uuid, studentuuid, programlistname)

	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to create Stduent program list")
	}

	err = program_name_list.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("STUDENT PROGRAM LIST CREATED SUCCESFULLY")

}

func AddToProgramList(program_name, student_uuid string) bool {
	updated := true

	// .dbread"UPDATE artist_t SET check_s = ? WHERE artist_n = ?", "2021-05-20", 42

	fmt.Println("The Program Out In: ", program_name)

	programlist := GetStudentPrograms(student_uuid)

	AppendProgramList(program_name, student_uuid, programlist)

	return updated

}

func AppendProgramList(program_name, student_uuid string, program_list []string) bool {

	// dbread := dbcode.SqlRead()
	fmt.Println("The Program Name: ", program_name)
	var updated bool = true
	var program_list_in_one []string

	for _, item := range program_list {
		if len(item) > 1 {
			program_list_in_one = append(program_list_in_one, item)
		} else {
			continue
		}
	}

	program_list_in_one = append(program_list_in_one, program_name)

	dbupdate := dbcode.SqlRead().DB

	var list_of_strings []string
	count := 1

	fmt.Println("the student uuid: ", student_uuid)
	for _, item := range program_list_in_one {

		if count == len(list_of_strings) {
			out := fmt.Sprintf("\"%s\"", item)
			list_of_strings = append(list_of_strings, out)
		} else {
			out := fmt.Sprintf("\"%s\",", item)
			list_of_strings = append(list_of_strings, out)
		}

		count += 1

	}

	programlist := fmt.Sprintf("%s", list_of_strings)

	stmt, err := dbupdate.Prepare("UPDATE studentprogramlist SET program_list = ? WHERE student_uuid = ?")
	fmt.Println("400", student_uuid)

	if err != nil {

		log.Fatal(err)
	}

	defer stmt.Close()

	_, erre := stmt.Exec(programlist, student_uuid)

	if erre != nil {
		log.Fatal(erre)

	}

	return updated

}

func GetStudentProgramData(programlist []string, students_uuid string) ([]AllCourceData, bool) {

	var programdata ProgramStruct

	var courcedata []CourceStruct
	var allcourcedataout AllCourceData

	var allcourcedataoutlist []AllCourceData

	var programsavailable ProgramAvailable
	var available []bool

	cunt := 0
	for _, program := range programlist {
		cunt += 1
		fmt.Println("Count Print", cunt)

		is_present, dataout, _ := GetProgramsStudents(students_uuid, "one", program)

		if is_present {
		

			var programdataacams StudentProgramData = dataout

			programdata.UUID = programdataacams.UUID
			programdata.Student_UUID = programdataacams.Student_UUID
			programdata.Program_Name = programdataacams.Program_Name
			programdata.First_Name = programdataacams.First_Name
			programdata.Last_Name = programdataacams.Last_Name
			programdata.Email = programdataacams.Email
			programdata.Payment_Method = programdataacams.Payment_Method
			programdata.Paid = programdataacams.Paid
			programdata.Approved = programdataacams.Approved
			programdata.Applied = programdataacams.Applied
			programdata.Completed = programdataacams.Completed
			programdata.Date = programdataacams.Date

			courcedata = GetFromProgramCources(students_uuid, program)

			allcourcedataout.ProgramStruct = programdata
			allcourcedataout.Cource_Struct = courcedata

			allcourcedataoutlist = append(allcourcedataoutlist, allcourcedataout)

			available = append(available, true)
		} else {
			available = append(available, false)
		}

	}

	programsavailable.Available = available[0]

	// fmt.Println("There are programs: ", available[0])
	// fmt.Println("Student data out 2 how is it working: ", allprogramdata)
	return allcourcedataoutlist, available[0]

}
