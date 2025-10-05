package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"
)

type MessageAdmin struct {
	StInfo   StudentInfo
	Messages []MessageOut
}

type MessageOut struct {
	UUID         string
	Student_UUID string
	Sender_Name  string
	Message      string
	Sender       bool
	Seen_Student bool
	Seen_Admin   bool
	Date         string
}



func DeleteMessage(message_uuid string) bool {

	dbdelete := dbcode.SqlRead().DB

	deleted := true

	delete, err := dbdelete.Prepare("DELETE FROM messages WHERE uuid = ?")
	if err != nil {
		errout := fmt.Sprintf("line 30 erro: %s", err)
		ErrorPrintOut("messages", "DeleteMessages", errout)
	}
	defer delete.Close()

	_, errd := delete.Exec(message_uuid)
	if errd != nil {
		errout := fmt.Sprintf("line 42 erro: %s", err)
		ErrorPrintOut("messages", "DeleteMessages", errout)
	}

	return deleted

}

func UpdateMessages(student_uuid string) bool {
	dbupdate := dbcode.SqlRead().DB

	updated := true

	update, err := dbupdate.Prepare("UPDATE messages SET seen_admin = ? WHERE sender_uuid = ? ")

	if err != nil {
		errout := fmt.Sprintf("line 30 erro: %s", err)
		ErrorPrintOut("messages", "UpdateMessages", errout)
	}

	defer update.Close()

	_, errupdate := update.Exec(true, student_uuid)

	if errupdate != nil {
		errout := fmt.Sprintf("line 48 error: %s", err)
		ErrorPrintOut("messaes", "UpdateMessages", errout)
	}

	return updated

}

func ReadMessage(student_uuid string) MessageAdmin {
	dbread := dbcode.SqlRead().DB
	var message_out MessageOut

	student_data := GetStudentAllDetails(student_uuid)

	var message_out_list []MessageOut

	statement, err := dbread.Query("select uuid,sender_uuid,sender_name,sender,message,seen_student,seen_admin,date from messages")

	if err != nil {
		errorout := fmt.Sprintf("line 22: %s", err)
		ErrorPrintOut("messasges.go", "ReadMessages", errorout)
	}

	defer statement.Close()

	for statement.Next() {
		err = statement.Scan(&message_out.UUID, &message_out.Student_UUID, &message_out.Sender_Name, &message_out.Sender, &message_out.Message, &message_out.Seen_Student, &message_out.Seen_Admin, &message_out.Date)

		if err != nil {
			errorout := fmt.Sprintf("line 48: %s", err)
			ErrorPrintOut("messages.go", "ReadMessages", errorout)
		}

		if message_out.Student_UUID == student_uuid {
			message_out_list = append(message_out_list, message_out)

		}

	}

	student_msg := MessageAdmin{
		StInfo:   student_data,
		Messages: message_out_list,
	}

	fmt.Println(student_msg)

	return student_msg

}

func LoadMessages() []MessageOut {
	var message_out_list []MessageOut
	var message_out MessageOut

	var message_seen []MessageOut
	var message_unseen []MessageOut

	dbread := dbcode.SqlRead().DB

	statement, err := dbread.Query("select uuid,sender_uuid,sender_name, sender, message,seen_admin,date from messages")

	if err != nil {
		log.Fatal(err)
	}

	defer statement.Close()

	for statement.Next() {
		err := statement.Scan(&message_out.UUID, &message_out.Student_UUID, &message_out.Sender_Name, &message_out.Sender, &message_out.Message, &message_out.Seen_Admin, &message_out.Date)

		if err != nil {
			error_out := fmt.Sprintf("line 42: %s", err)
			ErrorPrintOut("messages", "LoadMessages", error_out)
		}

		if message_out.Sender {
			if !message_out.Seen_Admin {
				if len(message_unseen) > 0 {
					for _, message := range message_unseen {
						if message_out.Student_UUID != message.Student_UUID {
							message_unseen = append(message_unseen, message_out)
						}
					}
				} else {
					message_unseen = append(message_unseen, message_out)
				}

			} else {
				if len(message_seen) > 0 {
					for _, message := range message_seen {
						if message_out.Student_UUID != message.Student_UUID {
							message_seen = append(message_seen, message_out)
						}
					}
				} else {
					message_seen = append(message_seen, message_out)
				}
			}

		}

		if len(message_unseen) > 0 {
			for _, message := range message_unseen {
				if len(message_out_list) > 0 {
					for _, messageout := range message_out_list {
						if message.Student_UUID != messageout.Student_UUID {
							message_out_list = append(message_out_list, message)
						}
					}
				} else {
					message_out_list = append(message_out_list, message)
				}
			}
		} else {
			for _, message := range message_unseen {
				message_out_list = append(message_out_list, message)
			}
		}

		if len(message_seen) > 0 {
			for _, message := range message_seen {
				if len(message_out_list) > 0 {
					for _, messageout := range message_out_list {
						if messageout.Student_UUID != message.Student_UUID {
							message_out_list = append(message_out_list, message)
						}
					}
				} else {
					message_out_list = append(message_out_list, message)
				}
			}
		} else {
			for _, message := range message_seen {
				message_out_list = append(message_out_list, message)
			}
		}
	}

	fmt.Println(message_out_list)

	return message_out_list

}

