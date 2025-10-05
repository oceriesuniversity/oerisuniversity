package routes 



import (
	"log"
	"fmt"
	"html/template"
	"net/http"
	"github.com/kamukwamba/oerisuniversity/dbcode"
	"strings"

)


func UpdateAllProgramDetailsR(w http.ResponseWriter, r *http.Request){

	tpl = template.Must(template.ParseGlob("templates/*.html"))


	program_code := r.URL.Query().Get("programCode")
	program_name := r.URL.Query().Get("programName")

	programData := ProgramDataEntry{
		Name: program_name,
		Code: program_code,
	}

	fmt.Println("Name: ", program_name, "\n Code:", program_code)
	
	err := tpl.ExecuteTemplate(w, "updateprgogramcard", programData)

	if err != nil {
		log.Fatal(err)
	}

}


func UpdateAllProgramDetails(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	old_program_code := r.URL.Query().Get("oldprogramcode")

	new_program_code := r.FormValue("program_code")
	program_name := r.FormValue("program_name")


	
	err := UpdateTheCourceNamesTable(old_program_code, new_program_code)
	if err != nil{
		fmt.Println("Calling the fuction returned a failed code")
	}

	err = ChangeProgramDataName(old_program_code, new_program_code, program_name)

	if err != nil{
		fmt.Println("Failed to change the Program Data Names")
	}

	err = ChangeCourceTable(old_program_code, new_program_code)

	if err != nil {
		fmt.Println("Failed to update cource table")
	}


	err = UpdateStudentProgram(old_program_code, new_program_code)

	if err != nil{
		fmt.Println("Failed to update student program list")
	}

	err = ChangeTableName(old_program_code, new_program_code)
	if err != nil {
		fmt.Println("Failed to add to new data base")
	}





}


func UpdateItOut(uuid, old_code, new_code string,index int, program_list []string)error{

	new_list := program_list
	new_list[index] = new_code


	dbread := dbcode.SqlRead().DB

	stmt, err := dbread.Prepare("UPDATE studentprogramlist SET program_list = ? WHERE uuid = ?")

	if err != nil {
		fmt.Println("Failed to prepare query for studentprogramlist")
	}

	defer stmt.Close()

	_, err_out := stmt.Exec(new_list, uuid) 

	if err_out != nil {
		fmt.Println("Failed to update student program list")
	}

	return nil

}

func UpdateStudentProgram(old_program_code, new_program_code string) error{

	student_data_out := GetStudentProgramsAll()
	var uuid string

	if len(student_data_out) != 0{
	for _, item := range student_data_out{
		uuid = item.UUID
		student_programs := item.Program_List

		for i, item_p := range student_programs{
			if item_p == old_program_code{
				index := i
				err := UpdateItOut(uuid, old_program_code, new_program_code, index, student_programs)
				if err != nil {
					fmt.Println("Failed to update the student program data")
				}

			}else{
				continue
			}
		}
	}	
	}else{
		fmt.Println("No Students in list")
	}


	return nil

}


type StudentPList struct{
	UUID string
	Program_List []string
}

func GetStudentProgramsAll() []StudentPList {
	dbread := dbcode.SqlRead().DB

	var uuid string
	var program_list string

	var listout []string
	var update_list []StudentPList

	var student_p_data StudentPList

	stmt, err := dbread.Query("select uuid, program_list from studentprogramlist")

	if err != nil {
		fmt.Println("Failed to get from studentprogramlist")
	}

	defer stmt.Close()

	for stmt.Next(){
		err = stmt.Scan(&uuid, &program_list)
		trimedlist := strings.Trim(program_list, "[]")

		list_out := strings.Split(trimedlist, ",")

		for _, item := range list_out {
			trimedlistone := strings.Trim(item, "\"")
			trimedlisttwo := strings.Trim(trimedlistone, " \"")

			if len(trimedlisttwo) > 1 {
				listout = append(listout, trimedlisttwo)

				student_p_data = StudentPList{
					UUID: uuid,
					Program_List: listout,
				}

				update_list = append(update_list, student_p_data )


			} else {
				continue
			}
		}

		if err != nil {
		fmt.Println("FAILED TO GET STUDENT PROGRAM LIST")
		}
	}

	


	return update_list
}



func ChangeTableName(old_program_code,  new_program_code string) error{


	dbread := dbcode.SqlRead().DB


	query := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", old_program_code, new_program_code)
    
    _, err := dbread.Exec(query)
    if err != nil {
        return fmt.Errorf("failed to rename table: %v", err)
    }
    
    fmt.Printf("Table renamed from %s to %s\n", old_program_code, new_program_code)
    return nil


}



// func ProgramExist(program_code, table_name string) bool{
// 	var exists bool
// 	db := dbcode.SqlRead().DB

// 	defer db.Close()



//     err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", program_code).Scan(&exists)
//     if err != nil {
//         return fmt.Println("error checking existence: %v", err)
//     }


    
//     if !exists {
//         return fmt.Println("user with id %d does not exist", program_code)
//     }

//     return exists

// }


//
func UpdateTheCourceNamesTable(old_program_code, new_program_code string) error{

	dbread := dbcode.SqlRead().DB

	defer dbread.Close()

	
	stmt, err:= dbread.Prepare("UPDATE CourseNames SET programCode = ? WHERE programCode = ?")
	if err != nil {
		fmt.Println("Failed to Update the program name for the cource names: ", err)
	}

	_, err = stmt.Exec(new_program_code, old_program_code)

	if err != nil {
		fmt.Println("Failed to update the Cource Names Table")
		return err
	}


	return nil

}

//

func ChangeProgramDataName(old_program_code, new_program_code, new_program_name string) error{



	dbread := dbcode.SqlRead().DB

	defer dbread.Close()


	stmt, err:= dbread.Prepare("UPDATE ProgramData SET programName = ?,programCode = ? WHERE programCode = ?")
	if err != nil {
		fmt.Println("Failed to Update the program name for the ProgramData: ", err)
		return err
	}

	_, err = stmt.Exec(new_program_name, new_program_code, old_program_code)

	if err != nil {
		fmt.Println("Failed to update the ProgramData Table")
		return err
	}

	return nil
}

//
func ChangeCourceTable(old_program_code, new_program_code string) error{


	dbread := dbcode.SqlRead().DB

	defer dbread.Close()

	// program_present := ProgramExist(old_program_code)

	// if program_present{
	stmt, err:= dbread.Prepare("UPDATE cource_table SET program_name = ? WHERE program_name = ?")
	if err != nil {
		fmt.Println("Failed to Update the cource_table for the cources: ", err)
		return err
	}

	_, err = stmt.Exec(new_program_code, old_program_code)

	if err != nil {
		fmt.Println("Failed to update the cource_table Table")
		return err
	}


	// }else{
	// 	fmt.Println("No program in the cource table")
	// }

	return nil
}

