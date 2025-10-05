package routes


import (
	"encoding/json"

	"github.com/kamukwamba/oerisuniversity/dbcode"
)

type Course struct {
	ProgramCode    string   `json:"program_code"`
	CourseName     string   `json:"course_name"`
	CourseCode     string   `json:"course_code"`
	ModuleLink     string   `json:"module_link"`
	VideoLinks     []string `json:"video_links"`
	AssessmentLink string   `json:"assessment_link"`
}

func CreateTableCourseDetails() error {

	db := dbcode.SqlRead().DB

	query := `
	CREATE TABLE IF NOT EXISTS coursematerial (
		program_code TEXT PRIMARY KEY,
		course_name TEXT NOT NULL,
		course_code TEXT NOT NULL,
		module_link TEXT,
		video_links TEXT,
		assessment_link TEXT
	)`
	_, err := db.Exec(query)
	return err
}

func InsertCourseDetails(course Course) error {
	// Convert video links to JSON
	db := dbcode.SqlRead().DB
	videoLinksJSON, err := json.Marshal(course.VideoLinks)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO coursematerial (
		program_code, 
		course_name, 
		course_code,
		module_link, 
		video_links, 
		assessment_link
	) VALUES (?, ?, ?, ?, ?, ?)`

	_, err = db.Exec(
		query,
		course.ProgramCode,
		course.CourseName,
		course.CourseCode,
		course.ModuleLink,
		string(videoLinksJSON),
		course.AssessmentLink,
	)
	return err
}

func GetCourseDetails(courseCode string) (Course, error) {
	var course Course
	var videoLinksJSON string
	db := dbcode.SqlRead().DB

	query := `
	SELECT 
		program_code, 
		course_name, 
		module_link, 
		video_links, 
		assessment_link 
	FROM courses 
	WHERE program_code = ?`

	err := db.QueryRow(query, courseCode).Scan(
		&course.ProgramCode,
		&course.CourseName,
		&course.ModuleLink,
		&videoLinksJSON,
		&course.AssessmentLink,
	)
	if err != nil {
		return course, err
	}

	// Convert JSON video links back to slice
	err = json.Unmarshal([]byte(videoLinksJSON), &course.VideoLinks)
	return course, err
}

func UpdateCourseDetails(course Course) error {
	// Convert video links to JSON
	db := dbcode.SqlRead().DB
	videoLinksJSON, err := json.Marshal(course.VideoLinks)
	if err != nil {
		return err
	}

	query := `
	UPDATE coursematerial SET
		course_name = ?,
		module_link = ?,
		video_links = ?,
		assessment_link = ?
	WHERE program_code = ?`

	_, err = db.Exec(
		query,
		course.CourseName,
		course.ModuleLink,
		string(videoLinksJSON),
		course.AssessmentLink,
		course.ProgramCode,
	)
	return err
}

func DeleteCourseDetails(courseCode string) error {
	db := dbcode.SqlRead().DB
	query := `DELETE FROM coursematerial WHERE course_code = ?`
	_, err := db.Exec(query, courseCode)
	return err
}

func GetCourseDetailsByProgram(programCode string) ([]Course, error) {
	var courses []Course
	db := dbcode.SqlRead().DB
	// Query all courses where program_code starts with the given prefix
	query := `
	SELECT 
		program_code, 
		course_name, 
		course_code,
		module_link, 
		video_links, 
		assessment_link 
	FROM coursematerial 
	WHERE program_code LIKE ? || '%'`

	rows, err := db.Query(query, programCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var course Course
		var videoLinksJSON string

		err := rows.Scan(
			&course.ProgramCode,
			&course.CourseName,
			&course.ModuleLink,
			&videoLinksJSON,
			&course.AssessmentLink,
		)
		if err != nil {
			return nil, err
		}

		// Convert JSON video links back to slice
		err = json.Unmarshal([]byte(videoLinksJSON), &course.VideoLinks)
		if err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}
