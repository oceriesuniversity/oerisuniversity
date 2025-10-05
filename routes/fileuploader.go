package routes

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"fmt"
)


func AdminDownLoadAsignment(w http.ResponseWriter, r *http.Request){

	uuid := r.URL.Query().Get("student_uuid")
	cource_name := r.URL.Query().Get("cource_name")
	file_name := r.URL.Query().Get("file_name")
	
	fmt.Println(uuid)
	fmt.Println(cource_name)
	fmt.Println(file_name)



	dbFilePath := fmt.Sprintf("assesmentFiles/%s/%s/%s", uuid, cource_name, file_name)


	fmt.Println("The Download path",dbFilePath)

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

func ReadStudentAssesments(uuid, cource_name string) []string{
	

	var filesList []string

	stID := CleanStudentUUID(uuid)

	parts := []string{"assesmentFiles", stID, cource_name}
	result := strings.Join(parts, "/")

	
	files, err := os.ReadDir(result)

	if err != nil {
		fmt.Println("FAILED TO READ FILES", err)
		
	}

	for _, file :=  range files{
		filesList = append(filesList, file.Name())
	}


	return filesList

}


func UploadAssesment(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form.
	err := r.ParseMultipartForm(10 << 20) 
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		return
	}

	
	student_uuid := r.URL.Query().Get("uuid")
	cource_name := r.URL.Query().Get("cource_name")
	cource_code := r.URL.Query().Get("course_code")


	stID := CleanStudentUUID(student_uuid)
	

	file, handler, err := r.FormFile("file")

	fileName := handler.Filename

	CreateFileDirectory(cource_name, cource_code,student_uuid, fileName)


	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure the file has a .pdf extension.
	if filepath.Ext(handler.Filename) != ".pdf" {
		http.Error(w, "Only PDF files are allowed", http.StatusBadRequest)
		return
	}

	// Create the destination file.

	parts := []string{"assesmentFiles", stID, cource_name}
	result := strings.Join(parts, "/")

	savePath := filepath.Join(result, handler.Filename)
	err = os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	destFile, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	// Copy the uploaded file to the destination file.
	_, err = io.Copy(destFile, file)
	if err != nil {
		http.Error(w, "Failed to copy file", http.StatusInternalServerError)
		return
	}

	// Respond to the client.

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, "fileuploaded", nil)

	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
}