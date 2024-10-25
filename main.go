package main

import (
	"fmt"
	"net/http"
	"strings"
)

var taskItems = []string{
	"watch go crash course",
	"watch nana's golang course",
	"Reward myself with a donut",
}

func main() {
	http.HandleFunc("/", greetUser)
	http.HandleFunc("/tasks", listTasks)
	http.HandleFunc("/add", addTask)

	// Start the server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failed to start:", err)
	}
}

func greetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "###### Welcome to our todo list application #######")
}

func listTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "List of my todos:")

	for index, task := range taskItems {
		fmt.Fprintf(w, "%d: %s\n", index+1, task)
	}
}

func addTask(w http.ResponseWriter, r *http.Request) {
	newTask := r.URL.Query().Get("task")
	if strings.TrimSpace(newTask) == "" {
		http.Error(w, "Please provide a task to add", http.StatusBadRequest)
		return
	}

	taskItems = append(taskItems, newTask)
	fmt.Fprintf(w, "Added task: %s\n", newTask)
}
