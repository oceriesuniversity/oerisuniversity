package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/routes"



	_ "github.com/mattn/go-sqlite3"
)

func Enviroment(w http.ResponseWriter, r *http.Request) {

	result := os.Getenv("Working")
	fmt.Fprintln(w, result)
}

func main() {


	//Email Configaration

	//LOAD DATA TABLES
	dbcode.LoadDB()


	routes.LoadAssesments()
	routes.LoadCource()
	routes.LoadExamTable()
	routes.LoadAdminUsers()
	routes.LoadAssesmentTable()
	routes.CreateVisitorTable()
	routes.CreateEmailSenderTem()

	//NEW DATA BASES
	routes.CreateProgramDB()
	routes.CreateCourseDB()

	fmt.Println("::SERVER STARTED::PORT::8080")

	router := http.NewServeMux()

	//ENVIROMENTAL VARIABLES

	router.HandleFunc("/env", Enviroment)

	//MAIN SCREEN
	router.HandleFunc("/", routes.HomePage)
	router.HandleFunc("/aboutus", routes.AboutUs)
	router.HandleFunc("/programs", routes.Programs)
	router.HandleFunc("/enroll", routes.Enrollment)
	router.HandleFunc("/confirmenrrol", routes.ConfirmEnrollment)

	router.HandleFunc("/programcards", routes.Programcards)
	router.HandleFunc("/news", routes.NewsPage)
	router.HandleFunc("/readnewsstory", routes.ReadNewsRoute)
	router.HandleFunc("/delete", routes.DeleteMessageRouter)
	router.HandleFunc("/getstudentdata", routes.ChangeStudentPassword)
	router.HandleFunc("/updtestudentpassword", routes.ChangePassword)
	router.HandleFunc("/closeupdatedata", routes.CloseUpdateData)

	//ADMIN DASHBOARD
	//NEW ROUTES START
	router.HandleFunc("/adminlogin", routes.AdminLogin)
	router.HandleFunc("/admindashboard", routes.AdminDashboard)
	router.HandleFunc("/createNewProgram", routes.CreateNewProgramR)
	router.HandleFunc("/studentdata", routes.StudentData)
	router.HandleFunc("/programdetails", routes.ProgramDetails)
	router.HandleFunc("/updatecourcedata", routes.UpdateCourceData)
	router.HandleFunc("/createnewcourse", routes.AddCourceData)

	//NEW ROUTES END

	router.HandleFunc("/studentcenter", routes.StudentCenter)
	router.HandleFunc("/sendbulleting", routes.SendBulleting)

	router.HandleFunc("/adminNews", routes.AdminNews)
	router.HandleFunc("/getstudent", routes.GetAllStudentMsg)
	router.HandleFunc("/adminMessages", routes.AdminMessagesPage)
	router.HandleFunc("/studentedataadmin", routes.StudentProfileData)
	router.HandleFunc("/readmessageadmin", routes.ReadMessageAdmin)
	router.HandleFunc("/courceupdateadmin", routes.ApproveCourceUpdate)
	router.HandleFunc("/approve", routes.ApproveProgram)
	router.HandleFunc("/create_cource_data", routes.CreateCourseData)
	router.HandleFunc("/close_cource_div", routes.CloseCreateCourseData)
	router.HandleFunc("/getstudymaterial", routes.GetStudyMaterial)
	router.HandleFunc("/get_admin_users", routes.AdminUsers)
	router.HandleFunc("/createadmin", routes.CreatAdminUser)
	router.HandleFunc("/updateadminuser", routes.UpdateAdminUsers)
	router.HandleFunc("/getupdate", routes.GetUpateAdmin)
	router.HandleFunc("/loadform", routes.LoadAdminForm)
	router.HandleFunc("/saveucdaata", routes.UpdateProgramDetails)

	///
	router.HandleFunc("/updateprogramdetails", routes.UpdateAllProgramDetailsR)
	router.HandleFunc("/updateprogramdata", routes.UpdateAllProgramDetails)

	///
	router.HandleFunc("/gradeca", routes.GradeCA)
	router.HandleFunc("/delete_cource_data", routes.DeleteCource)
	router.HandleFunc("/get_matrics", routes.Matrics)
	router.HandleFunc("/deleteschool", routes.DeleteEmail)
	router.HandleFunc("/downloadfile", routes.AdminDownLoadAsignment)

	router.HandleFunc("/deleteadmin", routes.DeleteAdmin)
	router.HandleFunc("/createnews", routes.Create_News)
	router.HandleFunc("/curiculum", routes.Curiculum)
	router.HandleFunc("/databackup", dbcode.BackUpData)
	router.HandleFunc("/downloadassignment", routes.DownloadAssesments)
	router.HandleFunc("/deletestuden", routes.DeleteStudent)
	router.HandleFunc("/createadminemail", routes.CreateEmailData)
	router.HandleFunc("/faq", routes.FAQ)
	router.HandleFunc("/offerings", routes.Offerings)
	router.HandleFunc("/resetpassword", routes.ForgotPassword)
	router.HandleFunc("/confirmreset", routes.ConfirmStudentId)

	//EXAM
	router.HandleFunc("/createexampage", routes.CreatePage)///currntly using
	router.HandleFunc("/addexam", routes.AddExam)
	router.HandleFunc("/examdetails", routes.AddExamDetails)

	router.HandleFunc("/takeexam", routes.TakeExam)
	router.HandleFunc("/submitexam", routes.SubmitExam)
	router.HandleFunc("/grade_exam", routes.GradeExam)
	router.HandleFunc("/gradeexamination", routes.SaveGrades)
	router.HandleFunc("/examddfdea", routes.GetParticularExam)
	router.HandleFunc("/courcecompleted", routes.CourceCompleted)
	router.HandleFunc("/updatequestion", routes.UpdateQuestion)
	router.HandleFunc("/deletequestion", routes.DeleteQuestion)
	router.HandleFunc("/saveUpdateQuestion", routes.SaveQuestionUpdates)

	//STUDENT PORTAL ROUTES
	router.HandleFunc("/confirmlogin", routes.ConfirmStudentLogin)
	router.HandleFunc("/studentportal", routes.StudentPortal)
	router.HandleFunc("/approvecource", routes.ApproveCource)
	router.HandleFunc("/messages/{id}", routes.Messages)
	router.HandleFunc("/completed", routes.ProgramCompleted)
	router.HandleFunc("/studentsettings", routes.StudentSettings)
	router.HandleFunc("/studentlogoout/{id}", routes.StudentLogOut)
	router.HandleFunc("/proceed", routes.StudentProcced)
	router.HandleFunc("/contactinstitution/{id}", routes.ContactInstitution)
	router.HandleFunc("/watchvideo", routes.WatcVideo)
	router.HandleFunc("/sendmessage", routes.SendMsg)
	router.HandleFunc("/studentprofileportal/{id}", routes.StudentProfilePortal)
	router.HandleFunc("/close_assesment_div", routes.CloseAssesmentDiv)
	router.HandleFunc("/applytpproceed", routes.ApplyProceed)
	// router.HandleFunc("/close_admin_div", routes.CloseAdmintDiv)
	router.HandleFunc("/error", routes.ErrorPage)
	router.HandleFunc("/assesmentsubmit", routes.UploadAssesment)
	router.HandleFunc("/forgotpassword", routes.PasswordResetPage)
	router.HandleFunc("/passwordreset", routes.ResetPassword)

	router.HandleFunc("/login", routes.LoginPage)
	router.HandleFunc("/handinassesment", routes.HandInAssesment)
	router.HandleFunc("/grade_assesment", routes.GradeAssesment)
	router.HandleFunc("/deleteassesmentresults", routes.DeleteAssesmentAdmin)
	router.HandleFunc("/deletenews", routes.DeleteNewsRoute)
	router.HandleFunc("/viewadminnews", routes.ViewAdminNews)



	//LOAD ASSETS

	fs := http.FileServer(http.Dir("assets"))
	router.Handle("/assets/", http.StripPrefix("/assets/", fs))

	ns := http.FileServer(http.Dir("news"))
	router.Handle("/news/", http.StripPrefix("/news/", ns))

	//RUN SERVER
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), router))

}
