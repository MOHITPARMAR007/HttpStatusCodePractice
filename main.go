package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	City    string `json:"city"`
}

var users = []User{
	{Name: "mohit", Company: "google", City: "Indore"},
}

func main() {
	fmt.Println("ha bhai sab sahi chal raha h ")
	router := mux.NewRouter()
	router.HandleFunc("/getUsers", getUsers).Methods(http.MethodGet)
	router.HandleFunc("/createUser", createUser).Methods(http.MethodPost)
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
