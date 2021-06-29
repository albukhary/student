package main

import (
	"fmt"
	"log"
	"os"

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

func main() {
	// variable of type pointer to a database
	var db *sqlx.DB
	var err error

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
	insertStudent := `INSERT INTO student (name, email, age) VALUES ($1, $2, $3);`

	// Insert persons
	db.MustExec(insertStudent, "Lazizbek", "lazizbek@gmail.com", 21)
	db.MustExec(insertStudent, "Zafar aka", "zafarAka@novalab.com", 23)
	db.MustExec(insertStudent, "Izzat aka", "izzatAka@novalab.com", 23)

}
