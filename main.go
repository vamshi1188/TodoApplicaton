package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var taskItems = []Task{
	{Description: "watch go crash course", Completed: false},
	{Description: "watch nana's golang course", Completed: false},
	{Description: "Reward myself with a donut", Completed: false},
}

func main() {
	http.HandleFunc("/", renderHomePage)
	http.HandleFunc("/tasks/add", addTask)
	http.HandleFunc("/tasks/edit", editTask)
	http.HandleFunc("/tasks/complete", markTaskComplete)
	http.HandleFunc("/tasks/delete", deleteTask)

	// Start the server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failed to start:", err)
	}
}

func renderHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Todo List</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f0f0f0;
					margin: 0;
					padding: 20px;
				}
				h1 {
					color: #333;
					text-align: center;
				}
				.container {
					width: 60%;
					margin: 0 auto;
					background-color: #fff;
					padding: 20px;
					border-radius: 10px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
				}
				ul {
					list-style-type: none;
					padding: 0;
				}
				li {
					padding: 10px;
					border-bottom: 1px solid #ddd;
					display: flex;
					justify-content: space-between;
					align-items: center;
				}
				li:last-child {
					border-bottom: none;
				}
				.completed {
					text-decoration: line-through;
					color: #888;
				}
				a {
					text-decoration: none;
					color: #007bff;
					margin-left: 10px;
				}
				a:hover {
					color: #0056b3;
				}
				.add-task-form, .edit-form {
					display: flex;
					justify-content: space-between;
					margin-top: 20px;
				}
				.add-task-form input, .edit-form input {
					width: 80%;
					padding: 10px;
					border: 1px solid #ddd;
					border-radius: 5px;
				}
				.add-task-form button, .edit-form button {
					padding: 10px 20px;
					background-color: #007bff;
					color: #fff;
					border: none;
					border-radius: 5px;
					cursor: pointer;
				}
				.add-task-form button:hover, .edit-form button:hover {
					background-color: #0056b3;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>###### Welcome to our Todo List Application #######</h1>
				<h2>Todo List</h2>
				<ul>
					{{range $index, $task := .}}
						<li>
							<span class="{{if $task.Completed}}completed{{end}}">{{$task.Description}}</span>
							<div>
								<a href="/tasks/complete?id={{$index}}&completed={{not $task.Completed}}">
									[{{if $task.Completed}}Undo Complete{{else}}Complete{{end}}]
								</a>
								<a href="/tasks/delete?id={{$index}}">[Delete]</a>
								<form action="/tasks/edit" method="get" class="edit-form" style="display:inline;">
									<input type="hidden" name="id" value="{{$index}}" />
									<input type="text" name="task" placeholder="Edit task" />
									<button type="submit">Edit</button>
								</form>
							</div>
						</li>
					{{end}}
				</ul>
				<h2>Add a New Task</h2>
				<form action="/tasks/add" method="get" class="add-task-form">
					<input type="text" name="task" placeholder="New task" required />
					<button type="submit">Add Task</button>
				</form>
			</div>
		</body>
		</html>
	`

	t := template.New("homePage")
	t, err := t.Parse(tmpl)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	t.Execute(w, taskItems)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	newTask := r.URL.Query().Get("task")
	if strings.TrimSpace(newTask) == "" {
		http.Error(w, "Please provide a task to add", http.StatusBadRequest)
		return
	}

	taskItems = append(taskItems, Task{Description: newTask, Completed: false})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func editTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	newTask := r.URL.Query().Get("task")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 || id >= len(taskItems) {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(newTask) == "" {
		http.Error(w, "Please provide a new task description", http.StatusBadRequest)
		return
	}

	taskItems[id].Description = newTask
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func markTaskComplete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	completedStr := r.URL.Query().Get("completed")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 || id >= len(taskItems) {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	completed, err := strconv.ParseBool(completedStr)
	if err != nil {
		http.Error(w, "Invalid completed value", http.StatusBadRequest)
		return
	}

	taskItems[id].Completed = completed
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 || id >= len(taskItems) {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskItems = append(taskItems[:id], taskItems[id+1:]...)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
