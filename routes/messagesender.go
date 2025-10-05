package routes

import (
	
	"gopkg.in/gomail.v2"
	"fmt"
	"html/template"
	"log"
	"net/http"
	

	"github.com/kamukwamba/oerisuniversity/dbcode"

	"github.com/kamukwamba/oerisuniversity/encription"
)



func CreateEmailSenderTem() {

	dbconn := dbcode.SqlRead().DB
	
	defer dbconn.Close()
	
	tabelCreate := `create table if not exists messengerData(uuid blob not null, email text, password text)`
	
	_, err := dbconn.Exec(tabelCreate)
	if err != nil {
		log.Printf("%q: %s\n", err, tabelCreate)
	}

}



type ApplicationApprovedSender struct {
	Email string
	Password string
	UUID string

}



func GetEmailData() (ApplicationApprovedSender, bool){

	dbcone := dbcode.SqlRead().DB
	
	isPresent := true
	var email string 
	var password  string 
	var uuid string
	
	
	
	stmt, err := dbcone.Prepare("select uuid,email,password from messengerData limit 1")
	
	if err != nil{
	fmt.Println("Failed to get the email and password: ", err)
	isPresent = false
	}
	
	defer stmt.Close()
	
	err = stmt.QueryRow().Scan(&uuid, &email, &password)
	
	if err != nil {
		fmt.Print("Failed to execute db command there are no files: ", err)
		isPresent = false
	
	}
	
	data_out := ApplicationApprovedSender{
		Email: email,
		Password: password,
		UUID: uuid,
		
	}
	
	return data_out, isPresent

}

func CreateEmailData(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	
	email := r.FormValue("email")
	password := r.FormValue("password")
	uuid := encription.Generateuudi()
	
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	dbconn :=  dbcode.SqlRead().DB
	
	//Check if there is an email alread present
	
	
	_, isPresent := GetEmailData()
	
	fmt.Println(isPresent)
	
	
	if isPresent ==  false{
	
		stmt, err := dbconn.Prepare("insert into messengerData(uuid,email, password) values(?,?,?)")
	
		if err != nil{
		fmt.Println("Failed to get the email and password", err)
		}
		
		defer stmt.Close()
		
		_, err = stmt.Exec(uuid,email, password)
		
		if err != nil {
		fmt.Print("Failed to execute db command create::", err)
		}
		
		
		data_out := ApplicationApprovedSender{
				UUID: uuid,
				Email: email,
				Password: password,
		
		}
	
		
		//debug failure to laod templates

		err = tpl.ExecuteTemplate(w, "schoolemailCreate", data_out)

		if err != nil {
			log.Fatal(err)
		}
	
		
		
	
	}else{
		err := tpl.ExecuteTemplate(w, "schoolemailpresent", nil)

		if err != nil {
			log.Fatal(err)
		}
		
		
		
	}
	
	
	

}


func DeleteEmail(w http.ResponseWriter, r *http.Request){

	uuid :=  r.URL.Query().Get("uuid")

	fmt.Println("UUID OUT: ", uuid)
	
	dbconn := dbcode.SqlRead().DB
	stmt, err := dbconn.Prepare("Delete from  messengerData where uuid = ?")
	
	if err != nil{
	fmt.Println("Failed to get the email and password")
	}
	
	defer stmt.Close()
	
	_, err = stmt.Exec(uuid)
	
	if err != nil {
	fmt.Print("Failed to execute db command")
	}
	
}


func SendMsgToAdminEmail(email string) {


	
}

func SendToAdmin(stemail string){
	
	
	
	var email string
	
	dbconn := dbcode.SqlRead().DB
	
	stmt, err := dbconn.Prepare("select email from admin where auth = ?")
	
	if err != nil {
	fmt.Println("Failed to get email from admi, error out: ",err)
	}
	
	defer stmt.Close()
	
	err = stmt.QueryRow("admin").Scan(&email)
	
	if err != nil {
		fmt.Println("Failed to QueryRow error out: ", err)
		
	}
	
	
	
	
	data_out, _ := GetEmailData()
	
	from := data_out.Email 
	password :=  data_out.Password
	
	
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	subject := "STUDENT APPLICATION SUCCESFUL"
	body := fmt.Sprintf("STUDENT APPLICATION WAS SUCCESFUL, STUDENT EMAILL: %s", stemail)
	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Create a new SMTP dialer
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// Send the email
	if err = d.DialAndSend(m); err != nil {
		fmt.Println("Failed to send email error out: ", err)
	}

	fmt.Println("Email sent successfully!")
	

}

func SendEmail(to string) error {
	// Email configuration
	
	
	data_out, _ := GetEmailData()
	
	from := data_out.Email 
	password :=  data_out.Password
	
	
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	subject := "Application Successfull"
	body := fmt.Sprintf("Congratulations your application was successfull your username: %s and password: %s", to,to)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Create a new SMTP dialer
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// Send the email
	err := d.DialAndSend(m)

	if err != nil{
		fmt.Println("FAILED TO SEND EMAIL CHECK \"messagesender.go\" ", err)
	}
		
	
	
	
	SendToAdmin(to)

	fmt.Println("Email sent successfully!")
	return nil
}


func ResetMessage(newpassowrd, to string) error{
	// Email configuration
	
	
	data_out, _ := GetEmailData()
	
	from := data_out.Email 
	password :=  data_out.Password
	
	
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	subject := "Application Successfull"
	body := fmt.Sprintf("Congratulations password reset  successfull. Your new password is : %s", newpassowrd)
	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Create a new SMTP dialer
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	
	SendToAdmin(to)

	fmt.Println("Email sent successfully!")
	return nil
}

