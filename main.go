package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
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

		if strings.ToLower(response) == "y" {
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
  scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the title for the item:          ")
	scanner.Scan()
  title := scanner.Text()

	fmt.Print("Enter the description for the item:    ")
	scanner.Scan()
  description := scanner.Text()

	fmt.Print("Enter the priority of the item (1-5):  ")
  
  var priority int
  _, err := fmt.Scanln(&priority)
  for err != nil {
    fmt.Println("Invalid input! Please try again!")
    _, err = fmt.Scanln(&priority)  
  }

	if priority < 1 {
		priority = 1
	}
	if priority > 5 {
		priority = 5
	}

	newItem := Item{
		Title:       title,
		Description: description,
		Priority:    priority,
		Completed:   false,
	}

	todolist.Items = append(todolist.Items, newItem)

  PrintList(todolist.Items)
	return SaveList(todolist, path)
}

func RemoveItem(todolist *TodoList, title string, path string) error {
	search := strings.ToLower(title)
	var index int

	for i := range todolist.Items {
		if strings.ToLower(todolist.Items[i].Title) == search {
			index = i
		}
	}

	todolist.Items = append(todolist.Items[:index], todolist.Items[index+1:]...)

	return SaveList(todolist, path)
}

func RemoveAll(todolist *TodoList, path string) error {
	var response string
	fmt.Print("Are you sure you want to remove all items from the Todo List? (y/n)\n>> ")
	fmt.Scanln(&response)

	if strings.ToLower(response) == "y" {
		todolist = &TodoList{}
	}

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

func CompleteAll(todolist *TodoList, path string) error {
	for i := range todolist.Items {
		todolist.Items[i].Completed = true
	}

	return SaveList(todolist, path)
}

func CheckCompleted(todolist *TodoList) int {
	completedCount := 0

	for i := range todolist.Items {
		if todolist.Items[i].Completed {
			completedCount++
		}
	}
	return completedCount
}

func RemoveCompleted(todolist *TodoList, path string) error {
	var response string
	fmt.Print("\n\nWould you like to delete all completed items? (y/n)\n>> ")
	fmt.Scanln(&response)

	if strings.ToLower(response) == "y" {
		// Filter out completed items
		var updatedItems []Item
		for _, item := range todolist.Items {
			if !item.Completed {
				updatedItems = append(updatedItems, item)
			}
		}
		todolist.Items = updatedItems

		// Save the updated list
		if err := SaveList(todolist, path); err != nil {
			return fmt.Errorf("Failed to save updated list: %v", err)
		}

		fmt.Println("All completed items have been removed.")
		PrintList(todolist.Items)
	}
	return nil
}

func PrintList(items []Item, headers ...string) {
	if len(headers) > 0 {
		header := headers[0]
		fmt.Printf("\nSorted by %s!\n", header)
	}
	fmt.Println("========================================================")
	if len(items) == 0 {
		fmt.Println("Todo list is empty :(")
		fmt.Println("========================================================")
	} else {
		for _, item := range items {
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
}

func PrintSortedList(todolist *TodoList, sortBy string) {
	sortedItems := make([]Item, len(todolist.Items))
	copy(sortedItems, todolist.Items)

	switch strings.ToLower(sortBy) {
	case "title":
		sort.SliceStable(sortedItems, func(i, j int) bool {
			return strings.ToLower(sortedItems[i].Title) < strings.ToLower(sortedItems[j].Title)
		})
	case "priority":
		sort.SliceStable(sortedItems, func(i, j int) bool {
			return sortedItems[i].Priority < sortedItems[j].Priority
		})
	case "completed":
		sort.SliceStable(sortedItems, func(i, j int) bool {
			return !sortedItems[i].Completed && sortedItems[j].Completed
		})
	default:
		fmt.Println("Invalid sort field. Valid options are: title, priority, completed.")
		return
	}

	PrintList(sortedItems, strings.ToLower(sortBy))
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
		"r",
		"",
		"indicate a todolist item to be removed",
	)
	remAllPtr := flag.Bool("ra",
    false, 
    "remove all items from the Todo List",
  )

	compPtr := flag.String(
		"c",
		"",
		"indicate a todolist item that has been completed",
	)
	compAllPtr := flag.Bool(
		"ca",
		false,
		"indicate that all items have been completed",
	)
	sortPtr := flag.String(
		"s",
		"",
		"sort the todo list by a specific field (title, priority, completed)",
	)

	flag.Parse()

	fmt.Println("\n __ __        ___         _       _    _        _  ")
	fmt.Println("|  \\  \\ _ _  |_ _| ___  _| | ___ | |  <_> ___ _| |_")
	fmt.Println("|     || | |  | | / . \\/ . |/ . \\| |_ | |<_-<  | | ")
	fmt.Println("|_|_|_|`_. |  |_| \\___/\\___|\\___/|___||_|/__/  |_| ")
	fmt.Println("       <___'\n")

	if *addPtr {
		err := AddItem(todolistPath, todolist)

		if err != nil {
			fmt.Println("Error adding item!")
		} else {
			fmt.Println("Item added successfully!")
		}
	}

	if *sortPtr != "" {
		PrintSortedList(todolist, *sortPtr)
	}

	if *compPtr != "" {
		err := CompleteItem(todolist, *compPtr, todolistPath)
		if err == nil {
			fmt.Printf("Successfully completed %s\n", *compPtr)
			PrintList(todolist.Items)
		} else {
			fmt.Println("Error completing item!", err)
		}
	}

	if *compAllPtr {
		err := CompleteAll(todolist, todolistPath)
		if err == nil {
			fmt.Println("Successfully completed all items!")
			PrintList(todolist.Items)
		} else {
			fmt.Println("Error!", err)
		}
	}

	if *remPtr != "" {
		err := RemoveItem(todolist, *remPtr, todolistPath)
		if err == nil {
			fmt.Printf("Successfully removed %s!\n", *remPtr)
			PrintList(todolist.Items)
		} else {
			fmt.Println("Error completing item!", err)
		}
	}

  if *remAllPtr {
    err := RemoveAll(todolist, todolistPath)
    if err != nil {
      fmt.Println("Error deleting all items!", err)
    }
  }

	if flag.NFlag() == 0 {
		PrintList(todolist.Items)
	}
	numComplete := CheckCompleted(todolist)
	if numComplete > 0 {
		fmt.Printf("\t !!! You have completed %d items !!!", numComplete)
		RemoveCompleted(todolist, todolistPath)
	}
}
