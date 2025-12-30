package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	Id      int    `json:"Id`
	Name    string `json:"name"`
	Company string `json:"company"`
	City    string `json:"city"`
}

var users = []User{
	{Id: 1, Name: "mohit", Company: "google", City: "Indore"},
}

func main() {
	fmt.Println("ha bhai sab sahi chal raha h ")
	router := mux.NewRouter()
	router.HandleFunc("/getUsers", getUsers).Methods(http.MethodGet)
	router.HandleFunc("/createUser", createUser).Methods(http.MethodPost)
	router.HandleFunc("/getUser/{id}", getUserById).Methods(http.MethodGet)
	router.HandleFunc("/getTime", getTime).Methods(http.MethodGet)
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		http.Error(w, `{"error":"Body required"}`, http.StatusBadRequest) // 400
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest) // 400
		return
	}

	if user.Name == "" || user.Company == "" || user.City == "" {
		http.Error(w, `{"error":"All fields are required"}`, http.StatusBadRequest) // 400
		return
	}

	if user.Name == "error" {
		http.Error(w, `{"error":"Internal Server Error"}`, http.StatusInternalServerError) // 500
		return
	}

	users = append(users, user)
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(user)
}
func getUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// url se id ninkali apan ne ider
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}

	// loop lagaya h har ek user k liye
	for _, user := range users {
		if user.Id == id {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user)
			return // match milte hi exit
		}
	}

	// User not found
	http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
}

func getTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// query param "tz" se timezone le lo
	tz := r.URL.Query().Get("tz")
	if tz == "" {
		http.Error(w, `{"error":"tz query param required"}`, http.StatusBadRequest)
		return
	}

	// time.LoadLocation se timezone load karo
	// time.LoadLocation go me inbuild module h time k liye
	loc, err := time.LoadLocation(tz)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"invalid timezone %s"}`, tz), http.StatusNotFound)
		return
	}

	// current time us timezone ke hisaab se
	currentTime := time.Now().In(loc)

	// JSON response
	// response ko aapne hisab se edit kiya h
	resp := map[string]string{
		"timezone": tz,
		"time":     currentTime.Format(time.RFC3339), // ISO format
	}
	json.NewEncoder(w).Encode(resp)
}
