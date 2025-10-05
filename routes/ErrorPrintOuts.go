package routes

import (
	"fmt"
	"net/http"
)

func ErrorPrintOut(file_name, function, the_error string) {

	error_out_string := fmt.Sprintf("FILE NAME:  %s \nFUNCTION NAME: %s \nTHE ERROR: %s", file_name, function, the_error)
	fmt.Println(error_out_string)
}

func ErrorPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("An error occurred while processing your request."))
}
