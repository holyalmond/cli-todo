package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
)

const (
	AddTaskOption = "1"
	ListTasksOption = "2"
	ToggleTaskOption = "3"
	DeleteTaskOption = "4"
	ExitOption = "5"
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

	var tasks []Task

	for {
		clearScreen()
		showMenu()
		fmt.Print("> ")
		choice := readInput(reader)

		clearScreen()
		showMenu()

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

func readInput(reader *bufio.Reader) string{
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
	fmt.Println()
}

func addTask(reader *bufio.Reader, tasks *[]Task) {
	fmt.Print("Enter task: ")
	taskName := readInput(reader)
	task := Task{Name: taskName, Done: false}
	*tasks = append(*tasks, task)
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

func toggleTask(reader *bufio.Reader, tasks *[]Task) {
	if len(*tasks) == 0 {
    	fmt.Println("No tasks to toggle")
    	return
	}	
	listTasks(*tasks)
	fmt.Println()
	fmt.Print("Enter task number: ")
	taskNumStr:= readInput(reader)

	taskIndex, err := strconv.Atoi(taskNumStr)
	if err != nil {
		fmt.Println("Invalid number")
		return
	}

	if taskIndex < 1 || taskIndex > len(*tasks) {
    	fmt.Println("No task with that number")
    	return
	}


	(*tasks)[taskIndex-1].Done = !(*tasks)[taskIndex-1].Done
	fmt.Printf("%d) %s [%s]\n", taskIndex, (*tasks)[taskIndex-1].Name, (*tasks)[taskIndex-1].Status())
}

func deleteTask(reader *bufio.Reader, tasks *[]Task) {
	if len(*tasks) == 0 {
    	fmt.Println("No tasks to delete")
    	return
	}	
	listTasks(*tasks)
	fmt.Println()
	fmt.Print("Enter task number: ")
	taskNumStr:= readInput(reader)

	taskIndex, err := strconv.Atoi(taskNumStr)
	if err != nil {
		fmt.Println("Invalid number")
		return
	}

	if taskIndex < 1 || taskIndex > len(*tasks) {
    	fmt.Println("No task with that number")
    	return
	}

	taskName := (*tasks)[taskIndex-1].Name

	*tasks = append((*tasks)[:taskIndex-1], (*tasks)[taskIndex:]...)
	fmt.Printf("Task \"%s\" deleted\n", taskName)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}