# Todo-CMD - The simple command-line Todo List App!

This app acts as a simple program to add items to a to-do list and modify that list in many ways!

## Installation
- Clone this repository anywhere you'd like
- Make sure all files are up to date
- In the directory, run the command
```
go build -o todo main.go
```
- For Linux users, use command:
```
sudo mv usr/local/bin ./todo
```
to be able to use this program anywhere on your PC!

## Usage
- After a successful installation, you can run the app in terminal using the commands `todo` (if you have the file in your usr/local/bin) or `./todo`
- Simply running the command will display the program header and whatever is currently in your To-Do list
- To add an item to your list, use the flag `-add` to specify that you would like to add an item
    - You can then enter the title, description and priority of that item
- To remove an item from the list, use the flag `-r` followed by the title (case-insensitive) of the item to be removed
    - To remove all items from the list, should you ever wish to do so, can be done using the `-ra` flag
    - Don't worry if you accidentally type this command though, the program will ask you to make sure that you really want to delete everything!
- To mark an item as completed, you can use the `-c` flag followed by the title (case-insensitive) of the item that has been completed
    - This will mark the item as completed
    - When you print the list again, you will be prompted as to whether you would like to remove any completed items from the list
    - To complete all items in the list, you can simply use the `-ca` flag
- To display the items in a particular order, use the `-s` flag followed by `title`, `completed`, or `priority` to print the list ordered by the respective variable

## Setup
The program makes use of a YAML file containing the To-Do list, and if no file exists, it will create one

