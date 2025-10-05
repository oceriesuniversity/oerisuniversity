package routes

import (
	"log"

	"crypto/rand"
	"fmt"
	"html/template"
	"math/big"
	"net/http"
	"net/smtp"
	"os"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	// "github.com/kamukwamba/oerisuniversity/dbcode"
)

func ForgotPassword(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "forgotpassword.html", nil)

	if err != nil {
		log.Fatal(err)
	}

}

func ConfirmStudentId(w http.ResponseWriter, r *http.Request) {

	student_email := r.FormValue("studentemail")

	dbread := dbcode.SqlRead()
	var redirectName string

	stmt, err := dbread.DB.Prepare("select email from studentdata where email = ?")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	var email_out string

	err = stmt.QueryRow(student_email).Scan(&email_out)

	if err != nil {
		fmt.Println("Email Dose Not Exist")
	} else {
		errOut := CreateUpdatePassword(email_out)
		if errOut {
			redirectName = "login_div_failed"

		} else {
			redirectName = "login_div_succesfull"
		}

	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, redirectName, nil)

	if err != nil {
		log.Fatal(err)
	}

}

// Create And  Update password
func CreateUpdatePassword(email string) bool {
	dbread := dbcode.SqlRead().DB
	length := 12 // Default password length
	errorOut := false

	newpassowrd := generatePassword(length)
	encryptPassword, err := HashPassword(newpassowrd)

	statement := fmt.Sprintf(`UPDATE studentdata SET password = %s WHERE email = %s `, encryptPassword, encryptPassword)

	_, err = dbread.Prepare(statement)

	if err != nil {
		fmt.Println("Error Failed to Save Reset Password")

		errorOut = true
	} else {
		errReset := SendPasswordResetEmail(email, newpassowrd)
		if errReset != nil {
			errorOut = true
		}
	}

	return errorOut
}

func generatePassword(length int) string {
	// All possible characters to use in password
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789" +
		"!@#$%^&*()-_=+,.?/:;{}[]~"

	password := make([]byte, length)
	for i := 0; i < length; i++ {
		// Get a random index from the chars string
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err) // In a real app, handle this more gracefully
		}
		password[i] = chars[num.Int64()]
	}

	return string(password)
}

func SendPasswordResetEmail(to, newpassword string) error {
	from := os.Getenv("SMTP_FROM")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	subject := "Password Reset Request"
	body := fmt.Sprintf(`
		Hello,
		
		We received a request to reset your password. Click the link below to proceed:
		
		%s
		
		This link will expire in 1 hour. If you didn't request this, please ignore this email.
		
		Thanks,
		Oceris Team
	`, newpassword)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}
