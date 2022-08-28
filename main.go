package main

import (
	"btrain/pkg/storage"
	"fmt"
)

func main() {
	password := "password"

	var dbURL string = fmt.Sprintf("postgres://postgres:" + password + "@65.108.250.159:5432/tasks")
	db, err := storage.New(dbURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	id, err := db.NewTask(storage.CreateTask(0, 0, "Do sth", "Do something"))
	catchErr(err)
	tasks, err := db.Tasks(0, 0)
	catchErr(err)
	fmt.Println(tasks)
	db.UpdateTask(id, 0, "Do sth New", "Do something New")
	tasks, err = db.Tasks(0, 0)
	catchErr(err)
	fmt.Println(tasks)
	db.DeleteTask(id)
	tasks, err = db.Tasks(0, 0)
	catchErr(err)
	fmt.Println(tasks)
}

func catchErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