type LoadMsg struct {
	Admin AdminInfo
	Msg   []MessageOut
	Admin_Name string
}

func AdminMessagesPage(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))


	user_name, err := GetUserName(r)

	messages_out := LoadMessages()

	display_data := LoadMsg{
		Admin_Name: user_name,
		Msg:   messages_out,
	}


	err = tpl.ExecuteTemplate(w, "messageAdmin.html", display_data)

	if err != nil {
		log.Fatal(err)
	}
}




func GetUUIDList()[]StudentInfo{

	var data StudentInfo
	var data_list []StudentInfo

	dbread := dbcode.SqlRead().DB

	query := `SELECT uuid, first_name, last_name FROM studentdata`
	stmt, err := dbread.Query(query)


	if err != nil {
	    fmt.Println("Failed to prepare statement:", err)
	}
	defer stmt.Close()


	if err != nil {
	    fmt.Println("Failed to execute query:", err)
	}

	for stmt.Next(){
        err = stmt.Scan(&data.UUID, &data.First_Name, &data.Last_Name) 
        if err != nil{
        	fmt.Println("Failed to scan stmt")
        }

        data_list = append(data_list, data)
    }


    return data_list


}

func SendBulleting(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	dbread := dbcode.SqlRead().DB

	defer dbread.Close()

	message := r.FormValue("message_bulleting")

	data_list := GetUUIDList()
	uuid := encription.Generateuudi()
	dateout := fmt.Sprintf("%s", time.Now())

	for _, data := range data_list{

		sender_name := fmt.Sprintf("%s %s", data.First_Name, data.Last_Name)
		student_uuid := data.UUID
		sender := "admin"


		stmt, err := dbread.Prepare("insert into messages (uuid, sender_uuid,sender_name, sender,message,seen_student,seen_admin,date) values(?,?,?,?,?,?,?,?)")

		if err != nil {
			fmt.Println("Faild to call prepare statement")
		}

		defer stmt.Close()

		_, err = stmt.Exec(uuid, student_uuid,sender_name, sender, message, false, true, dateout)
		if err != nil {
			fmt.Println("Failed to send bulleting")
		}
	}
	
}



func GetAllStudentMsg(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	var msg MessageOut
	var msgLs []MessageOut
	var uuid string

	r.ParseForm()

	email := r.FormValue("email")
	dbread := dbcode.SqlRead().DB

	stmt, err := dbread.Prepare("SELECT uuid FROM studentdata WHERE email = ?")

	if err != nil {
		fmt.Println("FAILED TO PREPARE STATEMENT", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(email).Scan(&uuid)
	if err != nil {
		fmt.Println("FAILED TO QUERYROW: ", err)
	}

	stmt_two, err_two := dbread.Query("SELECT select uuid,sender_uuid,sender_name, sender, message,seen_admin,date FROM messages WHERE sender_uuid = ?", uuid)

	if err_two != nil {
		fmt.Println("PREPARE STATEMENT FAILED: ", err_two)
	}

	defer stmt_two.Close()

	for stmt_two.Next() {
		err_two = stmt_two.Scan(&msg.UUID, &msg.Student_UUID, &msg.Sender_Name, &msg.Message, &msg.Sender, &msg.Seen_Student, &msg.Seen_Admin, &msg.Date)

		if err_two != nil {
			fmt.Println("FAILED TO SCAN: ", err_two)
		}

		msgLs = append(msgLs, msg)
	}

	massageData := LoadMsg{
		Msg: msgLs,
	}

	err = tpl.ExecuteTemplate(w, "messageslog", massageData)

	if err != nil {
		log.Fatal(err)
	}
}
