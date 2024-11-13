package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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

func SaveList(todolist *TodoList, path string) error {
	file := Must(os.Create(path))

	defer file.Close()

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(todolist)
}

func LoadList(path string) (*TodoList, error) {
	todolist := &TodoList{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		var response string
		fmt.Print("No Todolist detected!! Would you like to create one? (y/n)\n>> ")
		fmt.Scanln(&response)

		if response == "y" || response == "Y" {
			err := SaveList(todolist, path)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, nil
		}

	} else {
		// Open the todolist file
		file := Must(os.Open(path))

		defer file.Close()

		// Decode the YAML into the todolist struct
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
		Completed:   false,
	}

	todolist.Items = append(todolist.Items, newItem)

	return SaveList(todolist, path)
}

func RemoveItem(todolist *TodoList, title string, path string) error {
  search := strings.ToLower(title)
  var index int

  for i := range todolist.Items {
    if strings.ToLower(todolist.Items[i].Title) == search{
      index = i
    }
  }

  todolist.Items = append(todolist.Items[:index], todolist.Items[index+1:]...)
  
  return SaveList(todolist, path)
}

func CompleteItem(todolist *TodoList, title string, path string) error {
	completed_item := strings.ToLower(title)

	for i := range todolist.Items {
		if strings.ToLower(todolist.Items[i].Title) == completed_item {
			todolist.Items[i].Completed = true

			return SaveList(todolist, path)
		}
	}

	return fmt.Errorf("Item not found!")
}

func CheckCompleted(todolist *TodoList) {
	completedCount := 0

	for i := range todolist.Items {
		if todolist.Items[i].Completed {
			completedCount++
		}
	}

	if completedCount > 0 {
		fmt.Printf("\n\t !!! You have %d completed item(s) !!!\n\n", completedCount)
	}
}

// func RemoveCompleted(todolist *TodoList)
func PrintList(todolist *TodoList) {
	fmt.Println("\n __ __        ___         _       _    _        _  ")
	fmt.Println("|  \\  \\ _ _  |_ _| ___  _| | ___ | |  <_> ___ _| |_")
	fmt.Println("|     || | |  | | / . \\/ . |/ . \\| |_ | |<_-<  | | ")
	fmt.Println("|_|_|_|`_. |  |_| \\___/\\___|\\___/|___||_|/__/  |_| ")
	fmt.Println("       <___'                                        ")
	fmt.Println("========================================================")
	for _, item := range todolist.Items {
    var comp string
		if item.Completed {
			comp = "☑"
		} else {
			comp = "☐"
		}
		fmt.Printf("Title:          %s\n", strings.ToUpper(item.Title))
		fmt.Println("--------------------------------------------------------")
		fmt.Printf("Description:    %s\n", item.Description)
		fmt.Printf("Priority:       %d\n", item.Priority)
		fmt.Printf("Completed:      %s\n", comp)
		fmt.Println("========================================================")

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

	addPtr := flag.Bool(
    "add", 
    false, 
    "used to add a new item to the list",
  )
  remPtr := flag.String(
    "remove",
    "",
    "indicate a todolist item to be removed",
  )
	compPtr := flag.String(
		"completed",
		"",
		"indicate a todolist item that has been completed",
	)

	flag.Parse()

	if *addPtr {
		err := AddItem(todolistPath, todolist)

		if err != nil {
			fmt.Println("Error adding item!")
		} else {
			fmt.Println("Item added successfully!")
		}
	}

	if *compPtr != "" {
		err := CompleteItem(todolist, todolistPath, *compPtr)
		if err == nil {
			fmt.Printf("Successfully completed %s\n", *compPtr)
		} else {
			fmt.Println("Error completing item!", err)
		}
	}

  if *remPtr != "" {
    err := RemoveItem(todolist, *remPtr, todolistPath)
    if err == nil {
      fmt.Printf("Successfully removed %s!\n", *remPtr)
    } else {
      fmt.Println("Error completing item!", err)
    }
  }

	PrintList(todolist)
	CheckCompleted(todolist)
}
