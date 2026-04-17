package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

//go:embed static/index.html
var staticFiles embed.FS

type Task struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}

var (
	tasks         []Task
	taskIDCounter int
	mu            sync.Mutex
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/tasks", handleTasks)
	mux.HandleFunc("/api/tasks/", handleTaskByID)
	mux.HandleFunc("/api/reset", handleReset)
	mux.HandleFunc("/", handleIndex)

	log.Printf("task manager listening on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		log.Fatal(err)
	}
}

func handleReset(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	tasks = []Task{}
	taskIDCounter = 0
	mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "page not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		mu.Lock()
		defer mu.Unlock()
		if tasks == nil {
			tasks = []Task{}
		}
		_ = json.NewEncoder(w).Encode(tasks)
	case http.MethodPost:
		var body struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Name) == "" {
			http.Error(w, `{"error":"invalid name"}`, http.StatusBadRequest)
			return
		}
		mu.Lock()
		taskIDCounter++
		t := Task{ID: taskIDCounter, Name: strings.TrimSpace(body.Name)}
		tasks = append(tasks, t)
		mu.Unlock()
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(t)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// expected: api/tasks/{id} or api/tasks/{id}/complete
	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	// PATCH /api/tasks/{id}/complete
	if r.Method == http.MethodPatch && len(parts) == 4 && parts[3] == "complete" {
		mu.Lock()
		defer mu.Unlock()
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Complete = !tasks[i].Complete
				_ = json.NewEncoder(w).Encode(tasks[i])
				return
			}
		}
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	// DELETE /api/tasks/{id}
	if r.Method == http.MethodDelete && len(parts) == 3 {
		mu.Lock()
		defer mu.Unlock()
		for i := range tasks {
			if tasks[i].ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
