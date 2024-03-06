package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var tmpl *template.Template

var db *sql.DB

func getMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/studentinfo?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func init() {
	tmpl = template.Must(template.ParseFiles("crudForm.html"))
}

type studentinfo struct {
	Sid    string
	Name   string
	Course string
}

func crudHandler(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	student := studentinfo{
		Sid:    r.FormValue("sid"),
		Name:   r.FormValue("name"),
		Course: r.FormValue("course"),
	}
	if r.FormValue("submit") == "Insert" {
		sid, _ := strconv.Atoi(student.Sid)
		_, err := db.Exec("insert into student (sid, name, course) values(?, ?, ?)", sid, student.Name, student.Course)
		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: err.Error()})
		} else {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: "Record Inserted Successfully"})
		}

	} else if r.FormValue("submit") == "Read" {
		data := []string{}
		rows, err := db.Query("select * from student")
		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: err.Error()})
		} else {
			stud := studentinfo{}
			data = append(data, "<table border=1>")
			data = append(data, "<tr><th>Student Id</th><th>Student Name</th><th>Student Course</th></tr>")
			for rows.Next() {
				rows.Scan(&stud.Sid, &stud.Name, &stud.Course)
				data = append(data, fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", stud.Sid, stud.Name, stud.Course))
			}
			data = append(data, "</table>")
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: strings.Trim(fmt.Sprint(data), "[]")})
		}
	} else if r.FormValue("submit") == "Update" {
		sid, _ := strconv.Atoi(student.Sid)
		result, err := db.Exec("update student set name=?, course=? where sid=?", student.Name, student.Course, sid)
		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: err.Error()})
		} else {
			_, err := result.RowsAffected()
			if err != nil {
				tmpl.Execute(w, struct {
					Success bool
					Message string
				}{Success: true, Message: "Record not Updated"})
			} else {
				tmpl.Execute(w, struct {
					Success bool
					Message string
				}{Success: true, Message: "Record Updated"})
			}
		}
	} else if r.FormValue("submit") == "Delete" {
		sid, _ := strconv.Atoi(student.Sid)
		result, err := db.Exec("delete from student where sid=?", sid)
		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: err.Error()})
		} else {
			_, err := result.RowsAffected()
			if err != nil {
				tmpl.Execute(w, struct {
					Success bool
					Message string
				}{Success: true, Message: "Record not Deleted"})
			} else {
				tmpl.Execute(w, struct {
					Success bool
					Message string
				}{Success: true, Message: "Record Deleted"})
			}

		}
	}

	fmt.Println(student)

}

func main() {
	http.HandleFunc("/", crudHandler)
	http.ListenAndServe(":8080", nil)
}
