package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Course struct {
	CourseId    int     `json:"course_id"`
	CourseName  string  `json:"course_name"`
	CoursePrice float64 `json:"course_price"`
	Author      *Author `json:"author"`
}

type Author struct {
	Fullname string `json:"fullname"`
	Website  string `json:"website"`
}

var courses []Course

func seedCourses() {
	courses = append(courses, Course{
		CourseId:    1,
		CourseName:  "Course1",
		CoursePrice: 1000,
		Author: &Author{
			Fullname: "vaibhav",
			Website:  "https://www.linkedin.com/in/neogvaibhav/",
		},
	})
}

func main() {
	fmt.Println("Starting the server...")

	r := mux.NewRouter()

	seedCourses()

	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", createCourse).Methods("POST")
	r.HandleFunc("/courses/{id}", updateCourse).Methods("PUT")

	log.Fatal(http.ListenAndServe(":4000", r))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creator of this route is awesome and he welocme's you!")
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	var newCourse Course
	err := json.NewDecoder(r.Body).Decode(&newCourse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, course := range courses {
		if course.CourseId == newCourse.CourseId {
			http.Error(w, "Course with the same ID already exist", http.StatusBadRequest)
			return
		}
	}

	newCourse.CourseId = len(courses) + 1

	courses = append(courses, newCourse)
	json.NewEncoder(w).Encode(newCourse)
}

func updateCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	courseID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var updatedCourse Course
	err = json.NewDecoder(r.Body).Decode(&updatedCourse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i, course := range courses {
		if course.CourseId == courseID {
			courses[i].CourseName = updatedCourse.CourseName
			courses[i].CoursePrice = updatedCourse.CoursePrice
			courses[i].Author = updatedCourse.Author
			json.NewEncoder(w).Encode(courses[i])
			return
		}
	}
	http.Error(w, "Course not found", http.StatusNotFound)
}
