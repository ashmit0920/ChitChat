package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sync"
)

var (
	templates = template.Must(template.ParseFiles(
		"templates/navbar.html",
		"templates/login.html",
		"templates/signup.html",
	))
	users      = make(map[string]string)
	usersMutex sync.Mutex
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/register", registerHandler)

	loadUsers()

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

func loadUsers() {
	file, err := os.Open("users.json")
	if err != nil {
		fmt.Println("Error opening users file:", err)
		return
	}
	defer file.Close()

	usersMutex.Lock()
	defer usersMutex.Unlock()
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		fmt.Println("Error decoding users file:", err)
	}
}

func saveUsers() {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	file, err := os.Create("users.json")
	if err != nil {
		fmt.Println("Error creating users file:", err)
		return
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(users)
	if err != nil {
		fmt.Println("Error encoding users file:", err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "signup.html", nil)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	usersMutex.Lock()
	storedPassword, ok := users[username]
	usersMutex.Unlock()

	if !ok || storedPassword != password {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	usersMutex.Lock()
	users[username] = password
	usersMutex.Unlock()

	saveUsers()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
