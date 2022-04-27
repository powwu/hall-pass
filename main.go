package main

import (
	"net/http"
	"log"
	"fmt"
	"database/sql"
	"time"

	_ "github.com/genjidb/genji/driver"

	"github.com/labstack/echo/v4"
)

func main() {
	db, err := sql.Open("genji", "./hall.db")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer db.Close()

	router := echo.New()
	router.Debug = true

	router.Static("/", "./web")
	router.POST("/out", func(c echo.Context) error {
		classroom := c.FormValue("classroom")
		studentId := c.FormValue("studentid")
		destination := c.FormValue("destination")
		t := time.Now()
		currentTime := t.Format("3:04PM")

		// studentExists, err := db.Query(`SELECT * FROM Students WHERE studentId = (?)`, studentId)

		_, err := db.Exec(`INSERT INTO Students (out, studentId, classroom, destination, timeOut, timeIn) VALUES (?, ?, ?, ?, ?, ?)`, true, studentId, classroom, destination, currentTime, "Student not yet in...")
		if err != nil {
			log.Fatalf("%+v", err)
		}
		return c.Redirect(http.StatusMovedPermanently, "/")
	})

	router.POST("/in", func(c echo.Context) error {
		studentId := c.FormValue("studentId")
		t := time.Now()
		currentTime := t.Format("3:04PM")

		_, err := db.Exec(`UPDATE Students set out = (?), timeIn = (?) WHERE studentId = (?) AND out = (?)`, false, currentTime, studentId, true)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		return c.Redirect(http.StatusMovedPermanently, "/admin")
	})

	router.GET("/admin", func(c echo.Context) error {
		var outStudents string
		students, err := db.Query("SELECT studentId, classroom, destination, timeOut FROM Students WHERE out = true")
		if err != nil {
			log.Fatalf("%+v", err)
		}
		for students.Next() {
			var studentId, classroom, destination, timeOut string
			students.Scan(&studentId, &classroom, &destination, &timeOut)
			stext := fmt.Sprintf("STUDENT: %s<br>CLASSROOM: %s<br>DESTINATION: %s<br>TIME OUT: %s<br><form action=\"/in\" method=\"post\" enctype\"multipart/form-data\"><input type=\"hidden\" id=\"studentId\" name=\"studentId\" value=\"%s\"><input type=\"submit\" value=\"Student In\"></form><br><br>", studentId, classroom, destination, timeOut, studentId)
			outStudents += stext
		}

		return c.HTML(http.StatusOK, "<h1>Current Out Students</h1>" + outStudents)
	})

	router.Logger.Fatal(router.Start(":8080"))
}
