package routes

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/http"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"
)

const (
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits           = "0123456789"
	specialChars     = "!@#$%^&*()-_=+,.?/:;{}[]~"
)

func generatePassword2(length int) (string, error) {
	allChars := lowercaseLetters + uppercaseLetters + digits + specialChars

	password := make([]byte, length)
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}
		password[i] = allChars[index.Int64()]
	}

	return string(password), nil
}

func CheckUserCridentials(email string) (bool, string) {

	dbread := dbcode.SqlRead().DB
	var email_out string

	present := true

	stmt, err := dbread.Prepare("SELECT email from studentcridentials WHERE email = ? ")

	if err != nil {
		fmt.Println("PREPARE STATEMENT FAILED", err)
		present = false
	}

	defer stmt.Close()

	err = stmt.QueryRow(email).Scan(&email_out)

	if err != nil {
		fmt.Println("FAILED TO READ QueryRow", err)
		present = false

	}

	return present, email_out

}

func UpdateCridentials(email string) {

	dbread := dbcode.SqlRead().DB
	newpassowrd, _ := generatePassword2(8)
	resetSucc := true
	key := encription.GetKey()
	encryptPassword, _ := encription.EncryptData(newpassowrd, key)

	stmt, err := dbread.Prepare("UPDATE studentcridentials SET password = ? WHERE email = ?")

	if err != nil {
		fmt.Println("PREPARE STATEMENT ERROR: ", err)
		resetSucc = false
	}

	defer stmt.Close()

	_, err = stmt.Exec(encryptPassword)

	if err != nil {
		fmt.Println("Execu FAILED ", err)
		resetSucc = false
	}

	ResetMessage(email, newpassowrd)

	fmt.Println(resetSucc)

}
func PasswordResetPage(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "resetpassword.html", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))
	r.ParseForm()

	var templateName string

	emailIn := r.FormValue("email")
	present, _ := CheckUserCridentials(emailIn)

	if present {
		UpdateCridentials(emailIn)
		templateName = "resetsuccessful"
	} else {
		templateName = "emaildoes"

	}
	err := tpl.ExecuteTemplate(w, templateName, nil)

	if err != nil {
		log.Fatal(err)
	}

}
