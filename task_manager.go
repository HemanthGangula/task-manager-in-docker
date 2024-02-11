package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	ID       int
	Name     string
	Complete bool
}

var tasks []Task
var taskIdCounter int

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Task Manager")
		fmt.Println("1. Add Task")
		fmt.Println("2. List Tasks")
		fmt.Println("3. Mark Task as Complete")
		fmt.Println("4. Exit")
		fmt.Print("Enter your choice: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("Enter task name: ")
			taskName, _ := reader.ReadString('\n')
			taskName = strings.TrimSpace(taskName)
			addTask(taskName)
		case "2":
			listTasks()
		case "3":
			fmt.Print("Enter task ID to mark as complete: ")
			taskIDInput, _ := reader.ReadString('\n')
			taskIDInput = strings.TrimSpace(taskIDInput)
			taskID, err := strconv.Atoi(taskIDInput)
			if err != nil {
				fmt.Println("Invalid task ID")
				continue
			}
			markTaskAsComplete(taskID)
		case "4":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func addTask(name string) {
	taskIdCounter++
	task := Task{
		ID:       taskIdCounter,
		Name:     name,
		Complete: false,
	}
	tasks = append(tasks, task)
	fmt.Println("Task added successfully.")
}

func listTasks() {
	if len(tasks) == 0 {
		fmt.Println("No tasks.")
		return
	}
	fmt.Println("Tasks:")
	for _, task := range tasks {
		completeStatus := "Incomplete"
		if task.Complete {
			completeStatus = "Complete"
		}
		fmt.Printf("ID: %d, Name: %s, Status: %s\n", task.ID, task.Name, completeStatus)
	}
}

func markTaskAsComplete(id int) {
	found := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Complete = true
			fmt.Println("Task marked as complete.")
			found = true
			break
		}
	}
	if !found {
		fmt.Println("Task not found.")
	}
}

