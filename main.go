package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	AddTaskOption    = "1"
	ListTasksOption  = "2"
	ToggleTaskOption = "3"
	DeleteTaskOption = "4"
	ExitOption       = "5"
	taskFile         = "tasks.json"
)

type Task struct {
	Name string
	Done bool
}

func (t Task) Status() string {
	if t.Done {
		return "Done"
	}
	return "Not done"
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	tasks := loadTasks()

	for {
		clearScreen()
		showMenu()
		fmt.Print("> ")
		choice := readInput(reader)
		fmt.Println()

		switch choice {
		case AddTaskOption:
			addTask(reader, &tasks)
		case ListTasksOption:
			listTasks(tasks)
		case ToggleTaskOption:
			toggleTask(reader, &tasks)
		case DeleteTaskOption:
			deleteTask(reader, &tasks)
		case ExitOption:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice")
		}

		fmt.Println("\nPress ENTER to continue...")
		reader.ReadString('\n')
	}
}

func readInput(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func showMenu() {
	fmt.Println("\nTo-Do CLI")
	fmt.Println("1) Add a task")
	fmt.Println("2) List all tasks")
	fmt.Println("3) Toggle task status")
	fmt.Println("4) Delete a task")
	fmt.Println("5) Exit")
}

func addTask(reader *bufio.Reader, tasks *[]Task) {
	fmt.Print("Enter task: ")
	taskName := readInput(reader)
	task := Task{Name: taskName, Done: false}
	*tasks = append(*tasks, task)
	saveTasks(*tasks)
	fmt.Printf("Task added: %s\n", taskName)
}

func listTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks yet")
		return
	}
	fmt.Println("Your tasks:")
	for i, t := range tasks {
		fmt.Printf("%d) %s [%s]\n", i+1, t.Name, t.Status())
	}
}

func selectTask(reader *bufio.Reader, tasks []Task, action string) (int, bool) {
	if len(tasks) == 0 {
		fmt.Printf("No tasks to %s\n", action)
		return 0, false
	}

	listTasks(tasks)
	fmt.Println()
	fmt.Printf("Enter task number to %s: ", action)
	taskNumStr := readInput(reader)

	taskIndex, err := strconv.Atoi(taskNumStr)
	if err != nil || taskIndex < 1 || taskIndex > len(tasks) {
		fmt.Println("Invalid number")
		return 0, false
	}

	return taskIndex - 1, true
}

func toggleTask(reader *bufio.Reader, tasks *[]Task) {
	taskIndex, ok := selectTask(reader, *tasks, "toggle")
	if !ok {
		return
	}

	(*tasks)[taskIndex].Done = !(*tasks)[taskIndex].Done
	saveTasks(*tasks)
	fmt.Printf("%d) %s [%s]\n", taskIndex+1, (*tasks)[taskIndex].Name, (*tasks)[taskIndex].Status())
}

func deleteTask(reader *bufio.Reader, tasks *[]Task) {
	taskIndex, ok := selectTask(reader, *tasks, "delete")
	if !ok {
		return
	}

	taskName := (*tasks)[taskIndex].Name

	*tasks = append((*tasks)[:taskIndex], (*tasks)[taskIndex+1:]...)
	saveTasks(*tasks)
	fmt.Printf("Task \"%s\" deleted\n", taskName)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func saveTasks(tasks []Task) {
	file, err := os.Create(taskFile)
	if err != nil {
		fmt.Println("Error saving tasks: ", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tasks); err != nil {
		fmt.Println("Error encoding tasks:", err)
	}
}

func loadTasks() []Task {
	file, err := os.Open(taskFile)
	if err != nil {
		return []Task{}
	}
	defer file.Close()

	var tasks []Task
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		fmt.Println("Error decoding tasks:", err)
		return []Task{}
	}
	return tasks
}
