package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
  "io"

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
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()


}

func LoadList(path string) (*TodoList, error) {
	todolist := &TodoList{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := SaveList(path, todolist)
		if err != nil {
			return nil, err
		}
	} else {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(todolist)
		if err != nil {
			return nil, err
		}
	}
	return todolist, nil
}

func addItem(todolist *TodoList) {
	in := `{"items":[]}`
	json.Unmarshal([]byte(in), todolist)

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

	j, _ := json.Marshal(todolist)
	fmt.Println(string(j))
}

func PrintList(todolist *TodoList) {
	for _, item := range todolist.Items {
		fmt.Printf("Title:          %s\n", item.Title)
		fmt.Printf("Description:    %s\n", item.Description)
		fmt.Printf("Priority:       %d\n", item.Priority)
	}
}

func main() {
	todolistPath := "todolist.json"
	todolist, err := LoadList(todolistPath)
	if err != nil {
		print(" !!! Error loading todo list: ", err)
	}

	addPtr := flag.Bool("add", false, "used to add a new item to the list")

	flag.Parse()

	if *addPtr {
		addItem(todolist)
	}

	SaveList(todolistPath, todolist)
	PrintList(todolist)
}
