package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "modernc.org/sqlite"
)

var db *sql.DB

type Task struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

func main() {
	// Veritabanını başlat
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	// API rotalarını /api öneki ile gruplayalım
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/tasks", getTasks).Methods("GET")
	api.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	api.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	api.HandleFunc("/tasks", createTask).Methods("POST")
	// Diğer PUT, DELETE rotaları...

	// CORS ayarları: Geliştirme sırasında frontend'den gelen isteklere izin ver
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // SvelteKit dev sunucusu
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}).Handler(r)

	log.Println("Backend sunucusu http://localhost:8080 adresinde başlatılıyor...")
	// Sunucuyu CORS middleware ile başlat
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,        
        "name" TEXT,
        "completed" BOOLEAN
    );`
	statement, _ := db.Prepare(createTableSQL)
	statement.Exec()
}

// getTasks ve createTask fonksiyonları önceki cevaptaki ile aynı.
// ... Buraya diğer handler fonksiyonlarını (getTasks, createTask, vb.) ekleyin ...
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT id, name, completed FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}
	json.NewEncoder(w).Encode(tasks)
	// write get tasks response
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Retrieved %d tasks", len(tasks))
	if len(tasks) == 0 {
		http.Error(w, "No tasks found", http.StatusNotFound)
		return
	}
	log.Printf("Tasks: %+v", tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	statement, _ := db.Prepare("INSERT INTO tasks (name, completed) VALUES (?, ?)")
	result, err := statement.Exec(task.Name, task.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	task.ID = int(id)
	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	statement, _ := db.Prepare("UPDATE tasks SET name = ?, completed = ? WHERE id = ?")
	_, err := statement.Exec(task.Name, task.Completed, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	task.ID = intID
	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	statement, _ := db.Prepare("DELETE FROM tasks WHERE id = ?")
	_, err := statement.Exec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content
	log.Printf("Deleted task with ID %s", id)
}
