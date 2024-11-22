package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var (
	tasks  = []Task{}
	nextID = 1
	mu     sync.Mutex
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/tasks", tasksHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, nil)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTasks(w)
	case "POST":
		createTask(w, r)
	case "DELETE":
		deleteTask(w, r)
	case "PUT":
		updateTask(w, r)
	}
}

func getTasks(w http.ResponseWriter) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	title := r.FormValue("title")
	task := Task{ID: nextID, Title: title, Completed: false}
	nextID++
	tasks = append(tasks, task)
	w.WriteHeader(http.StatusCreated)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	id := r.URL.Query().Get("id")
	for i, task := range tasks {
		if strconv.Itoa(task.ID) == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	id := r.URL.Query().Get("id")
	completed := r.FormValue("completed") == "true"
	for i, task := range tasks {
		if strconv.Itoa(task.ID) == id {
			tasks[i].Completed = completed
			break
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
