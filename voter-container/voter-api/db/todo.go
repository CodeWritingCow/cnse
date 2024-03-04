package db

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ToDoItem struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	IsDone bool   `json:"done"`
}

// DbMap is a type alias for a map of ToDoItems.  The key
// will be the ToDoItem.Id and the value will be the ToDoItem
type DbMap map[int]ToDoItem

// ToDo is the struct that represents the main object of our
// todo app.  It contains a map of ToDoItems and the name of
// the file that is used to store the items.
//
// This is just a mock, so we will only be managing an in memory
// map
type ToDo struct {
	toDoMap DbMap
	//more things would be included in a real implementation
}

// New is a constructor function that returns a pointer to a new
// ToDo struct.  It takes a single string argument that is the
// name of the file that will be used to store the ToDo items.
// If the file doesn't exist, it will be created.  If the file
// does exist, it will be loaded into the ToDo struct.
func New() (*ToDo, error) {

	//Now that we know the file exists, at at the minimum we have
	//a valid empty DB, lets create the ToDo struct
	toDo := &ToDo{
		toDoMap: make(map[int]ToDoItem),
	}

	// We should be all set here, the ToDo struct is ready to go
	// so we can support the public database operations
	return toDo, nil
}

// DeleteItem accepts an item id and removes it from the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be removed from the DB
//		(2) The DB file will be saved with the item removed
//		(3) If there is an error, it will be returned
func (t *ToDo) DeleteItem(id int) error {

	// we should if item exists before trying to delete it
	// this is a good practice, return an error if the
	// item does not exist

	//Now lets use the built-in go delete() function to remove
	//the item from our map
	delete(t.toDoMap, id)

	return nil
}

// DeleteAll removes all items from the DB.
// It will be exposed via a DELETE /todo endpoint
func (t *ToDo) DeleteAll() error {
	//To delete everything, we can just create a new map
	//and assign it to our existing map.  The garbage collector
	//will clean up the old map for us
	t.toDoMap = make(map[int]ToDoItem)

	return nil
}

// UpdateItem accepts a ToDoItem and updates it in the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be updated in the DB
//		(2) The DB file will be saved with the item updated
//		(3) If there is an error, it will be returned
func (t *ToDo) UpdateItem(item ToDoItem) error {

	// Check if item exists before trying to update it
	// this is a good practice, return an error if the
	// item does not exist
	_, ok := t.toDoMap[item.Id]
	if !ok {
		return errors.New("item does not exist")
	}

	//Now that we know the item exists, lets update it
	t.toDoMap[item.Id] = item

	return nil
}

// PrintItem accepts a ToDoItem and prints it to the console
// in a JSON pretty format. As some help, look at the
// json.MarshalIndent() function from our in class go tutorial.
func (t *ToDo) PrintItem(item ToDoItem) {
	jsonBytes, _ := json.MarshalIndent(item, "", "  ")
	fmt.Println(string(jsonBytes))
}

// PrintAllItems accepts a slice of ToDoItems and prints them to the console
// in a JSON pretty format.  It should call PrintItem() to print each item
// versus repeating the code.
func (t *ToDo) PrintAllItems(itemList []ToDoItem) {
	for _, item := range itemList {
		t.PrintItem(item)
	}
}

// JsonToItem accepts a json string and returns a ToDoItem
// This is helpful because the CLI accepts todo items for insertion
// and updates in JSON format.  We need to convert it to a ToDoItem
// struct to perform any operations on it.
func (t *ToDo) JsonToItem(jsonString string) (ToDoItem, error) {
	var item ToDoItem
	err := json.Unmarshal([]byte(jsonString), &item)
	if err != nil {
		return ToDoItem{}, err
	}

	return item, nil
}
