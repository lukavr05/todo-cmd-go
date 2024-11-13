package main

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type TodoList struct {
	Name  string `yaml:"name"`
	Items []Item `yaml:"items"`
}

type Item struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Priority    int    `yaml:"priority"`
	Completed   bool   `yaml:"completed"`
}

func SaveList(path string, todolist *TodoList) error {
	file := Must(os.Create(path))

	defer file.Close()

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(todolist)
}

func LoadList(path string) (*TodoList, error) {
	todolist := &TodoList{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		var response string
		fmt.Print("No Todolist detected!! Would you like to create one? (y/n)")
		fmt.Scanln(&response)

		if response == "y" || response == "Y" {
			err := SaveList(path, todolist)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, nil
		}

	} else {
		// Open the config file
		file := Must(os.Open(path))

		defer file.Close()

		// Decode the JSON into the config struct
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(todolist); err != nil {
			return nil, err
		}
	}

	return todolist, nil
}

func AddItem(path string, todolist *TodoList) error {
	var title, description string
	var priority int

	fmt.Print("Enter the title for the item:  ")
	fmt.Scanln(&title)

	fmt.Print("Enter the description for the item:  ")
	fmt.Scanln(&description)

	fmt.Print("Enter the priority of the item (1-5):  ")
	fmt.Scanln(&priority)

	newItem := Item{
		Title:       title,
		Description: description,
		Priority:    priority,
	}

	todolist.Items = append(todolist.Items, newItem)

	return SaveList(path, todolist)
}

func PrintList(todolist *TodoList) {
	for _, item := range todolist.Items {
		fmt.Printf("Title:          %s\n", item.Title)
		fmt.Printf("Description:    %s\n", item.Description)
		fmt.Printf("Priority:       %d\n", item.Priority)
	}
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func main() {
	todolistPath := "todolist.yaml"
	todolist := Must(LoadList(todolistPath))

	addPtr := flag.Bool("add", false, "used to add a new item to the list")

	flag.Parse()

	if *addPtr {
    err := AddItem(todolistPath, todolist)

		if err != nil {
			fmt.Println("Error adding item!")
		} else {
			fmt.Println("Item added successfully!")
		}
	}
	PrintList(todolist)
}
