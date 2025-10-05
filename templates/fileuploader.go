package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type User struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Age      int    `json:"age"`
	ImageURL string `json:"image_url,omitempty"`
}

type Response struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	User    *User  `json:"user,omitempty"`
}

const (
	maxUploadSize = 10 << 20 // 10 MB
	uploadPath    = "./uploads"
)

func main() {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		log.Fatal("Failed to create upload directory:", err)
	}

	// Serve static files (HTML, CSS, JS, uploaded images)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadPath))))
	http.HandleFunc("/", serveHomePage)
	http.HandleFunc("/upload", uploadHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		sendError(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get form values
	name := r.FormValue("name")
	surname := r.FormValue("surname")
	ageStr := r.FormValue("age")

	if name == "" || surname == "" || ageStr == "" {
		sendError(w, "All fields are required", http.StatusBadRequest)
		return
	}

	age, err := strconv.Atoi(ageStr)
	if err != nil || age <= 0 {
		sendError(w, "Invalid age", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("image")
	if err != nil {
		sendError(w, "Failed to get uploaded file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		sendError(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	file.Seek(0, 0) // Reset file pointer

	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		sendError(w, "Only JPEG, PNG and GIF images are allowed", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	uniqueName := generateUniqueFilename(header.Filename)
	filePath := filepath.Join(uploadPath, uniqueName)

	// Create file on server
	dst, err := os.Create(filePath)
	if err != nil {
		sendError(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		sendError(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create user object
	user := &User{
		Name:     name,
		Surname:  surname,
		Age:      age,
		ImageURL: "/uploads/" + uniqueName,
	}

	// Log user data
	log.Printf("User registered: %+v\n", user)

	// Send success response with user data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "User information uploaded successfully!",
		User:    user,
	})
}

func generateUniqueFilename(originalName string) string {
	// Generate random string
	b := make([]byte, 16)
	rand.Read(b)
	randomStr := fmt.Sprintf("%x", b)[:8]
	
	// Get file extension
	ext := filepath.Ext(originalName)
	
	// Add timestamp for uniqueness
	timestamp := time.Now().Unix()
	
	return fmt.Sprintf("%d_%s%s", timestamp, randomStr, ext)
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{Error: message})
}