package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"

	_ "github.com/albukhary/student/docs"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// @title CRUD student API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email lazizbekkahramonov@sgmail.com
// @license.name Novalab Tech
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8083

//struct Student represents body of
type Student struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Age   int    `db:"age"`
}

// variable of type pointer to a database
var db *sqlx.DB
var err error

func main() {

	app := fiber.New()
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: false,
	}))
	app.Listen(":8083")

	//Loading environment variables for DATABASE connection
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	//open and connect to the database at the same time
	db, err = sqlx.Connect(dialect, dbURI)
	if err != nil {
		log.Fatal(err)
	}

	// API routes
	router := mux.NewRouter()

	router.HandleFunc("/students", getStudents).Methods("GET")
	router.HandleFunc("/student/{id}", getStudent).Methods("GET")

	router.HandleFunc("/create/student", createStudent).Methods("POST")

	router.HandleFunc("/delete/student/{id}", deleteStudent).Methods("DELETE")

	router.HandleFunc("/update/student/{id}", updateStudent).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8081", router))
}

// API Controllers

// swagger: route GET /students students listStudents
// Returns a list of students
// responses :
// 200: studentsListResponse
func getStudents(w http.ResponseWriter, r *http.Request) {
	var students []Student

	err = db.Select(&students, "SELECT * FROM student")

	json.NewEncoder(w).Encode(&students)
}

// swagger: route GET /student/{id} student getStudent
// Finds and returns a particular student with the requested ID
// responses:
// 200: studentResponse

// constroller of Person
func getStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var student Student

	id, err1 := strconv.Atoi(params["id"])
	if err1 != nil {
		log.Fatal(err1)
	}
	err = db.Get(&student, "SELECT id, name, email, age FROM student WHERE id=$1", id)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(student)
}

// swagger: route POST /create/student createStudent
// creates a student of given parameters and writes into the database
//
// Consumes:
// - application/json
// Produces:
// - application/json

// Postman will send student data as JSON
// and we will put it into student struct and then into database
func createStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	json.NewDecoder(r.Body).Decode(&student)

	// insert into Person query
	insertStudent := `INSERT INTO student (id, name, email, age) VALUES ($1, $2, $3, $4);`

	// Insert the student
	db.MustExec(insertStudent, student.ID, student.Name, student.Email, student.Age)

	// print the newly added student
	json.NewEncoder(w).Encode(&student)
}

// swagger: route DELETE /delete/student/{id} delete student deleteStudent
// Finds and deletes a student with the requested ID
// Consumes:
// - application/json
// Produces:
// - application/json

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var student Student

	id, err1 := strconv.Atoi(params["id"])
	if err1 != nil {
		log.Fatal(err1)
	}
	err = db.Get(&student, "SELECT id, name, email, age FROM student WHERE id=$1", id)
	if err != nil {
		fmt.Println("There is no student with such ID")
		log.Fatal(err)
	}

	deleteQuery := `DELETE FROM student WHERE id=$1`

	//execute deletion
	db.MustExec(deleteQuery, student.ID)

	json.NewEncoder(w).Encode(&student)
}

// swagger: route PUT /student/{id} update student updateStudent
// Updates the details of the student with the requested ID according to the requested parameters
// Consumes:
// - application/json
// Produces:
// -application/json

// Update controller
func updateStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	json.NewDecoder(r.Body).Decode(&student)

	// insert into Person query
	updateStudent := `UPDATE student SET name=$1, email=$2, age=$3;`

	// Insert the student into the database
	db.MustExec(updateStudent, student.Name, student.Email, student.Age)

	// print the newly updated student details
	json.NewEncoder(w).Encode(&student)
}
