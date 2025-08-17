package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"github.com/fatih/color"
)

const (
	AddTaskOption    = "1"
	ListTasksOption  = "2"
	ToggleTaskOption = "3"
	DeleteTaskOption = "4"
	ExitOption       = "5"
	taskFile         = "tasks.json"
)

var (
	green = color.New(color.FgGreen).Add(color.Bold).SprintFunc() 
	red = color.New(color.FgRed).Add(color.Bold).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan = color.New(color.FgCyan).SprintFunc()
	bold = color.New(color.Bold).SprintFunc()
)

type Task struct {
	Name string
	Done bool
}

func (t Task) Status() string {
	if t.Done {
		return green("Done")
	}
	return yellow("Not done")
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
			fmt.Println(cyan("Goodbye!"))
			return
		default:
			fmt.Println(red("Invalid choice"))
		}

		fmt.Println(bold("\nPress ENTER to continue..."))
		reader.ReadString('\n')
	}
}

func readInput(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func showMenu() {
	fmt.Println(bold("\nTo-Do CLI"))
	fmt.Println(cyan("1)") + " Add a task")
	fmt.Println(cyan("2)") + " List all tasks")
	fmt.Println(cyan("3)") + " Toggle task status")
	fmt.Println(cyan("4)") + " Delete a task")
	fmt.Println(cyan("5)") + " Exit")
}

func addTask(reader *bufio.Reader, tasks *[]Task) {
	fmt.Print("Enter task: ")
	taskName := readInput(reader)
	task := Task{Name: taskName, Done: false}
	*tasks = append(*tasks, task)
	saveTasks(*tasks)
	fmt.Printf("%s %s\n", green("Task added:"), bold(taskName))
}

func listTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println(red("No tasks yet"))
		return
	}
	fmt.Println(cyan("Your tasks:"))
	for i, t := range tasks {
		fmt.Printf("%s %s [%s]\n", 
			yellow(fmt.Sprintf("%d)", i+1)), 
			bold(t.Name), 
			t.Status())
	}
}

func selectTask(reader *bufio.Reader, tasks []Task, action string) (int, bool) {
	if len(tasks) == 0 {
		fmt.Printf("%s\n", red(fmt.Sprintf("No tasks to %s", action)))
		return 0, false
	}

	listTasks(tasks)
	fmt.Println()
	fmt.Printf("Enter task number to %s: ", action)
	taskNumStr := readInput(reader)

	taskIndex, err := strconv.Atoi(taskNumStr)
	if err != nil || taskIndex < 1 || taskIndex > len(tasks) {
		fmt.Println(red("Invalid number"))
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
		fmt.Printf("%s %s [%s]\n",
		yellow(fmt.Sprintf("%d)", taskIndex+1)),
		bold((*tasks)[taskIndex].Name),
		(*tasks)[taskIndex].Status(),
	)	
}

func deleteTask(reader *bufio.Reader, tasks *[]Task) {
	taskIndex, ok := selectTask(reader, *tasks, "delete")
	if !ok {
		return
	}

	taskName := (*tasks)[taskIndex].Name

	*tasks = append((*tasks)[:taskIndex], (*tasks)[taskIndex+1:]...)
	saveTasks(*tasks)
	fmt.Printf("%s \"%s\" %s\n", red("Task"), bold(taskName), red("deleted"))
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
