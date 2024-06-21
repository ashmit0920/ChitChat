package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	templates = template.Must(template.ParseFiles(
		"templates/navbar.html",
		"templates/login.html",
		"templates/signup.html",
		"templates/home.html",
		"templates/chatroom.html",
	))
	users      = make(map[string]string)
	usersMutex sync.Mutex
	chatrooms  = make(map[string][]string) // maps chatroom code to messages
	roomsMutex sync.Mutex
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/createroom", createRoomHandler)
	http.HandleFunc("/joinroom", joinRoomHandler)
	http.HandleFunc("/chatroom", chatRoomHandler)
	http.HandleFunc("/sendmessage", sendMessageHandler)

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
		data := struct {
			Error string
		}{
			Error: "Invalid credentials!",
		}
		templates.ExecuteTemplate(w, "login.html", data)
		return
	}

	// Set a cookie with the username
	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: username,
		Path:  "/",
	})

	http.Redirect(w, r, "/home", http.StatusSeeOther)
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	username := cookie.Value

	data := struct {
		Username string
	}{
		Username: username,
	}
	templates.ExecuteTemplate(w, "home.html", data)
}

func generateRoomCode() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	username := cookie.Value

	roomCode := generateRoomCode()

	roomsMutex.Lock()
	chatrooms[roomCode] = []string{}
	roomsMutex.Unlock()

	http.Redirect(w, r, fmt.Sprintf("/chatroom?room=%s&username=%s", roomCode, username), http.StatusSeeOther)
}

func joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		roomCode := r.FormValue("roomCode")
		roomsMutex.Lock()
		_, exists := chatrooms[roomCode]
		roomsMutex.Unlock()

		if exists {
			http.Redirect(w, r, fmt.Sprintf("/chatroom?room=%s", roomCode), http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		}
		return
	}

	templates.ExecuteTemplate(w, "home.html", nil)
}

func chatRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomCode := r.URL.Query().Get("room")
	if roomCode == "" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	username := cookie.Value

	roomsMutex.Lock()
	messages := chatrooms[roomCode]
	roomsMutex.Unlock()

	data := struct {
		RoomCode string
		Username string
		Messages []string
	}{
		RoomCode: roomCode,
		Username: username,
		Messages: messages,
	}

	templates.ExecuteTemplate(w, "chatroom.html", data)
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	roomCode := r.FormValue("roomCode")
	username := r.FormValue("username")
	message := r.FormValue("message")

	roomsMutex.Lock()
	messages, ok := chatrooms[roomCode]
	if !ok {
		roomsMutex.Unlock()
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	fullMessage := fmt.Sprintf("%s: %s", username, message)
	chatrooms[roomCode] = append(messages, fullMessage)
	roomsMutex.Unlock()

	// Send only the chat messages as HTML
	for _, msg := range chatrooms[roomCode] {
		fmt.Fprintf(w, "<div class='chat-message'>%s</div>", msg)
	}
}
