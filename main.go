package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Model for courses - file

type Course struct {
	CourseId    uuid.UUID `json:"courseid"`
	CourseName  string    `json:"name"`
	CoursePrice int       `json:"price"`
	Author      *Author   `json:"author"`
}

type Author struct {
	FullName string `json:"name"`
	Website  string `json:"website"`
}

// fake DB

var courses []Course

// middleware / helper - file

func (course *Course) IsEmpty() bool {
	//return course.CourseId == "" && course.CourseName == ""
	return course.CourseName == ""
}

func main() {

	r := mux.NewRouter()

	//seeding
	courses = append(courses, Course{
		CourseName:  "Kotlin",
		CoursePrice: 200,
		Author: &Author{
			FullName: "Aakash",
			Website:  "aakash.me",
		},
	})

	courses = append(courses, Course{
		CourseName:  "Golang",
		CoursePrice: 300,
		Author: &Author{
			FullName: "Ali",
			Website:  "coursera.com",
		},
	})

	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getSingleCourse).Methods("GET")
	r.HandleFunc("/course", createSingleCourse).Methods("POST")
	r.HandleFunc("/course/{id}", updateSingleCourse).Methods("PUT")
	r.HandleFunc("/course/{id}", deleteSingleCourse).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}

//controllers - file

// serve home route
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to Course Management System<h1/>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all courses")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(courses)
}

func getSingleCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get single course")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	sentId, err := uuid.Parse(params["id"])

	if err != nil {
		json.NewEncoder(w).Encode("Invalid uuid")
		return
	}

	for _, course := range courses {
		if course.CourseId == sentId {
			json.NewEncoder(w).Encode(course)
			return
		}
	}

	json.NewEncoder(w).Encode("Course not found with id " + params["id"])
	return
}

func createSingleCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create single course")
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course)

	if course.IsEmpty() {
		json.NewEncoder(w).Encode("Invalid json")
		return
	}

	course.CourseId = uuid.New()
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
	return
}

func updateSingleCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update single course")
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		json.NewEncoder(w).Encode("Invalid uuid")
		return
	}

	var newCourse Course
	for index, course := range courses {
		if course.CourseId == id {
			courses = append(courses[:index], courses[index+1:]...)
			_ = json.NewDecoder(r.Body).Decode(&newCourse)
			if newCourse.IsEmpty() {
				json.NewEncoder(w).Encode("Invalid body")
				return
			} else {
				newCourse.CourseId = id
				courses = append(courses, newCourse)
				json.NewEncoder(w).Encode(newCourse)
				return
			}
		}
	}
}

func deleteSingleCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete single course")
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		json.NewEncoder(w).Encode("Invalid uuid")
		return
	}

	for index, course := range courses {
		if course.CourseId == id {
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode("Deleted successfully")
			return
		}
	}
}
