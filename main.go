package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//struct Person for person table in database
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

	// insert into Person query
	//	insertStudent := `INSERT INTO student (name, email, age) VALUES ($1, $2, $3);`

	// Insert persons
	//	db.MustExec(insertStudent, "Lazizbek", "lazizbek@gmail.com", 21)
	//	db.MustExec(insertStudent, "Zafar aka", "zafarAka@novalab.com", 23)
	//	db.MustExec(insertStudent, "Izzat aka", "izzatAka@novalab.com", 23)

	// API routes
	router := mux.NewRouter()

	router.HandleFunc("/students", getStudents).Methods("GET")
	router.HandleFunc("/student/{id}", getStudent).Methods("GET")

	router.HandleFunc("/create/student", createStudent).Methods("POST")

	//	router.HandleFunc("/delete/student/{id}", deleteStudent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

// API Controllers

// controller of Persons
func getStudents(w http.ResponseWriter, r *http.Request) {
	var students []Student

	err = db.Select(&students, "SELECT * FROM student")

	json.NewEncoder(w).Encode(&students)
}

// constroller of Person
func getStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var student Student

	// find the first match from database
	//row := db.QueryRow("SELECT  FROM student WHERE id=$1", params["ID"])
	//	err = row.Scan(&student.ID, &student.Name, &student.Email, &student.Age)

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
